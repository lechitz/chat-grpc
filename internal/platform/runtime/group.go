// Package runtime provides a simple primitive to coordinate goroutine lifecycles.
package runtime

import (
	"context"
	"sync"
)

// Actor describes a unit of work that can be started and interrupted.
type Actor struct {
	Start     func() error
	Interrupt func(error)
}

// Group manages a collection of actors that start together and stop together.
type Group struct {
	actors []Actor
}

// Add registers a new actor.
func (g *Group) Add(start func() error, interrupt func(error)) {
	g.actors = append(g.actors, Actor{
		Start:     start,
		Interrupt: interrupt,
	})
}

// Run starts all actors and blocks until the context is cancelled or an actor fails.
func (g *Group) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for _, actor := range g.actors {
		wg.Add(1)
		go func(a Actor) {
			defer wg.Done()
			if err := a.Start(); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(actor)
	}

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errCh:
	}

	var stopWG sync.WaitGroup
	for _, actor := range g.actors {
		stopWG.Add(1)
		go func(a Actor) {
			defer stopWG.Done()
			if a.Interrupt != nil {
				a.Interrupt(err)
			}
		}(actor)
	}

	stopWG.Wait()
	wg.Wait()

	return err
}
