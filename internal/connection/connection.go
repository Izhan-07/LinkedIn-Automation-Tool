package connection

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"go.uber.org/zap"

	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/config"
	"linkedin-automation/internal/store"
	"linkedin-automation/pkg/logger"
)

type Manager struct {
	browser *browser.Browser
	store   *store.Store
	config  config.LimitsConfig
	log     *zap.Logger
}

func New(b *browser.Browser, s *store.Store, cfg config.LimitsConfig) *Manager {
	return &Manager{
		browser: b,
		store:   s,
		config:  cfg,
		log:     logger.Get().Named("connection"),
	}
}

// Connect sends a connection request to a profile
func (m *Manager) Connect(profileURL string, note string) error {
	m.log.Info("Visiting profile for connection...", zap.String("url", profileURL))

	if err := m.browser.NavigateTo(profileURL); err != nil {
		return err
	}

	// Wait for profile render
	m.browser.Page.MustWaitLoad()

	// Random scroll to simulate reading
	m.browser.Page.Mouse.Scroll(0, 300, 1)
	time.Sleep(time.Duration(rand.Intn(1000)+1000) * time.Millisecond)

	// Localize "Connect" button
	// LinkedIn UI varies: Connect is either primary action or under "More"
	// Selector strategy: Look for text "Connect" or "More" -> "Connect"

	els, err := m.browser.Page.Elements("button")
	if err != nil {
		return err
	}

	var connectBtn *rod.Element
	for _, el := range els {
		text, _ := el.Text()
		if text == "Connect" {
			connectBtn = el
			break
		}
	}

	if connectBtn == nil {
		// Try "More" menu
		m.log.Info("Connect button not found, checking 'More' menu...")
		moreBtn, err := m.browser.Page.Element("button[aria-label='More actions']")
		if err == nil {
			m.browser.Mouse.Click(moreBtn)
			time.Sleep(500 * time.Millisecond)
			// Look for Connect in dropdown
			connectBtn, _ = m.browser.Page.ElementR("div[role='button']", "Connect")
		}
	}

	if connectBtn == nil {
		// Already connected or pending
		m.log.Warn("Connect button not found. Skipping.")
		return errors.New("already connected or locked")
	}

	// Click Connect
	if err := m.browser.Mouse.Click(connectBtn); err != nil {
		return err
	}

	// Handle "Add a note" modal
	time.Sleep(1 * time.Second)

	if note != "" {
		addNoteBtn, err := m.browser.Page.ElementR("button", "Add a note")
		if err == nil {
			m.browser.Mouse.Click(addNoteBtn)
			time.Sleep(500 * time.Millisecond)

			textArea, err := m.browser.Page.Element("textarea[name='message']")
			if err == nil {
				// Typing simulation
				textArea.Type(note)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	// Send
	sendBtn, err := m.browser.Page.ElementR("button", "Send")
	if err != nil {
		return fmt.Errorf("failed to find Send button: %w", err)
	}

	if err := m.browser.Mouse.Click(sendBtn); err != nil {
		return err
	}

	m.log.Info("Connection request sent!")
	return nil
}
