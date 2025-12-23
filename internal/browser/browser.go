package browser

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"

	"linkedin-automation/internal/config"
	"linkedin-automation/internal/modules/mouse"
)

type Browser struct {
	Rod    *rod.Browser
	Page   *rod.Page
	Config config.BrowserConfig
	Mouse  *mouse.Mover
}

func New(cfg config.BrowserConfig) (*Browser, error) {
	// Launcher setup
	l := launcher.New().
		Headless(cfg.Headless).
		Devtools(false)

	if cfg.UserData != "" {
		l.UserDataDir(cfg.UserData)
	}

	// Masking User Agent (basic)
	l.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	url, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(url).MustConnect()

	// Create a new page (incognito or standard based on context, here standard)
	page := browser.MustPage()

	// Apply Stealth
	if cfg.Stealth {
		page.MustEvalOnNewDocument(stealth.JS)
	}

	// Set Viewport to standard desktop
	page.MustSetViewport(1920, 1080, 1, false)

	return &Browser{
		Rod:    browser,
		Page:   page,
		Config: cfg,
		Mouse:  mouse.New(page),
	}, nil
}

func (b *Browser) Close() {
	b.Rod.MustClose()
}

// NavigateTo wrapper with random delay
func (b *Browser) NavigateTo(url string) error {
	if err := b.Page.Navigate(url); err != nil {
		return err
	}
	// Wait for load + random delay
	b.Page.MustWaitLoad()
	time.Sleep(2 * time.Second)
	return nil
}
