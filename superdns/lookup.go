package superdns

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superdnsgo/pkg/model"
	"github.com/ironzhang/superdnsgo/superdns/routepolicy"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func mergeTags(srcTags map[string]string, dstTags map[string]string) map[string]string {
	if len(srcTags) <= 0 {
		return dstTags
	}
	if len(dstTags) <= 0 {
		return srcTags
	}

	m := make(map[string]string, len(srcTags)+len(dstTags))
	for k, v := range srcTags {
		m[k] = v
	}
	for k, v := range dstTags {
		m[k] = v
	}
	return m
}

type lookuper struct {
	service *model.ServiceModel
	route   *model.RouteModel
	policy  *routepolicy.Policy
	tags    map[string]string
}

func (p *lookuper) MakeClusterMap() map[string]model.Cluster {
	m := make(map[string]model.Cluster, len(p.service.Clusters))
	for _, c := range p.service.Clusters {
		m[c.Name] = c
	}
	return m
}

func (p *lookuper) MatchRouteScript(ctx context.Context, domain string, tags map[string]string, clusters map[string]model.Cluster) []model.Destination {
	tags = mergeTags(p.tags, tags)
	dests, err := p.policy.MatchRoute(domain, tags, clusters)
	if err != nil {
		tlog.Named("superdns").WithContext(ctx).Warnw("policy match route", "domain", domain, "tags", tags, "error", err)
		return nil
	}
	return dests
}

func (p *lookuper) MatchRoute(ctx context.Context, domain string, tags map[string]string, clusters map[string]model.Cluster) []model.Destination {
	if p.route.Strategy.EnableScript {
		dests := p.MatchRouteScript(ctx, domain, tags, clusters)
		if len(dests) > 0 {
			return dests
		}
	}
	return p.route.Strategy.DefaultDestinations
}

func (p *lookuper) Pick(dests []model.Destination) (cluster string, err error) {
	sum := 0.0
	r := rand.Float64()
	for _, dest := range dests {
		sum += dest.Percent
		if r < sum {
			return dest.Cluster, nil
		}
	}
	if len(dests) > 0 {
		return dests[0].Cluster, nil
	}
	if len(p.service.Clusters) > 0 {
		return p.service.Clusters[0].Name, nil
	}
	return "", ErrNoAvalibaleCluster
}

func (p *lookuper) Lookup(ctx context.Context, domain string, tags map[string]string) (model.Cluster, error) {
	clusters := p.MakeClusterMap()
	dests := p.MatchRoute(ctx, domain, tags, clusters)
	cname, err := p.Pick(dests)
	if err != nil {
		return model.Cluster{}, fmt.Errorf("%s domain can not pick cluster: %w", domain, err)
	}
	c, ok := clusters[cname]
	if !ok {
		return model.Cluster{}, fmt.Errorf("%s domain can not find %s cluster: %w", domain, cname, ErrClusterNotFound)
	}
	return c, nil
}
