package luaroute

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/ironzhang/superdnsgo/pkg/model"
)

func matchError(t testing.TB, err error, errstr string) bool {
	switch {
	case err == nil && errstr == "":
		return true
	case err != nil && errstr == "":
		return false
	case err == nil && errstr != "":
		return false
	case err != nil && errstr != "":
		matched, e := regexp.MatchString(errstr, err.Error())
		if e != nil {
			t.Fatalf("match: regexp match string: %v", e)
		}
		return matched
	}
	panic("never reach")
}

func TestPolicy(t *testing.T) {
	p := NewPolicy()

	// 加载脚本
	scripts := []string{"./testdata/test1.lua", "./testdata/test2.lua"}
	for _, script := range scripts {
		if err := p.Load(script); err != nil {
			t.Fatalf("load %q: %v", script, err)
		}
	}

	// 测试集群列表
	clusters := map[string]model.Cluster{
		"default@mock": model.Cluster{
			Name: "default@mock",
			Features: map[string]string{
				"Lidc":        "default_lidc",
				"Region":      "default_region",
				"Environment": "product",
			},
			Endpoints: make([]model.Endpoint, 10),
		},
		"hna-v": model.Cluster{
			Name: "hna-v",
			Features: map[string]string{
				"Lidc":        "hna",
				"Region":      "hn",
				"Environment": "product",
			},
			Endpoints: make([]model.Endpoint, 10),
		},
		"hnb-v": model.Cluster{
			Name: "hnb-v",
			Features: map[string]string{
				"Lidc":        "hnb",
				"Region":      "hn",
				"Environment": "product",
			},
			Endpoints: make([]model.Endpoint, 10),
		},
		"hbf-v": model.Cluster{
			Name: "hbf-v",
			Features: map[string]string{
				"Lidc":        "hbf",
				"Region":      "hb",
				"Environment": "product",
			},
			Endpoints: make([]model.Endpoint, 10),
		},
		"hna-sim000-v": model.Cluster{
			Name: "hna-sim000-v",
			Features: map[string]string{
				"Lidc":        "hna",
				"Region":      "hn",
				"Environment": "sim",
			},
			Endpoints: make([]model.Endpoint, 10),
		},
		"hna-sim001-v": model.Cluster{
			Name: "hna-sim001-v",
			Features: map[string]string{
				"Lidc":        "hna",
				"Region":      "hn",
				"Environment": "sim",
			},
			Endpoints: make([]model.Endpoint, 10),
		},
		"hna-sim002-v": model.Cluster{
			Name: "hna-sim002-v",
			Features: map[string]string{
				"Lidc":        "hna",
				"Region":      "hn",
				"Environment": "sim",
			},
			Endpoints: make([]model.Endpoint, 10),
		},
	}

	// 测试用例
	tests := []struct {
		domain       string
		tags         map[string]string
		clusters     map[string]model.Cluster
		destinations []model.Destination
		err          string
	}{
		{
			domain: "www.test1.com",
			tags: map[string]string{
				"Service":     "test-a",
				"Cluster":     "hna-v",
				"Lidc":        "hna",
				"Region":      "hn",
				"Environment": "product",
			},
			clusters: clusters,
			destinations: []model.Destination{
				{Cluster: "hna-v", Percent: 1},
			},
			err: "",
		},
		{
			domain: "www.test1.com",
			tags: map[string]string{
				"Service":     "test-a",
				"Cluster":     "hba-v",
				"Lidc":        "hba",
				"Region":      "hb",
				"Environment": "product",
			},
			clusters: clusters,
			destinations: []model.Destination{
				{Cluster: "hbf-v", Percent: 1},
			},
			err: "",
		},
		{
			domain: "www.test1.com",
			tags: map[string]string{
				"Service":        "test-a",
				"Cluster":        "hna-sim000-v",
				"Lidc":           "hna",
				"Region":         "hn",
				"Environment":    "sim",
				"X-Lane-Cluster": "hna-sim100-v",
			},
			clusters: clusters,
			destinations: []model.Destination{
				{Cluster: "hna-sim000-v", Percent: 1},
			},
			err: "",
		},
		{
			domain: "www.test1.com",
			tags: map[string]string{
				"Service":        "test-a",
				"Cluster":        "hna-sim000-v",
				"Lidc":           "hna",
				"Region":         "hn",
				"Environment":    "sim",
				"X-Lane-Cluster": "hna-sim001-v",
			},
			clusters: clusters,
			destinations: []model.Destination{
				{Cluster: "hna-sim001-v", Percent: 1},
			},
			err: "",
		},
		{
			domain: "www.test2.com",
			tags: map[string]string{
				"Service":     "test-a",
				"Cluster":     "hna-v",
				"Lidc":        "hna",
				"Region":      "hn",
				"Environment": "product",
			},
			clusters: clusters,
			destinations: []model.Destination{
				{Cluster: "hna-v", Percent: 0.5},
				{Cluster: "hnb-v", Percent: 0.5},
			},
			err: "",
		},
		{
			domain: "www.test2.com",
			tags: map[string]string{
				"Service":        "test-a",
				"Cluster":        "hna-sim000-v",
				"Lidc":           "hna",
				"Region":         "hn",
				"Environment":    "sim",
				"X-Lane-Cluster": "hna-sim001-v",
			},
			clusters: clusters,
			destinations: []model.Destination{
				{Cluster: "default@mock", Percent: 1},
			},
			err: "",
		},
		{
			domain: "www.not.find.com",
			tags: map[string]string{
				"Service":        "test-a",
				"Cluster":        "hna-sim000-v",
				"Lidc":           "hna",
				"Region":         "hn",
				"Environment":    "sim",
				"X-Lane-Cluster": "hna-sim001-v",
			},
			clusters:     clusters,
			destinations: nil,
			err:          "can not find",
		},
	}
	for i, tt := range tests {
		dests, err := p.MatchRoute(tt.domain, tt.tags, tt.clusters)
		if got, want := err, tt.err; !matchError(t, got, want) {
			t.Fatalf("%d: match route, domain=%q: error is not match, got %v, want %v", i, tt.domain, got, want)
		}
		if err != nil {
			continue
		}
		if got, want := dests, tt.destinations; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: match route, domain=%q: destinations is unexpected, got %v, want %v", i, tt.domain, got, want)
		}
	}
}
