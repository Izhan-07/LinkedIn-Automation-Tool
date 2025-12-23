package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"

	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/config"
	"linkedin-automation/pkg/logger"
)

type Authenticator struct {
	browser *browser.Browser
	config  config.LinkedInConfig
	log     *zap.Logger
}

func New(b *browser.Browser, cfg config.LinkedInConfig) *Authenticator {
	return &Authenticator{
		browser: b,
		config:  cfg,
		log:     logger.Get().Named("auth"),
	}
}

// Login attempts to log in or use existing session
func (a *Authenticator) Login(ctx context.Context) error {
	a.log.Info("Starting login process...")

	// 1. Try to load cookies
	if err := a.loadCookies(); err == nil {
		a.log.Info("Restored cookies, verifying session...")
		if a.checkSession() {
			a.log.Info("Session valid.")
			return nil
		}
		a.log.Warn("Session invalid, performing full login.")
	} else {
		a.log.Info("No session found, performing full login details:", zap.String("Error", err.Error()))
	}

	// 2. Perform fresh login
	if err := a.doLogin(); err != nil {
		return err
	}

	// 3. Save cookies
	return a.saveCookies()
}

func (a *Authenticator) doLogin() error {
	page := a.browser.Page

	a.log.Info("Navigating to login page...")
	if err := a.browser.NavigateTo("https://www.linkedin.com/login"); err != nil {
		return fmt.Errorf("navigation failed: %w", err)
	}

	// Stealth: Random delay
	time.Sleep(time.Duration(2000) * time.Millisecond)

	a.log.Info("Entering credentials...")

	// Username
	elUser, err := page.Element("#username")
	if err != nil {
		return fmt.Errorf("username field not found: %w", err)
	}
	if err := a.browser.Mouse.Click(elUser); err != nil {
		return err
	}
	if err := elUser.Type(a.config.Username); err != nil {
		return err
	}
	time.Sleep(time.Duration(500) * time.Millisecond)

	// Password
	elPass, err := page.Element("#password")
	if err != nil {
		return fmt.Errorf("password field not found: %w", err)
	}
	if err := a.browser.Mouse.Click(elPass); err != nil {
		return err
	}
	if err := elPass.Type(a.config.Password); err != nil {
		return err
	}
	time.Sleep(time.Duration(1000) * time.Millisecond)

	// Submit
	elBtn, err := page.Element("button[type=submit]")
	if err != nil {
		return fmt.Errorf("submit button not found: %w", err)
	}
	if err := a.browser.Mouse.Click(elBtn); err != nil {
		return err
	}

	// Wait for navigation and check for challenges
	a.log.Info("Waiting for navigation...")
	page.MustWaitLoad()

	// Check for checkpoints (2FA, captcha)
	// Simple check: are we on feed or challenge?
	url := page.MustInfo().URL
	if a.isChallengePage(url) {
		a.log.Warn("Challenge detected! Manual intervention might be required.", zap.String("url", url))
		// For this assignment, we might wait or fail.
		// Real implementation would pause and notify user.
		// Let's wait a bit to see if user manually solves it if headful.
		time.Sleep(10 * time.Second)
	}

	if !a.checkSession() {
		return errors.New("login failed: session not valid after attempt")
	}

	return nil
}

func (a *Authenticator) isChallengePage(url string) bool {
	// naive check
	return Contains(url, "challenge") || Contains(url, "checkpoint")
}

func (a *Authenticator) checkSession() bool {
	// Navigate to feed and check for "start a post" or nav bar
	// For speed, just check if we are on feed
	if err := a.browser.NavigateTo("https://www.linkedin.com/feed/"); err != nil {
		return false
	}

	// Check for a specific element that exists only when logged in
	// e.g., the "Me" icon or global nav
	_, err := a.browser.Page.Timeout(5 * time.Second).Element(".global-nav__me")
	return err == nil
}

func (a *Authenticator) saveCookies() error {
	cookies, err := a.browser.Page.Cookies([]string{})
	if err != nil {
		return err
	}

	data, err := json.Marshal(cookies)
	if err != nil {
		return err
	}

	return os.WriteFile("cookies.json", data, 0644)
}

func (a *Authenticator) loadCookies() error {
	data, err := os.ReadFile("cookies.json")
	if err != nil {
		return err
	}

	var cookies []*proto.NetworkCookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return err
	}

	return a.browser.Page.SetCookies(cookies)
}

// Utility
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
