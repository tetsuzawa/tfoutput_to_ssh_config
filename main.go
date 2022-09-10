package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"strings"
)

const sshConfigTemplate = `Host bench
    HostName {{ .IP_BENCH }}
    User isucon
    ServerAliveInterval 60
    ForwardAgent yes

Host a
    HostName {{ .IP_A }}
    User isucon
    ServerAliveInterval 60

Host b
    HostName {{ .IP_B }}
    User isucon
    ServerAliveInterval 60

Host c
    HostName {{ .IP_C }}
    User isucon
    ServerAliveInterval 60
`

func main() {
	ipMap := make(map[string]string)
	for _, key := range []string{"IP_BENCH", "IP_A", "IP_B", "IP_C"} {
		// ssh config のエラー回避
		ipMap[key] = "hoge"
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// IP_X = "xxx.xxx.xxx.xxx"
		rawLine := scanner.Text()
		line := strings.Replace(rawLine, " ", "", -1)
		line = strings.Replace(line, "\"", "", -1)
		// IP_X=xxx.xxx.xxx.xxx
		splitLine := strings.Split(line, "=")
		key, ip := splitLine[0], splitLine[1]

		switch key {
		case "IP_BENCH", "IP_A", "IP_B", "IP_C":
			ipMap[key] = ip
		default:
			fmt.Fprintf(os.Stderr, "unkown IP key: %s", key)
		}
	}

	t, err := template.New("ssh_config").Parse(sshConfigTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "template parse failed: %w", err)
		os.Exit(1)
	}

	if err = t.Execute(os.Stdout, ipMap); err != nil {
		fmt.Fprintf(os.Stderr, "template execution failed: %w", err)
		os.Exit(1)
	}
}
