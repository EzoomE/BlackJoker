package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func handleHelp() {
	fmt.Println("\n####################################\n&& SessionList > 获取在线Bot\n&& Shell {BotName} > 进入某Bot的Shell\n####################################\n")
}

func handleSessionList() {
	mu.Lock()
	if ObtainingOnlineNumberOfPeople == nil {
		fmt.Print("空")
	}
	fmt.Println(ObtainingOnlineNumberOfPeople)
	mu.Unlock()
}

func handleExit() {
	os.Exit(0)
}

func handleShell(input string) {
	FlagBool = false
	parts := strings.Fields(input)
	if len(parts) >= 2 {
		hostname := parts[1]
		if findSession(hostname) {
			fmt.Println("Before sending to flag channel") // 添加这行打印语句
			flag <- hostname
			fmt.Println("After sending to flag channel") // 添加这行打印语句
			CommandShell(hostname)
		}
	} else {
		fmt.Println("请输入BotName")
	}
}

func handleDefault() {
	fmt.Println("找不到该命令")
}

func UserInputLoop() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case commandeering := <-command:
			if commandeering == "exit" {
			}
		default:
		}
		fmt.Print("$ JokerOSshell?>> ")
		if scanner.Scan() {
			input := scanner.Text()

			switch {
			case input == "help":
				handleHelp()
			case input == "SessionList":
				handleSessionList()
			case input == "exit":
				handleExit()
			case strings.HasPrefix(input, "Shell"):
				handleShell(input)
			default:
				handleDefault()
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
		time.Sleep(time.Second * 2)
		fmt.Printf("$ %s Shell?>> ", hostname)
		if scanner.Scan() {
			commandString := scanner.Text()

			if commandString == "exit" {
				return
			}
			fmt.Println("After sending to command channel") // 添加这行打印语句
			select {
			case command <- commandString:
				FlagBool = true
			default:
				fmt.Println("Command channel is full, unable to send data")
			}
		} else {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading standard input:", err)
	}
}

func InteractiveMain() {
	time.Sleep(time.Second * 2)
	banner()
	UserInputLoop()
}
