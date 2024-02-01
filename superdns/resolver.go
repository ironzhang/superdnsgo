package superdns

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superlib/fileutil"
	"github.com/ironzhang/superlib/filewatch"
	"github.com/ironzhang/superlib/superutil/parameter"

	"github.com/ironzhang/superdnsgo/pkg/agentclient"
	"github.com/ironzhang/superdnsgo/pkg/model"
	"github.com/ironzhang/superdnsgo/superdns/lb"
	"github.com/ironzhang/superdnsgo/superdns/routepolicy"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	Pickup(domain, cluster string, endpoints []model.Endpoint) (model.Endpoint, error)
}

// Resolver 服务发现解析程序
type Resolver struct {
	Tags             map[string]string // 路由标签
	LoadBalancer     LoadBalancer      // 负载均衡器
	SkipPreloadError bool              // 忽略预加载错误

	once     sync.Once
	resolver *resolver
}

func (r *Resolver) init() {
	if r.resolver != nil {
		return
	}

	tlog.Named("superdns").Debugw("init superdns resolver", "resolver", r, "param", parameter.Param)
	if r.LoadBalancer == nil {
		r.LoadBalancer = &lb.WRLoadBalancer{}
	}
	r.resolver = newResolver(r.Tags, parameter.Param)
}

func (r *Resolver) clone() *Resolver {
	r.once.Do(r.init)

	return &Resolver{
		Tags:             r.Tags,
		LoadBalancer:     r.LoadBalancer,
		SkipPreloadError: r.SkipPreloadError,
		resolver:         r.resolver,
	}
}

// WithLoadBalancer 构建一个新的服务发现解析程序，并重置负载均衡器
func (r *Resolver) WithLoadBalancer(lb LoadBalancer) *Resolver {
	c := r.clone()
	c.LoadBalancer = lb
	return c
}

