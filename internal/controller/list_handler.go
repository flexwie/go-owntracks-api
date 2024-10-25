package controller

import (
	"html/template"
	"net/http"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/flexwie/owntracks-api/internal/db"
	"github.com/jmoiron/sqlx"
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

	locations := []db.Location{}
	err := c.Db.Select(&locations, "select * FROM location ORDER BY created DESC LIMIT 100")
	if err != nil {
		returnError(err, w)
		return
	}

	tmpl, err := template.New("list").Parse(t)
	if err != nil {
		returnError(err, w)
		return
	}

	type data struct {
		Locations []db.Location
	}

	err = tmpl.Execute(w, &data{
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
