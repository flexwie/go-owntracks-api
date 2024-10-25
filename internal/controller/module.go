package controller

import (
	"github.com/flexwie/go-common/http"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(http.AsRoute(NewCreatePointHandler)),
	fx.Provide(http.AsRoute(NewGetListHandler)),
	fx.Provide(http.AsRoute(NewMetricsHandler)),
)
