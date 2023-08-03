package model

// ServiceModel 服务模型
type ServiceModel struct {
	Domain   string             // 域名
	Clusters map[string]Cluster // 集群节点表
}

// RouteModel 路由模型
type RouteModel struct {
	Domain   string        // 域名
	Strategy RouteStrategy // 路由策略
}
