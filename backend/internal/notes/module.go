package notes

import "go.uber.org/fx"

var Module = fx.Module(
	"notes",
	fx.Provide(
		NewRepository,
		NewService,
		NewHandler,
	),
)
