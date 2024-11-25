package controller

import (
	"html/template"
	"iter"
	"net/http"
	"slices"
	"time"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/flexwie/owntracks-api/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/twpayne/go-polyline"
)

//go:embed list.html.tmpl
var t string

type GetListHandler struct {
	Log *log.Logger
	Db  *sqlx.DB
}

func NewGetListHandler(logger *log.Logger, db *sqlx.DB) *GetListHandler {
	return &GetListHandler{
		Log: logger.WithPrefix("GetList"),
		Db:  db,
	}
}

func (c *GetListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Log.Debug("got list request")

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

	// get date from query
	date := time.Now()
	var err error
	dateParam := r.URL.Query().Get("date")
	if dateParam != "" {
		date, err = time.Parse(time.DateOnly, dateParam)
		if err != nil {
			returnError(err, w)
			return
		}
	}

	c.Log.Info("fetching locations", "date", date.Format(time.DateOnly))

	locations := []db.Location{}
	err = c.Db.Select(&locations, "select * FROM location WHERE username=$1 AND created::date >= to_date($2, 'YYYY-mm-dd') ORDER BY created DESC", username, date.Format(time.DateOnly))
	if err != nil {
		returnError(err, w)
		return
	}

	coords := [][]float64{}
	for n := range Map(slices.Values(locations), func(loc db.Location) []float64 {
		return []float64{float64(loc.Lat), float64(loc.Lng)}
	}) {
		coords = append(coords, n)
	}
	line := string(polyline.EncodeCoords(coords))

	tmpl, err := template.New("list").Parse(t)
	if err != nil {
		returnError(err, w)
		return
	}

	type data struct {
		Line      string
		Locations []db.Location
	}

	err = tmpl.Execute(w, &data{
		Line:      line,
		Locations: locations,
	})
	if err != nil {
		returnError(err, w)
		return
	}

}

func (*GetListHandler) Pattern() string {
	return "/list"
}

func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for a := range seq {
			if !yield(f(a)) {
				return
			}
		}
	}
}
