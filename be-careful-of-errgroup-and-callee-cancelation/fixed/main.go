package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

// This is the source code used for the post https://peng.fyi/post/lessons-from-an-errgroup-and-context-mishap/

func main() {
	group, ctx := errgroup.WithContext(context.Background())

	// set up 3 managers.
	managers := make([]*Manager, 3)
	for i := range 3 {
		managers[i] = &Manager{
			id:     fmt.Sprintf("m-%d", i),
			ctx:    ctx,
			taskCh: make(chan func(), 10),
			stopCh: make(chan func()),
		}
		// simulate tasks to each manager.
		for range 5 {
			go func() {
				managers[i].taskCh <- func() {
					time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
				}
			}()
		}
	}

	// set up the director
	d := &Director{
		ctx:            ctx,
		group:          group,
		controlSession: &ControlSession{},
		managers:       managers,
	}

	d.Start()
}

type Director struct {
	ctx            context.Context
	group          *errgroup.Group
	managers       []*Manager
	controlSession *ControlSession
}

func (d *Director) Start() {
	for _, m := range d.managers {
		m := m
		d.group.Go(func() error {
			m.Start()
			return nil
		})
	}

	ctx, cancel := context.WithCancel(d.ctx)
	d.group.Go(func() error {
		defer cancel()
		d.monitorAnomaly()
		return nil
	})

	d.group.Go(func() error {
		d.controlSession.Start(ctx)
		return nil
	})

	doneCh := make(chan error)
	go func() {
		doneCh <- d.group.Wait()
	}()

	for {
		select {
		case <-d.ctx.Done():
			fmt.Println("Stopped director: context canceled")
			return
		case err := <-doneCh:
			fmt.Printf("Stopped director: %v\n", err)
			return
		}
	}
}

func (d *Director) monitorAnomaly() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if d.anomalyDetected() {
				fmt.Println("Anomaly detected.")
				d.notifyManagers()
				return
			}
		}
	}
}

func (d *Director) notifyManagers() {
	for _, m := range d.managers {
		m := m
		m.stopCh <- func() {
			fmt.Printf("Important message for %s!. Anomaly Detected.\n", m.id)
		}
	}
}

func (d *Director) anomalyDetected() bool {
	return true
}

// Manager is an actor that processes actions from the actionCh channel.
type Manager struct {
	ctx context.Context

	id string

	// a buffered channel for regular tasks
	taskCh chan func()

	// a unbuffered channel for the function to execute before exit
	stopCh chan func()
}

func (m *Manager) Start() {
	for {
		select {
		case f := <-m.taskCh:
			f()
		case f := <-m.stopCh:
			f()
			fmt.Printf("Stopped manager %s: stop function received.\n", m.id)
			return
		case <-m.ctx.Done():
			fmt.Printf("Stopped manager %s: context canceled.\n", m.id)
			return
		}
	}
}

// ControlSession represents a session to the control plane.
type ControlSession struct{}

func (s *ControlSession) Start(ctx context.Context) error {
	for {
		select {
		// For simplicity, we leave out activities of the session.
		case <-ctx.Done():
			fmt.Println("Stopped the control session.")
			return ctx.Err()
		}
	}
}
