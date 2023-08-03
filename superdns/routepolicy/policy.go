package routepolicy

import (
	"sync"

	"github.com/ironzhang/superdnsgo/pkg/model"
	"github.com/ironzhang/superdnsgo/superdns/routepolicy/luaroute"
)

// Policy 路由策略
type Policy struct {
	mu     sync.Mutex
	policy *luaroute.Policy
}

// NewPolicy 构建路由策略
func NewPolicy() *Policy {
	return &Policy{
		policy: luaroute.NewPolicy(),
	}
}

// Load 加载路由脚本
func (p *Policy) Load(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.policy.Load(path)
}

// MatchRoute 执行路由匹配
func (p *Policy) MatchRoute(domain string, tags map[string]string, clusters map[string]model.Cluster) ([]model.Destination, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.policy.MatchRoute(domain, tags, clusters)
}
