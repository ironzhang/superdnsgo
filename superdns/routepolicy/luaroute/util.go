package luaroute

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"

	"github.com/ironzhang/superdnsgo/pkg/model"
)

func makeMapTable(ls *lua.LState, m map[string]string) *lua.LTable {
	lt := ls.NewTable()
	for key, value := range m {
		lt.RawSetString(key, lua.LString(value))
	}
	return lt
}

func makeClusterTable(ls *lua.LState, c model.Cluster) *lua.LTable {
	lt := ls.NewTable()
	lt.RawSetString("Name", lua.LString(c.Name))
	lt.RawSetString("Features", makeMapTable(ls, c.Features))
	lt.RawSetString("EndpointNum", lua.LNumber(len(c.Endpoints)))
	return lt
}

func makeClusterMapTable(ls *lua.LState, clusters map[string]model.Cluster) *lua.LTable {
	lt := ls.NewTable()
	for name, cluster := range clusters {
		lt.RawSetString(name, makeClusterTable(ls, cluster))
	}
	return lt
}

func makeDestinations(lt *lua.LTable) ([]model.Destination, error) {
	var dests []model.Destination
	for i := 1; i <= lt.Len(); i++ {
		value := lt.RawGetInt(i)
		table, ok := value.(*lua.LTable)
		if !ok {
			return nil, fmt.Errorf("invalid return value, the %dth element is not a table", i)
		}
		if table.Len() != 2 {
			return nil, fmt.Errorf("invalid return value, the %dth table's len is not 2", i)
		}

		dest := model.Destination{
			Cluster: lua.LVAsString(table.RawGetInt(1)),
			Percent: float64(lua.LVAsNumber(table.RawGetInt(2))),
		}
		dests = append(dests, dest)
	}
	return dests, nil
}
