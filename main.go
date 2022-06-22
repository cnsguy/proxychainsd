package main

import (
	"fmt"
	"os"
	"sync"
)

func MakeConfig(fname string) (*Config, error) {
	conf := DefaultConfig()
	return conf, WriteConfig(fname, conf)
}

func readConfigWrap(fname string, createNew bool) (*Config, error) {
	config, err := ReadConfig(fname)

	if os.IsNotExist(err) && createNew {
		fmt.Printf("Config file '%s' doesn't exist yet, creating\n", fname)
		newConfig, err := MakeConfig(fname)

		if err != nil {
			return nil, fmt.Errorf("Failed to create config: %s", err.Error())
		} else {
			return newConfig, nil
		}
	} else if err != nil {
		return nil, fmt.Errorf("Failed to read config: %s", err)
	}

	return config, nil
}

func main() {
	config, err := readConfigWrap("config.json", true) // XXX parse cli args

	if err != nil {
		fmt.Printf("Failed to read config: %s\n", err.Error())
		return
	}

	var wg sync.WaitGroup

	for _, listenConf := range config.Listeners {
		server := NewServer(&listenConf, config.Chains)
		go server.Run(wg)
		wg.Add(1)
	}

	wg.Wait()
}
