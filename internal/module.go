package internal

import (
	"net/http"

	fhttp "github.com/flexwie/go-common/http"
	"github.com/flexwie/owntracks-api/internal/controller"
	"github.com/flexwie/owntracks-api/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	_ "github.com/lib/pq"
)

func WithBusinessLogic() fx.Option {
	addr := viper.GetString("addr")

	return fx.Options(
		fhttp.WithHttpFactory(addr),
		controller.Modules,
		db.WithDb,
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(func(*sqlx.DB) {}),
	)
}
