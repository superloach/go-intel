package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/superloach/go-intel"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const defVers = "a8ca614df70e09516b36f060ef0304464e29dc75"

var (
	// process options
	conc    = flag.Int("conc", 100, "number of concurrent jobs")
	timeout = flag.Duration("timeout", time.Second*5, "max request timeout")
	logint  = flag.Int("logint", 25, "interval to log progress")
	retries = flag.Int("retries", 10, "maximum times to retry on error")

	// network options
	base   = flag.String("base", "intel.ingress.com", "intel site url")
	secure = flag.Bool("secure", true, "use https")
	ua     = flag.String("ua", "Foo Bar Browser", "user agent")
	vers   = flag.String("version", defVers, "internal intel version")

	// auth options
	csrf   = flag.String("csrf", "", "csrf token")
	sessid = flag.String("sessid", "", "google(?) session id")

	// location options
	minlat = flag.Float64("minlat", 0.0, "minimum latitude")
	minlng = flag.Float64("minlng", 0.0, "minimum longitude")
	maxlat = flag.Float64("maxlat", 0.0, "maximum latitude")
	maxlng = flag.Float64("maxlng", 0.0, "maximum longitude")
	step   = flag.Float64("step", 0.0001, "lng/lat step")
	zoom   = flag.Int("zoom", 17, "zoom level 0-17")

	// db options
	dburl   = flag.String("dburl", "", "rethink url")
	dbname  = flag.String("dbname", "", "rethink db name")
	dbtable = flag.String("dbtable", "", "rethink table name")
)

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
	client.MaxTries = *retries

	client.Base = *base
	client.Secure = *secure
	client.UA = *ua
	client.Version = *vers

	client.CSRF = *csrf
	client.SessID = *sessid

	if *dbname == "" {
		panic("please provide a db name")
	}
	if *dbtable == "" {
		panic("please provide a db table")
	}

	db, err := r.Connect(r.ConnectOpts{
		Address: *dburl,
		Timeout: *timeout,
	})
	if err != nil {
		panic(err)
	}

	tileKeys := make([]string, 0)
	for lat := *minlat; lat <= *maxlat; lat += *step {
		for lng := *minlng; lng <= *maxlng; lng += *step {
			tileKey := intel.TileKey(lat, lng, *zoom)
			tileKeys = append(tileKeys, tileKey)
		}
	}
	tileKeys = intel.Dedup(tileKeys)

	portalIDMap := make(map[string]struct{})

	l := len(tileKeys)
	for i, tileKey := range tileKeys {
		fmt.Printf("%s - %d/%d (%.2f%%) - %d portals\n", tileKey, (i + 1), l, float64(i+1)/float64(l)*100, len(portalIDMap))

		portalIDs, err := client.PortalIDs([]string{tileKey})
		if err != nil {
			panic(err)
		}

		for _, id := range portalIDs {
			if _, ok := portalIDMap[id]; ok {
				continue
			}

			portalIDMap[id] = struct{}{}
		}
	}

	l = len(portalIDMap)

	jobs := make(chan struct{}, *conc)
	done := make(chan struct{})

	for portalID, _ := range portalIDMap {
		go func(guid string) {
			jobs <- struct{}{}

			portal, err := client.GetPortal(guid)
			if err != nil {
				<-jobs
				l--
				fmt.Println(err)
				return
			}

			_, err = r.DB(*dbname).Table(*dbtable).Insert(portal).Run(db)
			if err != nil {
				<-jobs
				l--
				fmt.Println(err)
				return
			}

			done <- <-jobs
		}(portalID)
	}

	i := 0
	for range done {
		i++
		if i%*logint == 0 {
			fmt.Printf("%d/%d (%.1f%%)\n", i, l, float64(i)/float64(l)*100)
		}
		if i == l {
			close(done)
		}
	}

	ol := len(portalIDMap)
	fmt.Printf(
		"lost %d of %d (%.1f%%)\n",
		ol-l, ol, float64(ol-l)/float64(ol)*100,
	)
}
