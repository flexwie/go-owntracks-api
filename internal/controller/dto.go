package controller

import "github.com/flexwie/owntracks-api/internal/db"

type LocationDto struct {
	Bat  int8    `json:"bat"`
	Lon  float32 `json:"lon"`
	Lat  float32 `json:"lat"`
	Alt  float32 `json:"alt"`
	Tid  string  `json:"tid"`
	Type string  `json:"_type"`
}

func (dto *LocationDto) ToModel(username string) *db.Location {
	return &db.Location{
		Lat:      dto.Lat,
		Lng:      dto.Lon,
		Alt:      dto.Alt,
		Username: username,
	}
}
