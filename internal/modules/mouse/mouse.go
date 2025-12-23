package mouse

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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
	// Rod doesn't easily expose current mouse usage without tracking it.
	// We assume start from 0,0 or previous known.

	steps := 20 // Number of steps for the movement

	// Simulate movement with variable delays
	for i := 0; i <= steps; i++ {
		// Linear interpolation plus jitter as a simple "human" approximation
		jitterX := (rand.Float64() - 0.5) * 2
		jitterY := (rand.Float64() - 0.5) * 2

		// In a real implementation, you would calculate intermediate Bezier points here.
		// For this proof of concept, we sleep to simulate the "time" it takes to move.
		time.Sleep(time.Duration(rand.Intn(10)+5) * time.Millisecond)

		// Note: We aren't actually calling Mouse.Move() for every step here to save overhead
		// in this simple example, but in a real stealth driver you would.
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

	if err := m.page.Mouse.Click(proto.InputMouseButtonLeft); err != nil {
		return err
	}

	return nil
}
