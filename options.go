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
	// 路由标签，为 nil 则读取文件 superoptions/tags.json
	Tags map[string]string

	// 预加载域名列表，为 nil 则读取文件 superoptions/preload.json
	PreloadDomains []string

	// 负载均衡器，为 nil 则使用 lb.WRLoadBalancer
	LoadBalancer superdns.LoadBalancer

	// 杂项信息，为 nil 则读取文件 superoptions/misc.json
	Misc *Misc
}

func (p *Options) setupDefaults() (err error) {
	if p.Tags == nil {
		p.Tags, err = readTags()
		if err != nil {
			return fmt.Errorf("read tags: %w", err)
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

	if p.Misc == nil {
		p.Misc, err = readMisc()
		if err != nil {
			return fmt.Errorf("read misc: %w", err)
		}
	}

	return nil
}

func readTags() (map[string]string, error) {
	var tags map[string]string
	const path = "superoptions/tags.json"
	if fileutil.FileExist(path) {
		err := fileutil.ReadJSON(path, &tags)
		if err != nil {
			return nil, err
		}
	}
	return tags, nil
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
