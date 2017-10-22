package data

import "github.com/satori/go.uuid"

type Garbanzo struct {
	Id           int
	APIUUID      uuid.UUID
	GarbanzoType GarbanzoType
	DiameterMM   float32
}
