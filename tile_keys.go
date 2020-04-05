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
	a := lat * math.Pi / 180
	b := math.Tan(a)
	c := 1 / math.Cos(a)
	d := math.Log(b + c)
	e := 1 - d / math.Pi
	f := edge(zoom)
	g := e / 2 * f
	h := math.Floor(g)
	i := int(h)
	return i
}

func lngToTile(lng float64, zoom int) int {
	a := edge(zoom)
	b := (lng + 180) / 360 * a
	c := math.Floor(b)
	d := int(c)
	return d
}

func TileKey(lat, lng float64, zoom int) string {
	return fmt.Sprintf(
		"%d_%d_%d_0_8_100",
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
