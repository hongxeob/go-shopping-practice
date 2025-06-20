package main

import (
	"context"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(func() context.Context { return context.Background() }),
	)
}
