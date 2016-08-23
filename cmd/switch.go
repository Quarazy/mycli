package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const gitConfigPath = ".git/config"

var GitConfigCmd = &NoArgCommand{
	run: gitConfigRun,
}

// gitProperty is a configuration setting from gitconfig
type gitProperty struct {
	root, key, value string
}

func (p gitProperty) String() string {
	return p.root + p.key
}

func gitConfigRun() {
	f, err := os.Open(gitConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	conf := make(map[string]gitProperty)
	var order []string
	var rootConf string
	var line string

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line = scanner.Text()
		if len(line) == 0 {
			continue
		} else if line[0:1] == "[" {
			rootConf = line[1 : len(line)-1]
		} else {
			i := strings.Index(line, "=")

			prop := gitProperty{
				rootConf,
				strings.TrimSpace(line[:i]),
				strings.TrimSpace(line[i+1:]),
			}

			conf[prop.String()] = prop
			order = append(order, prop.String())
		}
	}

	conf, order = merge(conf, order)
	writeConfig(conf, order)
}

func merge(
	conf map[string]gitProperty,
	order []string) (map[string]gitProperty, []string) {
	newProps := []gitProperty{
		gitProperty{
			"user",
			"email",
			"quarazy@gmail.com",
		},
		gitProperty{
			"user",
			"name",
			"Quarazy",
		},
	}

	// Add new props if it doesn't exist yet
	for _, prop := range newProps {
		if _, ok := conf[prop.String()]; !ok {
			conf[prop.String()] = prop
			order = append(order, prop.String())
		}
	}

	remoteUrl := `remote "origin"url`
	if remoteUrlProp, ok := conf[remoteUrl]; ok {
		i := strings.LastIndex(remoteUrlProp.value, "/")
		remoteUrlProp.value = "git@github.com-quarazy:Quarazy/" + remoteUrlProp.value[i+1:]
		conf[remoteUrl] = remoteUrlProp
	}

	return conf, order
}

// writeConfig creates a new git config file with the updated configuration
func writeConfig(conf map[string]gitProperty, order []string) {
	f, err := os.Create(gitConfigPath)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	var rootConf string
	for _, o := range order {
		prop := conf[o]
		if prop.root != rootConf {
			rootConf = prop.root
			fmt.Fprintf(writer, "[%s]\n", prop.root)
		}
		fmt.Fprintf(writer, "\t%s = %s\n", prop.key, prop.value)
	}

	writer.Flush()
}
