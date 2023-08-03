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

func pick(dests []model.Destination) (cluster string, err error) {
	sum := 0.0
	p := rand.Float64()
	for _, dest := range dests {
		sum += dest.Percent
		if p < sum {
			return dest.Cluster, nil
		}
	}
	if len(dests) > 0 {
		return dests[0].Cluster, nil
	}
	return "", ErrInvalidDestinations
}

func mergeZoneTags(srcTags map[string]string, zoneTags map[string]string) map[string]string {
	if len(srcTags) <= 0 {
		return zoneTags
	}
	if len(zoneTags) <= 0 {
		return srcTags
	}

	m := make(map[string]string, len(srcTags)+len(zoneTags))
	for k, v := range zoneTags {
		m[k] = v
	}
	for k, v := range srcTags {
		m[k] = v
	}
	return m
}

type lookuper struct {
	service *model.ServiceModel
	route   *model.RouteModel
	policy  *routepolicy.Policy
}

func (p *lookuper) MatchByRoutePolicy(ctx context.Context, domain string, zone Zone, tags map[string]string) []model.Destination {
	tags = mergeZoneTags(tags, zone.RouteTags)
	dests, err := p.policy.MatchRoute(domain, tags, p.service.Clusters)
	if err != nil {
		tlog.Named("superdns").WithContext(ctx).Warnw("policy match route", "domain", domain, "tags", tags, "error", err)
		return nil
	}
	return dests
}

func (p *lookuper) MatchByRouteRule(ctx context.Context, domain string, zone Zone, tags map[string]string, rule model.RouteRule) []model.Destination {
	var dests []model.Destination
	if p.route.Strategy.EnableScriptRoute || rule.EnableScriptRoute {
		dests = p.MatchByRoutePolicy(ctx, domain, zone, tags)
		if len(dests) > 0 {
			return dests
		}
	}
	dests = rule.LidcDestinations[zone.Lidc]
	if len(dests) > 0 {
		return dests
	}
	dests = rule.RegionDestinations[zone.Region]
	if len(dests) > 0 {
		return dests
	}
	if len(rule.EnvironmentDestinations) > 0 {
		return rule.EnvironmentDestinations
	}
	return p.route.Strategy.DefaultDestinations
}

func (p *lookuper) MatchRoute(ctx context.Context, domain string, zone Zone, tags map[string]string) []model.Destination {
	rule, ok := p.route.Strategy.RouteRules[zone.Environment]
	if ok {
		return p.MatchByRouteRule(ctx, domain, zone, tags, rule)
	}

	if p.route.Strategy.EnableScriptRoute {
		dests := p.MatchByRoutePolicy(ctx, domain, zone, tags)
		if len(dests) > 0 {
			return dests
		}
	}
	return p.route.Strategy.DefaultDestinations
}

func (p *lookuper) Lookup(ctx context.Context, domain string, zone Zone, tags map[string]string) (model.Cluster, error) {
	dests := p.MatchRoute(ctx, domain, zone, tags)
	cname, err := pick(dests)
	if err != nil {
		return model.Cluster{}, fmt.Errorf("%s domain can not match route: %w", domain, err)
	}
	c, ok := p.service.Clusters[cname]
	if !ok {
		return model.Cluster{}, fmt.Errorf("%s domain can not find %s cluster: %w", domain, cname, ErrClusterNotFound)
	}
	return c, nil
}
