package AttackModule

import "fmt"

var (
	CommandString string
)

func CommandShell(hostname string) string {
	for {
		fmt.Print("$ ", hostname, " Shell?>>")
		_, _ = fmt.Scan(&CommandString)
		return CommandString
	}
}
