package model

import "github.com/ironzhang/superlib/superutil/supermodel"

type State = supermodel.State
type Endpoint = supermodel.Endpoint
type Cluster = supermodel.Cluster
type Destination = supermodel.Destination
type RouteRule = supermodel.RouteRule
type RouteStrategy = supermodel.RouteStrategy
type ServiceModel = supermodel.ServiceModel
type RouteModel = supermodel.RouteModel

const (
	Enabled  = supermodel.Enabled
	Disabled = supermodel.Disabled
)
