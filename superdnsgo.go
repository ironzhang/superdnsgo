package superdnsgo

import (
	"context"
	"fmt"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superdnsgo/pkg/model"
	"github.com/ironzhang/superdnsgo/superdns"
)

var superdnsResolver = &superdns.Resolver{}

// Setup 初始化设置
func Setup(opts Options) (err error) {
	// 初始化选项设置默认值
	if err = opts.setupDefaults(); err != nil {
		tlog.Errorw("options setup defaults", "error", err)
		return fmt.Errorf("options setup defaults: %w", err)
	}

	// 构造服务发现解析程序
	superdnsResolver = &superdns.Resolver{
		Tags:             opts.Tags,
		LoadBalancer:     opts.LoadBalancer,
		SkipPreloadError: opts.Misc.SkipPreloadError,
	}

	// 预加载域名
	err = superdnsResolver.Preload(context.Background(), opts.PreloadDomains)
	if err != nil {
		tlog.Errorw("superdns resolver preload", "domains", opts.PreloadDomains, "error", err)
		return fmt.Errorf("superdns resolver preload: %w", err)
	}
	return nil
}

// AutoSetup 无参初始化
func AutoSetup() error {
	return Setup(Options{})
}

// WithLoadBalancer 构建一个新的服务发现解析程序，并重置负载均衡器
func WithLoadBalancer(lb superdns.LoadBalancer) *superdns.Resolver {
	return superdnsResolver.WithLoadBalancer(lb)
}

// LookupEndpoint 查找地址节点
func LookupEndpoint(ctx context.Context, domain string, tags map[string]string) (endpoint model.Endpoint, cluster string, err error) {
	return superdnsResolver.LookupEndpoint(ctx, domain, tags)
}
