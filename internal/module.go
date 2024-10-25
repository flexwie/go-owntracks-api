package internal

import (
	"github.com/flexwie/owntracks-api/internal/controller"
	"github.com/flexwie/owntracks-api/internal/db"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"tailscale.com/tsnet"

	_ "github.com/lib/pq"
)

var Modules = fx.Options(
	WithTsHttp,
	controller.Modules,
	db.WithDb,
	fx.Invoke(func(*tsnet.Server) {}),
	fx.Invoke(func(*sqlx.DB) {}),
)
