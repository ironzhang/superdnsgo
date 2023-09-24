package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superlib/fileutil"
	"github.com/ironzhang/superlib/superutil/parameter"

	"github.com/ironzhang/superdnsgo"
	"github.com/ironzhang/superdnsgo/pkg/model"
	"github.com/ironzhang/superdnsgo/superdns"
)

func writeJSON(filename string, v interface{}) error {
	dir := filepath.Dir(filename)
	os.MkdirAll(dir, os.ModePerm)
	return fileutil.WriteJSON(filename, v)
}

func writeFile(filename string, data string) error {
	dir := filepath.Dir(filename)
	os.MkdirAll(dir, os.ModePerm)
	return ioutil.WriteFile(filename, []byte(data), 0666)
}

func writeSuperdnsCfg() {
	param := parameter.Parameter{
		AgentServer:   "127.0.0.1:1789",
		ResourcePath:  "/var/superdns",
		WatchInterval: 1,
	}

	path := "superdns.conf"
	err := fileutil.WriteTOML(path, param)
	if err != nil {
		tlog.Errorw("write toml", "path", path, "error", err)
		return
	}
	tlog.Debugf("write %s success", path)
}

func writeServiceModel() {
	m := model.ServiceModel{
		Domain: "www.superdns.com",
		Clusters: map[string]model.Cluster{
			"hna": model.Cluster{
				Name: "hna",
				Features: map[string]string{
					"Lidc":        "hna",
					"Region":      "hn",
					"Environment": "product",
				},
				Endpoints: []model.Endpoint{
					{Addr: "192.168.1.1:8000", State: model.Enabled, Weight: 100},
					{Addr: "192.168.1.2:8000", State: model.Enabled, Weight: 100},
				},
			},
			"hnb": model.Cluster{
				Name: "hnb",
				Features: map[string]string{
					"Lidc":        "hnb",
					"Region":      "hn",
					"Environment": "product",
				},
				Endpoints: []model.Endpoint{
					{Addr: "192.168.2.1:8000", State: model.Enabled, Weight: 100},
					{Addr: "192.168.2.2:8000", State: model.Enabled, Weight: 100},
				},
			},
			"hba@vip": model.Cluster{
				Name: "hba@vip",
				Features: map[string]string{
					"Lidc":        "hba",
					"Region":      "hb",
					"Environment": "product",
				},
				Endpoints: []model.Endpoint{
					{Addr: "127.0.0.1:8000", State: model.Enabled, Weight: 100},
				},
			},
			"default@mock": model.Cluster{
				Name: "default@mock",
				Features: map[string]string{
					"Lidc":        "default-lidc",
					"Region":      "default-region",
					"Environment": "product",
				},
				Endpoints: []model.Endpoint{
					{Addr: "mock.endpoint.com:8000", State: model.Enabled, Weight: 100},
				},
			},
		},
	}

	path := "./superdns/services/www.superdns.com.json"
	err := writeJSON(path, m)
	if err != nil {
		tlog.Errorw("write json", "path", path, "error", err)
		return
	}
	tlog.Debugf("write %s success", path)
}

func writeRouteModel() {
	m := model.RouteModel{
		Domain: "www.superdns.com",
		Strategy: model.RouteStrategy{
			EnableScriptRoute: false,
			RouteRules: map[string]model.RouteRule{
				"product": model.RouteRule{
					EnableScriptRoute: false,
					LidcDestinations: map[string][]model.Destination{
						"hna": []model.Destination{{Cluster: "hna", Percent: 1}},
						"hnb": []model.Destination{{Cluster: "hnb", Percent: 1}},
					},
					RegionDestinations: map[string][]model.Destination{
						"hn": []model.Destination{
							{Cluster: "hna", Percent: 0.5},
							{Cluster: "hnb", Percent: 0.5},
						},
					},
					EnvironmentDestinations: []model.Destination{
						{Cluster: "hba@vip", Percent: 1},
					},
				},
			},
			DefaultDestinations: []model.Destination{
				{Cluster: "default@mock", Percent: 1},
			},
		},
	}

	path := "./superdns/routes/www.superdns.com.json"
	err := writeJSON(path, m)
	if err != nil {
		tlog.Errorw("write json", "path", path, "error", err)
		return
	}
	tlog.Debugf("write %s success", path)
}

func writeRouteScript() {
	path := "./superdns/routes/www.superdns.com.lua"
	err := writeFile(path, luascript)
	if err != nil {
		tlog.Errorw("write file", "path", path, "error", err)
		return
	}
	tlog.Debugf("write %s success", path)
}

func writeZone() {
	z := superdns.Zone{
		Environment: "product",
		Region:      "hn",
		Lidc:        "hna",
		RouteTags: map[string]string{
			"Service": "test",
			"Cluster": "hna-v",
		},
	}

	path := "./superoptions/zone.json"
	err := writeJSON(path, z)
	if err != nil {
		tlog.Errorw("write json", "path", path, "error", err)
		return
	}
	tlog.Debugf("write %s success", path)
}

func writeMisc() {
	m := superdnsgo.Misc{
		SkipPreloadError: false,
	}

	path := "./superoptions/misc.json"
	err := writeJSON(path, m)
	if err != nil {
		tlog.Errorw("write json", "path", path, "error", err)
		return
	}
	tlog.Debugf("write %s success", path)
}

func writePreload() {
	domains := []string{"www.superdns.com"}

	path := "./superoptions/preload.json"
	err := writeJSON(path, domains)
	if err != nil {
		tlog.Errorw("write json", "path", path, "error", err)
		return
	}
	tlog.Debugf("write %s success", path)
}

func main() {
	writeSuperdnsCfg()
	writeServiceModel()
	writeRouteModel()
	writeRouteScript()
	writeZone()
	writeMisc()
	writePreload()
}
