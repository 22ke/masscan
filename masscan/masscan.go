package masscan

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Masscan struct {
	lua.Super
	cfg     *config
	Command []string
	Result  string

	ctx    context.Context
	cancel context.CancelFunc
}

type ip struct {
	status    string
	protocol  string
	port      string
	ip        string
	timestamp string
}

func newMasscan(cfg *config) *Masscan {
	em := &Masscan{cfg: cfg}
	em.S = lua.INIT
	em.T = MASSCAN
	return em
}

func (m *Masscan) Init() {
	m.cfg.masscanpath = m.cfg.masscanpath + runtime.GOOS+"_masscan.exe"
	m.ctx, m.cancel = context.WithCancel(context.Background())
}

func (mass *Masscan) loop() {
loop:
	tk := time.NewTicker(time.Duration(mass.cfg.Period) * time.Second) //10秒扫描一次
	println("start scanning : ", time.Now().String())
	var cmd *exec.Cmd
	var out, err bytes.Buffer

	mass.Command = append(mass.Command, mass.cfg.Ip)
	mass.Command = append(mass.Command, "-p")
	mass.Command = append(mass.Command, mass.cfg.Port)
	mass.Command = append(mass.Command, "--rate")
	mass.Command = append(mass.Command, mass.cfg.Rate)
	mass.Command = append(mass.Command, "--exclude")
	mass.Command = append(mass.Command, mass.cfg.Exclude)
	mass.Command = append(mass.Command, "--wait")
	mass.Command = append(mass.Command, mass.cfg.Wait)
	mass.Command = append(mass.Command, "-oL")
	mass.Command = append(mass.Command, "-")


	cmd = exec.Command(mass.cfg.masscanpath, mass.Command...)
	fmt.Println("当前运行系统： ", runtime.GOOS)
	fmt.Println("Masscan => ", cmd.Args)
	fmt.Println("Masscan:", cmd)

	cmd.Stdout = &out
	cmd.Stderr = &err

	e := cmd.Run()
	if e != nil {
		println(e.Error())
		if err.Len() > 0 {
			fmt.Printf("masscan run err : %s\n", e.Error())
			println(err.String())
		}
	}
	mass.Result = out.String()
	println("")
	println(mass.Result)
	println("")
	mass.splitstring()

	<-tk.C

	goto loop
}

func (mass *Masscan) splitstring() {
	str := strings.Split(mass.Result, "\r\n") //去掉第一行和最后两行
	len := len(str)
	var ip ip
	for i := 1; i < len-2; i++ {
		//println("i:" , i ,str[i])
		iplist := strings.Split(str[i], " ")
		ip.status = iplist[0]
		ip.protocol = iplist[1]
		ip.port = iplist[2]
		ip.ip = iplist[3]
		//ip.timestamp = iplist[4]
		int64, err := strconv.ParseInt(iplist[4], 10, 64)
		if err != nil {
			println("timeStamp format err ", err.Error())
		}
		ip.timestamp = time.Unix(int64, 0).Format("20060102150405")
		println(ip.ip, ip.port, ip.protocol, ip.timestamp)
		//db.Setmasscandata(ip.ip, ip.port, ip.protocol, ip.timestamp)
	}

}

func (mass *Masscan) Monitor() {
loop:
	tk := time.NewTicker(100 * time.Second)
	//a := db.Getallmasscandata()
	//for i := 0; i < len(a); i++ {
	//	println(a[i].ID)
	//}
	<-tk.C
	goto loop
}

func (m *Masscan) Start() error {
	m.Init()
	go m.loop()
	m.S = lua.RUNNING
	m.U = time.Now()
	logger.Infof("%s masscan start successfully", m.Name())
	return nil
}

func (m *Masscan) Close() error {
	m.S = lua.CLOSE
	if m.cancel != nil {
		m.cancel()
	}
	//close(m.mailChan)
	return nil
}

func (m *Masscan) Name() string {
	return m.cfg.name
}
