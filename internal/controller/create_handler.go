package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
)

var locationPosted = promauto.NewCounter(prometheus.CounterOpts{
	Name: "owntracks_locations_received",
	Help: "Total number of location points received",
})

type CreatePointHandler struct {
	Log *log.Logger
	Db  *sqlx.DB
}

func NewCreatePointHandler(logger *log.Logger, db *sqlx.DB) *CreatePointHandler {
	return &CreatePointHandler{
		Log: logger.WithPrefix("CreatePoint"),
		Db:  db,
	}
}

func (c *CreatePointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Log.Debug("got point creation")

	// parse dto from body
	var locationDto LocationDto
	err := json.NewDecoder(r.Body).Decode(&locationDto)
	if err != nil {
		returnError(err, w)
		return
	}

	// make sure its a location request
	if locationDto.Type != "location" {
		c.Log.Debug("not a location request", "type", locationDto.Type)
		return
	}

	// TODO: create auth middleware for this
	// fetch user from header
	userHeaderName := viper.GetString("user-header")
	if userHeaderName == "" {
		c.Log.Error("no user header provided in config")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username := r.Header.Get(userHeaderName)
	if username == "" {
		c.Log.Error("no user header provided in request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//location := locationDto.ToModel(who.Node.Name)
	location := locationDto.ToModel(username)

	tx := c.Db.MustBegin()
	tx.NamedExec("INSERT INTO location(username, lat, lng, alt, created) VALUES (:username, :lat, :lng, :alt, :created)", location)
	tx.Commit()

	locationPosted.Inc()

	if _, err := io.Copy(w, r.Body); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
	}

	c.Log.Info("processed location")
}

func (*CreatePointHandler) Pattern() string {
	return "/"
}

func returnError(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
