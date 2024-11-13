package controller

import (
	"time"

	"github.com/flexwie/owntracks-api/internal/db"
)

type LocationDto struct {
	Bat  int8    `json:"bat"`   // battery level
	Lon  float32 `json:"lon"`   // longitude
	Lat  float32 `json:"lat"`   // latitude
	Alt  float32 `json:"alt"`   // altitude
	Tid  string  `json:"tid"`   // tracker id
	Type string  `json:"_type"` // request type
	Tst  int32   `json:"tst"`   // timestamp
	Vel  float32 `json:"vel"`   // velocity
}

func (dto *LocationDto) ToModel(username string) *db.Location {
	created := time.Unix(int64(dto.Tst), 0)

	return &db.Location{
		Lat:      dto.Lat,
		Lng:      dto.Lon,
		Alt:      dto.Alt,
		Username: username,
		Vel:      dto.Vel,
		Created:  created,
	}
}
