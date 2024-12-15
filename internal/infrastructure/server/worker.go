package server

import (
	"context"
	"eventsguard/internal/infrastructure/signals"
	"log"
	"time"

	"go.uber.org/fx"
)

type EventWorker struct {
	signalsBus signals.SignalsBus
}

func NewEventWorker(signalsBus signals.SignalsBus) Worker {
	return &EventWorker{
		signalsBus: signalsBus,
	}
}

func (w *EventWorker) Start(ctx context.Context) {
	log.Println("EventWorker started listening for events...")

	for {
		select {
		case <-ctx.Done():
			log.Println("EventWorker shutting down")
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (w *EventWorker) RegisterHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Println("Starting EventWorker...")
				go w.Start(ctx)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Stopping EventWorker...")
				return nil
			},
		},
	)
}
