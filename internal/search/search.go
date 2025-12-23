package search

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.uber.org/zap"

	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/store"
	"linkedin-automation/pkg/logger"
)

type Engine struct {
	browser *browser.Browser
	store   *store.Store
	log     *zap.Logger
}

func New(b *browser.Browser, s *store.Store) *Engine {
	return &Engine{
		browser: b,
		store:   s,
		log:     logger.Get().Named("search"),
	}
}

type SearchCriteria struct {
	Keywords string
	Type     string // "people", "jobs", etc.
	Count    int    // Number of profiles to scrape
}

// Run performs a search and collects profile URLs
func (e *Engine) Run(criteria SearchCriteria) ([]string, error) {
	e.log.Info("Starting search...", zap.String("keywords", criteria.Keywords))

	encodedKeywords := strings.ReplaceAll(criteria.Keywords, " ", "%20")
	// search URL for People
	url := fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s&origin=GLOBAL_SEARCH_HEADER", encodedKeywords)

	if err := e.browser.NavigateTo(url); err != nil {
		return nil, err
	}

	var collectedProfiles []string
	page := 1

	for len(collectedProfiles) < criteria.Count {
		e.log.Info("Processing search page", zap.Int("page", page))

		// Wait for results
		e.browser.Page.MustWaitLoad()
		time.Sleep(3 * time.Second) // Randomize this in prod

		// Scroll to load all items (lazy loading)
		if err := e.scrollPage(); err != nil {
			e.log.Warn("Failed to scroll page", zap.Error(err))
		}

		// Scrape current page
		profiles, err := e.scrapeProfiles()
		if err != nil {
			e.log.Error("Failed to scrape profiles", zap.Error(err))
		} else {
			e.log.Info("Found profiles on page", zap.Int("count", len(profiles)))
			for _, p := range profiles {
				if len(collectedProfiles) >= criteria.Count {
					break
				}
				// Dedupe in memory for this run (Database dedupe happens later or here)
				if !contains(collectedProfiles, p) {
					collectedProfiles = append(collectedProfiles, p)
				}
			}
		}

		if len(collectedProfiles) >= criteria.Count {
			break
		}

		// Pagination
		hasNext, err := e.nextPage()
		if err != nil {
			e.log.Warn("Failed to navigate next page", zap.Error(err))
			break
		}
		if !hasNext {
			e.log.Info("No more pages.")
			break
		}
		page++
	}

	return collectedProfiles, nil
}

func (e *Engine) scrapeProfiles() ([]string, error) {
	// Selectors typically look like .reusable-search__result-container or .entity-result__title-text a
	// Note: Class names change often. Using stable attribute selectors is better if possible.
	// As of 2024, likely .app-aware-link or similar inside search results.

	els, err := e.browser.Page.Elements(".reusable-search__result-container a.app-aware-link")
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, el := range els {
		link, err := el.Attribute("href")
		if err != nil || link == nil {
			continue
		}
		// Clean URL (remove query params)
		clean := strings.Split(*link, "?")[0]
		if !strings.Contains(clean, "/in/") {
			continue // skip non-profile links
		}
		urls = append(urls, clean)
	}
	return urls, nil
}

func (e *Engine) scrollPage() error {
	// Human-like random scroll
	// Scroll down in chunks
	for i := 0; i < 5; i++ {
		e.browser.Page.Mouse.Scroll(0, float64(rand.Intn(300)+200), 1)
		time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
	}
	// Scroll back up a tiny bit (human behavior)
	e.browser.Page.Mouse.Scroll(0, -100, 1)
	time.Sleep(1 * time.Second)
	return nil
}

func (e *Engine) nextPage() (bool, error) {
	// Find "Next" button
	// usually generic button with assertable text or aria-label="Next"
	btn, err := e.browser.Page.Element("button[aria-label='Next']")
	if err != nil {
		// Try finding by text if aria fails
		return false, nil // Assume end of list
	}

	if disabled, _ := btn.Attribute("disabled"); disabled != nil {
		return false, nil
	}

	if err := e.browser.Mouse.Click(btn); err != nil {
		return false, err
	}

	return true, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
