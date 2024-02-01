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
	param := parameter.Param

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
		Clusters: []model.Cluster{
			model.Cluster{
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
			model.Cluster{
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
			model.Cluster{
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
			model.Cluster{
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
			EnableScript: true,
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

func writeTags() {
	tags := map[string]string{
		"Environment": "product",
		"Region":      "hn",
		"Lidc":        "hna",
		"Service":     "test",
		"Cluster":     "hna-v",
	}

	path := "./superoptions/tags.json"
	err := writeJSON(path, tags)
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
	writeTags()
	writeMisc()
	writePreload()
}
