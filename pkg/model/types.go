package model

// State 地址状态
type State string

// 地址状态常量定义
const (
	Enabled  State = "enabled"
	Disabled State = "disabled"
)

// Endpoint 地址节点
type Endpoint struct {
	Addr   string // 地址，IP:Port
	State  State  // 状态
	Weight int    // 权重
}

// Cluster 集群节点
type Cluster struct {
	Name                 string            // 集群名
	Features             map[string]string // 集群特征
	Endpoints            []Endpoint        // 地址节点列表
	AvailableEndpointNum int               // 有效地址节点数
}

// Destination 目标节点
type Destination struct {
	Cluster string  // 目标集群名
	Percent float64 // 流量配比
}

// RouteRule 路由规则
type RouteRule struct {
	EnableScriptRoute       bool                     // 是否启用脚本路由
	LidcDestinations        map[string][]Destination // 机房级默认路由目标
	RegionDestinations      map[string][]Destination // 地域级默认路由目标
	EnvironmentDestinations []Destination            // 环境级默认路由目标
}

// RouteStrategy 路由策略
type RouteStrategy struct {
	EnableScriptRoute   bool                 // 是否启用脚本路由
	RouteRules          map[string]RouteRule // 路由规则，key 为环境名称
	DefaultDestinations []Destination        // 默认路由目标
}
