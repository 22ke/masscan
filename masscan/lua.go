package masscan

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
	"reflect"
)

var (
	MASSCAN = reflect.TypeOf((*Masscan)(nil)).String()
)

func (m *Masscan) Index(L *lua.LState, key string) lua.LValue {
	if key == "start" {
		return lua.NewFunction(m.start)
	}
	if key == "close" {
		return lua.NewFunction(m.close)
	}

	return lua.LNil
}

func (m *Masscan) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "name":
		m.cfg.name = lua.CheckString(L, val)
	case "ip":
		m.cfg.Ip = lua.CheckString(L, val)
	case "port":
		m.cfg.Port = lua.CheckString(L, val)
	case "rate":
		m.cfg.Rate = lua.CheckString(L, val)
	case "exclude":
		m.cfg.Exclude = lua.CheckString(L, val)
	case "wait":
		m.cfg.Wait = lua.CheckString(L, val)
	case "period":
		m.cfg.Period = lua.CheckInt(L, val)
	}
}

func (m *Masscan) start(L *lua.LState) int {
	if m.State() == lua.RUNNING {
		L.RaiseError("%s Masscan is already running", m.cfg.name)
		return 0
	}

	err := m.Start()
	if err != nil {
		L.RaiseError("Masscan sender start error: %v", err)
	}

	return 0
}

func (m *Masscan) close(L *lua.LState) int {
	if m.S == lua.CLOSE {
		L.RaiseError("%s Masscan is already closed", m.cfg.name)
		return 0
	}

	err := m.Close()
	if err != nil {
		L.RaiseError("Masscan close error: %v", err)
	}
	return 0
}

func newLuaMasscan(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name , MASSCAN)
	if proc.IsNil() {
		proc.Set(newMasscan(cfg))
		goto done
	}
	proc.Value.(*Masscan).cfg = cfg

done:
	L.Push(proc)
	return 1
}

func LuaInjectApi(env xcall.Env) {
	env.Set("masscan", lua.NewFunction(newLuaMasscan))
}