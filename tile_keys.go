package intel

import (
	"fmt"
	"math"
)

func edge(zoom int) (e float64) {
	if zoom < 0 {
		zoom = 0
	}
	if zoom > 17 {
		zoom = 17
	}

	switch zoom {
	case 0, 1, 2:
		e = 1
	case 3, 4:
		e = 40
	case 5, 6:
		e = 80
	case 7:
		e = 320
	case 8:
		e = 1000
	case 9, 10:
		e = 2000
	case 11:
		e = 4000
	case 12:
		e = 8000
	case 13, 14:
		e = 16000
	case 15, 16, 17:
		e = 32000
	default:
	}

	return
}

func latToTile(lat float64, zoom int) int {
	return int(math.Floor((1 - math.Log(math.Tan(lat*math.Pi/180)+1/math.Cos(lat*math.Pi/180))/math.Pi) / 2 * edge(zoom)))
}

func lngToTile(lng float64, zoom int) int {
	return int(math.Floor(lng+180) / 360 * edge(zoom))
}

func TileKey(lat, lng float64, zoom int) string {
	return fmt.Sprintf(
		"%d_%d_%d_8_8_100",
		zoom,
		lngToTile(
			lng,
			zoom,
		),
		latToTile(
			lat,
			zoom,
		),
	)
}
