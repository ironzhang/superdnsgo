package main

var luascript = `-- www.superdns.com
if MatchFuncs == nil then
	MatchFuncs = {}
end

MatchFuncs["www.superdns.com"] = function(tags, clusters)
	-- 仿真环境路由匹配规则
	if tags["Environment"] == "sim" then
		target = tags["X-Lane-Cluster"] -- 泳道集群
		cluster = clusters[target]
		if cluster ~= nil then
			return { {cluster.Name, 1} }
		end

		cluster = clusters["hna-sim000-v"] -- 基准集群
		if cluster ~= nil then
			return { {cluster.Name, 1} }
		end
	end

	return {}
end
`
