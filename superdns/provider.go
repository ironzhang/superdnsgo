package superdns

import (
	"sync/atomic"

	"github.com/ironzhang/superdnsgo/pkg/model"
)

type provider struct {
	domain  string
	service atomic.Value // *model.ServiceModel
	route   atomic.Value // *model.RouteModel
}

func (p *provider) StoreServiceModel(s *model.ServiceModel) {
	p.service.Store(s)
}

func (p *provider) LoadServiceModel() (*model.ServiceModel, bool) {
	s, ok := p.service.Load().(*model.ServiceModel)
	return s, ok
}

func (p *provider) StoreRouteModel(r *model.RouteModel) {
	p.route.Store(r)
}

func (p *provider) LoadRouteModel() (*model.RouteModel, bool) {
	r, ok := p.route.Load().(*model.RouteModel)
	if !ok {
		return &model.RouteModel{Domain: p.domain}, true
	}
	return r, ok
}
