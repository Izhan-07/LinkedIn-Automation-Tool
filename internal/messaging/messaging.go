package messaging

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"linkedin-automation/internal/browser"
	"linkedin-automation/pkg/logger"
)

type Messenger struct {
	browser *browser.Browser
	log     *zap.Logger
}

func New(b *browser.Browser) *Messenger {
	return &Messenger{
		browser: b,
		log:     logger.Get().Named("messaging"),
	}
}

// SendMessage sends a direct message to a profile
func (m *Messenger) SendMessage(profileURL, message string) error {
	m.log.Info("Sending message...", zap.String("url", profileURL))

	if err := m.browser.NavigateTo(profileURL); err != nil {
		return err
	}

	// Find "Message" button
	msgBtn, err := m.browser.Page.ElementR("button", "Message")
	if err != nil {
		return fmt.Errorf("message button not found: %w", err)
	}

	if err := m.browser.Mouse.Click(msgBtn); err != nil {
		return err
	}

	// Wait for chat box
	time.Sleep(1 * time.Second)

	// Selector for chat input area
	// usually .msg-form__contenteditable
	input, err := m.browser.Page.Element("div[role='textbox']")
	if err != nil {
		return err
	}

	if err := input.Type(message); err != nil {
		return err
	}

	// Send (Click Send button or Press Enter)
	sendBtn, err := m.browser.Page.Element("button[type='submit']")
	if err == nil {
		m.browser.Mouse.Click(sendBtn)
	} else {
		// Fallback to Enter key
		// m.browser.Page.KeyActions().Press(input.Enter).MustDo()
	}

	m.log.Info("Message sent.")
	return nil
}
