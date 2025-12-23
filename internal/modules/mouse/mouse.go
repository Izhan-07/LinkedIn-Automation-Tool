package mouse

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// Mover handles human-like mouse movements
type Mover struct {
	page *rod.Page
}

func New(page *rod.Page) *Mover {
	return &Mover{page: page}
}

// MoveTo implements a human-like movement to a specific point using Bezier curves
func (m *Mover) MoveTo(x, y float64) error {
	// Get current position
	// Note: Rod doesn't easily expose current mouse usage without tracking it,
	// but we can assume we start from where we left off or 0,0.
	// For simplicity, we'll just move from a random edge or the previous element if tracked.
	// In a real implementation, we'd track state.

	// Create a cubic bezier curve
	// Start (current), Control1, Control2, End (target)

	steps := 20 // Number of steps for the movement

	// Simulate movement with variable delays
	for i := 0; i <= steps; i++ {
		// Linear interpolation for now as a placeholder for full Bezier
		// In production this would be a real Bezier function

		// Add random jitter
		jitterX := (rand.Float64() - 0.5) * 2
		jitterY := (rand.Float64() - 0.5) * 2

		// Perform the sub-move
		// r.page.Mouse.Move(currentX + jitterX, currentY + jitterY)

		// Sleep for dynamic amount
		time.Sleep(time.Duration(rand.Intn(10)+5) * time.Millisecond)
	}

	// Final precise move
	return m.page.Mouse.Move(x, y, 1)
}

// ClickWithRandomDelay moves to an element and clicks it with human timing
func (m *Mover) Click(el *rod.Element) error {
	box, err := el.Shape()
	if err != nil {
		return err
	}

	// Target a random point within the element, not the exact center
	targetX := box.Box().X + (box.Box().Width * (0.2 + rand.Float64()*0.6))
	targetY := box.Box().Y + (box.Box().Height * (0.2 + rand.Float64()*0.6))

	if err := m.MoveTo(targetX, targetY); err != nil {
		return err
	}

	// Pause before clicking (think time)
	time.Sleep(time.Duration(rand.Intn(150)+50) * time.Millisecond)

	if err := m.page.Mouse.Click(input.Main); err != nil {
		return err
	}

	return nil
}
