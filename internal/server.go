package internal

import (
	"context"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	fhttp "github.com/flexwie/go-common/http"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"tailscale.com/tsnet"
)

var WithTsHttp = fx.Provide(
	newTsHttpServer,
	fx.Annotate(
		newServeMux,
		fx.ParamTags(`group:"routes"`),
	),
)

func newTsHttpServer(lc fx.Lifecycle, mux *http.ServeMux, logger *log.Logger) (*tsnet.Server, error) {
	authKey := viper.GetString("ts-auth-key")
	if authKey == "" {
		return nil, errors.New("ts-auth-key must be provided")
	}

	name := viper.GetString("ts-name")

	s := &tsnet.Server{
		Hostname: name,
		Logf:     logger.WithPrefix("ts").Debugf,
		UserLogf: logger.WithPrefix("ts").Infof,
		AuthKey:  authKey,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := s.ListenTLS("tcp", ":443")
			if err != nil {
				return err
			}

			go http.Serve(ln, mux)
			logger.Info("started tailscale http server", "addr", ":443", "hostname", s.Hostname)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Close()
		},
	})

	return s, nil
}

func newServeMux(routes []fhttp.Route, logger *log.Logger) *http.ServeMux {
	logger = logger.WithPrefix("routing")

	mux := http.NewServeMux()
	for _, route := range routes {
		logger.Debug("adding route", "from", route.Pattern(), "to", route)
		mux.Handle(route.Pattern(), route)
	}
	return mux
}
