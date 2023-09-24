package superdnsgo

import (
	"fmt"

	"github.com/ironzhang/superlib/fileutil"

	"github.com/ironzhang/superdnsgo/superdns"
	"github.com/ironzhang/superdnsgo/superdns/lb"
)

// Misc 杂项信息
type Misc struct {
	// 忽略预加载错误
	SkipPreloadError bool
}

// Options 初始化选项
type Options struct {
	// 分区信息，为 nil 则读取文件 superoptions/zone.json
	Zone *superdns.Zone

	// 杂项信息，为 nil 则读取文件 superoptions/misc.json
	Misc *Misc

	// 预加载域名列表，为 nil 则读取文件 superoptions/preload.json
	PreloadDomains []string

	// 负载均衡器，为 nil 则使用 lb.WRLoadBalancer
	LoadBalancer superdns.LoadBalancer
}

func (p *Options) setupDefaults() (err error) {
	if p.Zone == nil {
		p.Zone, err = readZone()
		if err != nil {
			return fmt.Errorf("read zone: %w", err)
		}
	}

	if p.Misc == nil {
		p.Misc, err = readMisc()
		if err != nil {
			return fmt.Errorf("read misc: %w", err)
		}
	}

	if p.PreloadDomains == nil {
		p.PreloadDomains, err = readPreloadDomains()
		if err != nil {
			return fmt.Errorf("read preload domains: %w", err)
		}
	}

	if p.LoadBalancer == nil {
		p.LoadBalancer = &lb.WRLoadBalancer{}
	}

	return nil
}

func readZone() (*superdns.Zone, error) {
	var zone superdns.Zone
	const path = "superoptions/zone.json"
	if fileutil.FileExist(path) {
		err := fileutil.ReadJSON(path, &zone)
		if err != nil {
			return nil, err
		}
	}
	return &zone, nil
}

func readMisc() (*Misc, error) {
	var misc Misc
	const path = "superoptions/misc.json"
	if fileutil.FileExist(path) {
		err := fileutil.ReadJSON(path, &misc)
		if err != nil {
			return nil, err
		}
	}
	return &misc, nil
}

func readPreloadDomains() (domains []string, err error) {
	const path = "superoptions/preload.json"
	if fileutil.FileExist(path) {
		err = fileutil.ReadJSON(path, &domains)
		if err != nil {
			return nil, err
		}
	}
	return domains, nil
}
