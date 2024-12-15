package signals

import (
	"eventsguard/internal/infrastructure/signals"
	"fmt"
)

type AuthUserDelayedSignalsHandler struct{}

func NewAuthUserDelayedSignalsHandler(
	asignalsBus signals.SignalsBus,
) *AuthUserDelayedSignalsHandler {
	asignalsBus.Subscribe("user:updated", func(args []interface{}) error {
		fmt.Println("Has cridat Callback 1 per user:updated:", args)
		return nil
	})

	asignalsBus.Subscribe("user:updated", func(args []interface{}) error {
		fmt.Println("Has cridat Callback 2 per user:updated:", args)
		return nil
	})

	return &AuthUserDelayedSignalsHandler{}
}
