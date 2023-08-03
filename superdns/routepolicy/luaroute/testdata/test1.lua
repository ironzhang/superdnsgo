if MatchFuncs == nil then
	MatchFuncs = {}
end

MatchFuncs["www.test1.com"] = function(tags, clusters)
	-- 生产环境路由匹配规则
	if tags["Environment"] == "product" then
		-- 先判断是否存在同名集群
		name = tags["Cluster"]
		for _, cluster in pairs(clusters) do
			if name == cluster.Name then
				return { {cluster.Name, 1} }
			end
		end

		-- 再判断是否有同 Lidc 的集群
		lidc = tags["Lidc"]
		for _, cluster in pairs(clusters) do
			if lidc == cluster.Features["Lidc"] then
				return { {cluster.Name, 1} }
			end
		end

		-- 再判断是否有同地域的集群
		region = tags["Region"]
		for _, cluster in pairs(clusters) do
			if region == cluster.Features["Region"] then
				return { {cluster.Name, 1} }
			end
		end

	end

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

	-- 默认规则
	cluster = clusters["default@mock"]
	if cluster ~= nil then
		return { {cluster.Name, 1} }
	end
	return {}
end
