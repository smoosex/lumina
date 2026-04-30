package bookmeta

import "go.uber.org/fx"

var Module = fx.Module("bookmeta", fx.Provide(NewDoubanClient))
