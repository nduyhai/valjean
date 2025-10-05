package main

import (
	"github.com/nduyhai/valjean/internal/infra/fxmodules"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		fxmodules.AllModules,
		fxmodules.ServerModule,
		fxmodules.LifecycleModule,
	)
	app.Run()
}
