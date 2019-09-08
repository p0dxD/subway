package subway

import (
        "encoding/json"
        "fmt"
        "log"
        "math"
        "net/http"
        "strconv"
        "strings"
        "path/filepath"	
        "io/ioutil"

        rtree "github.com/dhconnelly/rtreego"
        geojson "github.com/paulmach/go.geojson"
)

// Stations is an RTree housing the stations
var Stations = rtree.NewTree(2, 25, 50)

// Station is a wrapper for `*geojson.Feature` so that we can implement
// `rtree.Spatial` interface type.
type Station struct {
        feature *geojson.Feature
}

// Bounds implements `rtree.Spatial` so we can load
// stations into an `rtree.Rtree`.
func (s *Station) Bounds() *rtree.Rect {
        return rtree.Point{
                s.feature.Geometry.Point[0],
                s.feature.Geometry.Point[1],
        }.ToRect(1e-6)
}

// GeoJSON is a cache of the NYC Subway Station and Line data.
var GeoJSON = make(map[string][]byte)

// cacheGeoJSON loads files under data into `GeoJSON`.
func cacheGeoJSON() {
        filenames, err := filepath.Glob("data/*")
        if err != nil {
                log.Fatal(err)
        }
        for _, f := range filenames {
                name := filepath.Base(f)
                dat, err := ioutil.ReadFile(f)
                if err != nil {
                        log.Fatal(err)
                }
                GeoJSON[name] = dat
        }
}

// Init is called from the App Engine runtime to initialize the app.
func Init() {
        cacheGeoJSON()
        loadStations()
        http.HandleFunc("/data/subway-stations", subwayStationsHandler)
        http.HandleFunc("/data/subway-lines", subwayLinesHandler)
}

func subwayLinesHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-type", "application/json")
        w.Write(GeoJSON["subway-lines.geojson"])
}
// loadStations loads the geojson features from
// `subway-stations.geojson` into the `Stations` rtree.
func loadStations() {
        stationsGeojson := GeoJSON["subway-stations.geojson"]
        fc, err := geojson.UnmarshalFeatureCollection(stationsGeojson)
        if err != nil {
                // Note: this will take down the GAE instance by exiting this process.
                log.Fatal(err)
        }
        for _, f := range fc.Features {
                Stations.Insert(&Station{f})
        }
}

// subwayStationsHandler reads r for a "viewport" query parameter
// and writes a GeoJSON response of the features contained in
// that viewport into w.
func subwayStationsHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-type", "application/json")
        vp := r.FormValue("viewport")
        rect, err := newRect(vp)
        if err != nil {
                str := fmt.Sprintf("Couldn't parse viewport: %s", err)
                http.Error(w, str, 400)
                return
        }
        s := Stations.SearchIntersect(rect)
        fc := geojson.NewFeatureCollection()
        for _, station := range s {
                fc.AddFeature(station.(*Station).feature)
        }
        err = json.NewEncoder(w).Encode(fc)
        if err != nil {
                str := fmt.Sprintf("Couldn't encode results: %s", err)
                http.Error(w, str, 500)
                return
        }
}

// newRect constructs a `*rtree.Rect` for a viewport.
func newRect(vp string) (*rtree.Rect, error) {
        ss := strings.Split(vp, "|")
        sw := strings.Split(ss[0], ",")
        swLat, err := strconv.ParseFloat(sw[0], 64)
        if err != nil {
                return nil, err
        }
        swLng, err := strconv.ParseFloat(sw[1], 64)
        if err != nil {
                return nil, err
        }
        ne := strings.Split(ss[1], ",")
        neLat, err := strconv.ParseFloat(ne[0], 64)
        if err != nil {
                return nil, err
        }
        neLng, err := strconv.ParseFloat(ne[1], 64)
        if err != nil {
                return nil, err
        }
        minLat := math.Min(swLat, neLat)
        minLng := math.Min(swLng, neLng)
        distLat := math.Max(swLat, neLat) - minLat
        distLng := math.Max(swLng, neLng) - minLng

        r, err := rtree.NewRect(
                rtree.Point{
                        minLng,
                        minLat,
                },
                []float64{
                        distLng,
                        distLat,
                })
        if err != nil {
                return nil, err
        }
        return r, nil
}