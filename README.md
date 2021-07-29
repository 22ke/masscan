# masscan

```
local masscan = rock.masscan {
    ip = "10.0.0.5,192.168.193.133,172.31.61.10-172.31.61.202,61.152.230.35",
    port = "81,8080,8090,8001,22,25,58,80",
    rate = 10000,
    exclude = "10.0.0.33,10.0.0.35-10.0.0.38",
    wait = 3,
    period = 10,
    masscanpath = "masscan\\resource\\software\\",
}

masscan.start()
```
