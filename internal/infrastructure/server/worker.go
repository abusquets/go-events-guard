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

// NewEventWorker crea una nova instància de EventWorker com a punter.
func NewEventWorker(signalsBus signals.SignalsBus) Worker {
	return &EventWorker{
		signalsBus: signalsBus,
	}
}

// Start és la funció que inicia el processament d'esdeveniments.
func (w *EventWorker) Start(ctx context.Context) {
	log.Println("EventWorker started listening for events...")

	// Aquí subscrivim el worker a un tema específic del bus
	w.signalsBus.Subscribe("usuari:creat", func(args []interface{}) error {
		// Aquí pots processar els esdeveniments que arriben al worker
		log.Printf("EventWorker processing event: %v\n", args)
		return nil
	})

	// Simulem escoltar esdeveniments en un bucle infinit
	for {
		select {
		case <-ctx.Done():
			log.Println("EventWorker shutting down")
			return
		default:
			// Esperar un moment abans de continuar per evitar bloquejar el CPU
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// RegisterHooks afegeix els hooks del worker al cicle de vida d'Fx.
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
				// Aquí pots afegir lògica de parada, si cal
				return nil
			},
		},
	)
}
