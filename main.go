package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"go.i3wm.org/i3/v4"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "config.toml", "Configuration file")
}

func main() {
	flag.Parse()
	conf := readConfig(configFile)

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