// Preload 预加载
func (r *Resolver) Preload(ctx context.Context, domains []string) error {
	r.once.Do(r.init)

	// 执行预加载
	err := r.resolver.Preload(ctx, domains)
	if !r.SkipPreloadError && err != nil {
		return err
	}

	// 检查预加载结果
	if !r.SkipPreloadError {
		for _, domain := range domains {
			_, err = r.resolver.LookupCluster(ctx, domain, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// LookupEndpoint 查找地址节点
func (r *Resolver) LookupEndpoint(ctx context.Context, domain string, tags map[string]string) (model.Endpoint, string, error) {
	r.once.Do(r.init)

	c, err := r.resolver.LookupCluster(ctx, domain, tags)
	if err != nil {
		tlog.Named("superdns").WithContext(ctx).Errorw("lookup cluster", "domain", domain, "tags", tags, "error", err)
		return model.Endpoint{}, "", err
	}
	ep, err := r.LoadBalancer.Pickup(domain, c.Name, c.Endpoints)
	if err != nil {
		tlog.Named("superdns").WithContext(ctx).Errorw("load balancer pickup", "domain", domain, "cluster", c.Name, "error", err)
		return model.Endpoint{}, "", err
	}
	return ep, c.Name, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// 内部核心实现
///////////////////////////////////////////////////////////////////////////////////////////////////

// resolver 服务发现解析程序核心实现
type resolver struct {
	tags      map[string]string    // 路由标签
	param     parameter.Parameter  // 解析程序配置参数
	agent     *agentclient.Client  // agent 客户端
	watcher   *filewatch.Watcher   // 文件订阅程序
	policy    *routepolicy.Policy  // 路由策略
	mu        sync.Mutex           // 服务提供方映射表互斥锁
	providers map[string]*provider // 服务提供方映射表，key 为 domain
}

// newResolver 构造服务发现解析程序核心实现
func newResolver(tags map[string]string, param parameter.Parameter) *resolver {
	r := &resolver{
		tags:  tags,
		param: param,
		agent: agentclient.New(agentclient.Options{
			Addr:    param.Agent.Server,
			Timeout: time.Duration(param.Agent.Timeout) * time.Second,
		}),
		watcher:   filewatch.NewWatcher(time.Duration(param.Watch.WatchInterval) * time.Second),
		policy:    routepolicy.NewPolicy(),
		providers: make(map[string]*provider),
	}
	go r.running()
	return r
}

// Preload 预加载
func (r *resolver) Preload(ctx context.Context, domains []string) error {
	err := r.watchProviders(ctx, domains)
	if err != nil {
		tlog.Named("superdns").WithContext(ctx).Errorw("watch providers", "domains", domains, "error", err)
		return err
	}
	return nil
}

// LookupCluster 查找集群节点
func (r *resolver) LookupCluster(ctx context.Context, domain string, tags map[string]string) (model.Cluster, error) {
	// 订阅服务提供方
	p, err := r.watchProvider(ctx, domain)
	if err != nil {
		tlog.Named("superdns").WithContext(ctx).Errorw("watch provider", "domain", domain, "error", err)
		return model.Cluster{}, err
	}

	// 获取服务模型
	service, ok := p.LoadServiceModel()
	if !ok {
		tlog.Named("superdns").WithContext(ctx).Errorw("can not load service model", "domain", domain)
		return model.Cluster{}, r.serviceNotLoad(domain)
	}

	// 获取路由模型
	route, ok := p.LoadRouteModel()
	if !ok {
		tlog.Named("superdns").WithContext(ctx).Errorw("can not load route model", "domain", domain)
		return model.Cluster{}, r.routeNotLoad(domain)
	}

	// 查找集群
	c, err := (&lookuper{
		service: service,
		route:   route,
		policy:  r.policy,
		tags:    r.tags,
	}).Lookup(ctx, domain, tags)
	if err != nil {
		tlog.Named("superdns").WithContext(ctx).Errorw("lookup", "domain", domain, "tags", tags, "error", err)
		return model.Cluster{}, err
	}
	return c, nil
}

func (r *resolver) running() {
	t := time.NewTicker(time.Duration(r.param.Agent.KeepAliveInterval) * time.Second)
	for {
		select {
		case <-t.C:
			r.keepAlive()
		}
	}
}

func (r *resolver) keepAlive() {
	domains := r.listProviders()

	// 向 agent 发送订阅域名请求，以保持订阅的心跳
	err := r.agent.SubscribeDomains(context.Background(), domains, time.Duration(r.param.Agent.SubscribeTTL)*time.Second, true)
	if err != nil {
		tlog.Warnw("keep alive fail", "error", err)
		return
	}
	tlog.Debug("keep alive success")
}

func (r *resolver) listProviders() []string {
	var domains []string
	for domain := range r.providers {
		domains = append(domains, domain)
	}
	return domains
}

func (r *resolver) watchProviders(ctx context.Context, domains []string) error {
	// 向 agent 发送订阅域名请求
	err := r.agent.SubscribeDomains(ctx, domains,
		time.Duration(r.param.Agent.SubscribeTTL)*time.Second, false)
	if err != nil {
		tlog.WithContext(ctx).Warnw("subscribe domains", "domains", domains, "error", err)
		if !r.param.Agent.SkipError {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, domain := range domains {
		r.loadProvider(ctx, domain)
	}

	return nil
}

func (r *resolver) watchProvider(ctx context.Context, domain string) (*provider, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查服务提供方是否已存在
	p, ok := r.providers[domain]
	if ok {
		return p, nil
	}

	// 向 agent 发送订阅域名请求
	err := r.agent.SubscribeDomains(ctx, []string{domain},
		time.Duration(r.param.Agent.SubscribeTTL)*time.Second, false)
	if err != nil {
		tlog.WithContext(ctx).Warnw("subscribe domains", "domain", domain, "error", err)
		if !r.param.Agent.SkipError {
			return nil, err
		}
	}

	// 构建新的服务提供方
	return r.loadProvider(ctx, domain), nil
}

func (r *resolver) loadProvider(ctx context.Context, domain string) *provider {
	// 检查服务提供方是否已存在
	p, ok := r.providers[domain]
	if ok {
		return p
	}

	// 构建新的服务提供方对象
	p = &provider{domain: domain}

	// 订阅服务文件
	r.watcher.WatchFile(ctx, r.serviceFilePath(domain), func(path string) bool {
		var m model.ServiceModel
		err := fileutil.ReadJSON(path, &m)
		if err != nil {
			return false
		}
		p.StoreServiceModel(&m)
		return false
	})

	// 订阅路由文件
	r.watcher.WatchFile(ctx, r.routeFilePath(domain), func(path string) bool {
		var m model.RouteModel
		err := fileutil.ReadJSON(path, &m)
		if err != nil {
			return false
		}
		p.StoreRouteModel(&m)
		return false
	})

	// 订阅路由脚本
	r.watcher.WatchFile(ctx, r.routeScriptPath(domain), func(path string) bool {
		err := r.policy.Load(path)
		if err != nil {
			tlog.Named("superdns").Errorw("policy load", "path", path, "error", err)
		}
		return false
	})

	r.providers[domain] = p

	return p
}

func (r *resolver) serviceFilePath(domain string) string {
	filename := fmt.Sprintf("%s.json", domain)
	return filepath.Join(r.param.Watch.ResourcePath, "services", filename)
}

func (r *resolver) routeFilePath(domain string) string {
	filename := fmt.Sprintf("%s.json", domain)
	return filepath.Join(r.param.Watch.ResourcePath, "routes", filename)
}

func (r *resolver) routeScriptPath(domain string) string {
	filename := fmt.Sprintf("%s.lua", domain)
	return filepath.Join(r.param.Watch.ResourcePath, "routes", filename)
}

func (r *resolver) serviceNotLoad(domain string) error {
	return fmt.Errorf("can not load %s file", r.serviceFilePath(domain))
}

func (r *resolver) routeNotLoad(domain string) error {
	return fmt.Errorf("can not load %s file", r.routeFilePath(domain))
}
