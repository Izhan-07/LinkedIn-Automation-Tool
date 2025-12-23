package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/config"
	"linkedin-automation/internal/connection"
	"linkedin-automation/internal/search"
	"linkedin-automation/internal/store"
	"linkedin-automation/pkg/logger"
)

func main() {
	// Parse CLI flags
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	task := flag.String("task", "search-connect", "Task to run: search-connect, message")
	keywords := flag.String("keywords", "", "Search keywords (for search-connect)")
	limit := flag.Int("limit", 10, "Max actions to perform")
	flag.Parse()

	// 1. Load Config
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Init Logger
	if err := logger.Init(cfg.App.Debug); err != nil {
		panic(err)
	}
	log := logger.Get()
	defer log.Sync()

	log.Info("Starting LinkedIn Automation Bot", zap.String("task", *task))

	// 3. Init Store
	s, err := store.New(cfg.Database.Path)
	if err != nil {
		log.Fatal("Failed to init database", zap.Error(err))
	}

	// 4. Init Browser
	b, err := browser.New(cfg.Browser)
	if err != nil {
		log.Fatal("Failed to init browser", zap.Error(err))
	}
	defer b.Close()

	// 5. Auth
	authenticator := auth.New(b, cfg.LinkedIn)
	ctx := context.Background()
	if err := authenticator.Login(ctx); err != nil {
		log.Fatal("Authentication failed", zap.Error(err))
	}

	// 6. Execute Task
	switch *task {
	case "search-connect":
		if *keywords == "" {
			log.Fatal("Keywords required for search-connect task")
		}
		runSearchConnect(b, s, cfg, *keywords, *limit)
	case "message":
		log.Warn("Messaging task not fully wired in CLI yet (logic implemented)")
	default:
		log.Fatal("Unknown task", zap.String("task", *task))
	}

	log.Info("Task completed successfully.")
}

func runSearchConnect(b *browser.Browser, s *store.Store, cfg *config.Config, keywords string, limit int) {
	log := logger.Get()

	// Init Search
	searchEngine := search.New(b, s)
	criteria := search.SearchCriteria{
		Keywords: keywords,
		Count:    limit,
	}

	profiles, err := searchEngine.Run(criteria)
	if err != nil {
		log.Error("Search failed", zap.Error(err))
		return
	}

	// Init Connector
	connector := connection.New(b, s, cfg.Limits)

	count := 0
	for _, p := range profiles {
		if count >= limit {
			break
		}

		log.Info("Processing profile", zap.String("url", p))

		// Connect with a generic note
		note := "Hi, I'd like to connect to discuss potential opportunities."
		if err := connector.Connect(p, note); err != nil {
			log.Warn("Failed to connect", zap.String("profile", p), zap.Error(err))
			continue
		}

		count++
		// Cooldown between actions
		time.Sleep(10 * time.Second)
	}
}
