if MatchFuncs == nil then
	MatchFuncs = {}
end

MatchFuncs["www.test2.com"] = function(tags, clusters)
	-- 生产环境路由匹配规则
	if tags["Environment"] == "product" and tags["Region"] == "hn" then
		return { {"hna-v", 0.5}, {"hnb-v", 0.5} }
	end

	-- 默认规则
	cluster = clusters["default@mock"]
	if cluster ~= nil then
		return { {cluster.Name, 1} }
	end
	return {}
end
