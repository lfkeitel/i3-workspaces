package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"go.i3wm.org/i3/v4"
)

var (
	configFile string
	sway       bool
)

func init() {
	flag.StringVar(&configFile, "c", "config.toml", "Configuration file")
	flag.BoolVar(&sway, "sway", false, "Use sway instead of i3")
}

func main() {
	flag.Parse()
	conf := readConfig(configFile)

	if sway {
		i3.SocketPathHook = func() (string, error) {
			sock, exists := os.LookupEnv("SWAYSOCK")
			if !exists {
				return "", fmt.Errorf("SWAYSOCK environment variable not set")
			}
			return sock, nil
		}
	}

	events := i3.Subscribe(i3.WorkspaceEventType)

	log.Println("Subscribed to workspace events")

	for events.Next() {
		event := events.Event().(*i3.WorkspaceEvent)

		log.Printf("Received event: %s\n", event.Change)
		if event.Change != "init" {
			continue
		}

		current := event.Current
		wsName := current.Name
		wsID := current.ID
		wsConf, exists := conf[wsName]

		if !exists || len(current.Nodes) > wsConf.Threshold {
			continue
		}

		log.Printf("Received init for new workspace: %s\n", wsName)

		for _, cmd := range wsConf.Commands {
			i3.RunCommand(fmt.Sprintf("[con_id=%d] %s", wsID, cmd))
		}
	}

	log.Println(events.Close().Error())
}

type Config map[string]struct {
	Threshold int
	Commands  []string
}

func readConfig(filename string) Config {
	var conf Config
	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		log.Fatal(err)
	}
	return conf
}
