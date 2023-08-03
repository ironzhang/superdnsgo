package luaroute

import (
	"errors"
	"fmt"

	lua "github.com/yuin/gopher-lua"

	"github.com/ironzhang/superdnsgo/pkg/model"
)

// Policy 路由策略
type Policy struct {
	lstate *lua.LState
}

// NewPolicy 构建路由策略
func NewPolicy() *Policy {
	return &Policy{lstate: lua.NewState()}
}

// Load 加载路由脚本
func (p *Policy) Load(path string) error {
	err := p.lstate.DoFile(path)
	if err != nil {
		return fmt.Errorf("luaroute: do file %q: %w", path, err)
	}
	return nil
}

// MatchRoute 执行路由匹配
func (p *Policy) MatchRoute(domain string, tags map[string]string, clusters map[string]model.Cluster) ([]model.Destination, error) {
	// 查找脚本函数
	fn, err := p.lookupFunction("MatchFuncs", domain)
	if err != nil {
		return nil, fmt.Errorf("luaroute: lookup function: %w", err)
	}

	// 调用脚本函数
	err = p.lstate.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, makeMapTable(p.lstate, tags), makeClusterMapTable(p.lstate, clusters))
	if err != nil {
		return nil, fmt.Errorf("luaroute: call by param: %w", err)
	}

	// 获取脚本函数调用返回值
	ret := p.lstate.Get(-1)
	p.lstate.Pop(1)
	lt, ok := ret.(*lua.LTable)
	if !ok {
		return nil, errors.New("luaroute: pop value is not a table")
	}

	// 构建目标节点
	dests, err := makeDestinations(lt)
	if err != nil {
		return nil, fmt.Errorf("luaroute: make destinations: %w", err)
	}
	return dests, nil
}

func (p *Policy) lookupFunction(table, key string) (*lua.LFunction, error) {
	lt, ok := p.lstate.GetGlobal(table).(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("%q is not a table", table)
	}
	lf, ok := lt.RawGetString(key).(*lua.LFunction)
	if !ok {
		return nil, fmt.Errorf("can not find function %s[%q]", table, key)
	}
	return lf, nil
}
