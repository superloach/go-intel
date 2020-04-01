package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/superloach/intel"
)

const defVers = "a8ca614df70e09516b36f060ef0304464e29dc75"

var (
	conc    = flag.Int("conc", 10, "number of concurrent jobs")
	timeout = flag.Duration("timeout", time.Second*5, "max request timeout")
	logint  = flag.Int("logint", 25, "interval to log progress")

	base   = flag.String("base", "intel.ingress.com", "intel site url")
	secure = flag.Bool("secure", true, "use https")
	ua     = flag.String("ua", "Foo Bar Browser", "user agent")
	vers   = flag.String("version", defVers, "internal intel version")
	csrf   = flag.String("csrf", "", "csrf token")
	sessid = flag.String("sessid", "", "google(?) session id")

	lat  = flag.Float64("lat", 0.0, "latitude")
	lng  = flag.Float64("lng", 0.0, "longitude")
	zoom = flag.Int("zoom", 17, "zoom level 0-17")
)

func dedup(o []string) []string {
	n := make([]string, len(o))
	copy(n, o)

	if len(n) < 2 {
		return n
	}

	j := 0
	i := 1

	for i < len(n) {
		if n[i] != n[j] {
			j++
			n[j] = n[i]
		}
		i++
	}

	return n[:j+1]
}

func main() {
	flag.Parse()

	if *csrf == "" {
		panic("please provide a csrf token")
	}
	if *sessid == "" {
		panic("please provide a sessid")
	}

	client, err := intel.NewClient()
	if err != nil {
		panic(err)
	}
	client.Client.Timeout = *timeout
	client.Base = *base
	client.Secure = *secure
	client.UA = *ua
	client.Version = *vers
	client.CSRF = *csrf
	client.SessID = *sessid

	tileKey := intel.TileKey(*lat, *lng, *zoom)
	fmt.Println("tile key:", tileKey)

	portalIDs, err := client.PortalIDs([]string{tileKey})
	if err != nil {
		panic(err)
	}
	portalIDs = dedup(portalIDs)

	jobs := make(chan struct{}, *conc)
	done := make(chan struct{})

	portals := make([]*intel.Portal, 0)
	for _, portalID := range portalIDs {
		go func() {
			jobs <- struct{}{}
			portal, err := client.GetPortal(portalID)
			if err != nil {
				panic(err)
			}

			portals = append(portals, portal)
			done <- <-jobs
		}()
	}

	l := len(portalIDs)
	for range done {
		i := len(portals)
		if i % *logint == 0 {
			fmt.Printf("%d/%d (%.1f%%)\n", i, l, float64(i)/float64(l)*100)
		}
		if i == l {
			close(done)
		}
		i++
	}
}
