package data

import "fmt"

//go:generate stringer -type=GarbanzoType

type GarbanzoType int

const (
	DESI   GarbanzoType = 1001
	KABULI GarbanzoType = 1002
)

func GarbanzoTypeFromString(gType string) (GarbanzoType, error) {
	switch gType {
	case DESI.String():
		return DESI, nil
	case KABULI.String():
		return KABULI, nil
	default:
		return 0, fmt.Errorf("invalid garbanzo type: %s", gType)
	}
}
