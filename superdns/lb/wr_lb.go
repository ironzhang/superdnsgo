package lb

import (
	"math/rand"

	"github.com/ironzhang/superdnsgo/pkg/model"
)

// WRLoadBalancer 权重随机均衡算法
type WRLoadBalancer struct {
}

// Pickup 按权重随机挑选地址节点
func (p *WRLoadBalancer) Pickup(domain, cluster string, endpoints []model.Endpoint) (model.Endpoint, error) {
	if len(endpoints) <= 0 {
		return model.Endpoint{}, ErrInvalidEndpoints
	}

	// 按权重随机挑选节点
	total, availables := calcAvailableEndpoints(endpoints)
	if total > 0 && len(availables) > 0 {
		n := 0
		random := rand.Intn(total)
		for _, ep := range availables {
			n += ep.Weight
			if random < n {
				return ep, nil
			}
		}
	}

	// 随机挑选一个节点
	i := rand.Intn(len(endpoints))
	return endpoints[i], nil
}

func calcAvailableEndpoints(endpoints []model.Endpoint) (int, []model.Endpoint) {
	total := 0
	results := make([]model.Endpoint, 0, len(endpoints))
	for _, ep := range endpoints {
		if ep.State == model.Enabled {
			total += ep.Weight
			results = append(results, ep)
		}
	}
	return total, results
}
