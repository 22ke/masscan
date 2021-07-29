package masscan

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/utils"
)

type config struct {
	name        string
	Ip          string
	Port        string
	Rate        string
	Exclude     string
	Wait        string
	Period      int
	masscanpath string
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}

	tab.ForEach(func(key lua.LValue, val lua.LValue) {
		if key.Type() != lua.LTString {
			L.RaiseError("invalid options %s", key.Type().String())
			return
		}

		switch key.String() {
		case "name":
			cfg.name = val.String()
		case "ip":
			cfg.Ip = val.String()
		case "port":
			cfg.Port = val.String()
		case "rate":
			cfg.Rate = val.String()
		case "exclude":
			cfg.Exclude = val.String()
		case "wait":
			cfg.Wait = val.String()
		case "period":
			cfg.Period = utils.LValueToInt(val, 10)
		case "masscanpath":
			cfg.masscanpath = val.String()

		default:
			L.RaiseError("invalid options %s key", key.String())
		}
	})

	return cfg
}
