package signals

import (
	"go.uber.org/fx"
)

var Module = fx.Module("signals",
	fx.Provide(NewSignalsBus),
)
