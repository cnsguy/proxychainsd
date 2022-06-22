package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ErrConfigInvalidPort uint16
type ErrConfigParseFailed string
type ErrConfigInvalidChainLogic struct {
	Logic string
}

func (e ErrConfigInvalidPort) Error() string {
	return fmt.Sprintf("Invalid port in config: %d\n", e)
}

func (e ErrConfigInvalidChainLogic) Error() string {
	return fmt.Sprintf("Invalid chain logic in config: %s\n", e.Logic)
}

type ListenerConfig struct {
	BindIP       string `json:"bind_ip,validate:some_fun"`
	Port         uint16 `json:"port"`
	EnableSocks4 bool   `json:"enable_socks_4"`
	EnableSocks5 bool   `json:"enable_socks_5"`
}

type ChainConfigProxyEntry struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

type ChainConfig struct {
	Logic   string                  `json:"logic"`
	Entries []ChainConfigProxyEntry `json:"entries"`
}

type Config struct {
	Listeners []ListenerConfig `json:"listeners"`
	Chains    []ChainConfig    `json:"chains"`
}

func ValidLogic() []string {
	return []string{"random", "sequence"}
}

func DefaultConfig() *Config {
	return &Config{
		[]ListenerConfig{
			{
				"127.0.0.1",
				4242,
				true,
				true,
			},
		},
		[]ChainConfig{
			{
				Logic: "random",
				Entries: []ChainConfigProxyEntry{
					{
						"127.0.0.1",
						9050,
					},
				},
			},
		},
	}
}

func ValidateListenerPort(port uint16) error {
	if port == 0 {
		return ErrConfigInvalidPort(port)
	} else {
		return nil
	}
}

func ValidateChainLogic(logic string) error {
	if logic != "random" && logic != "sequence" {
		return ErrConfigInvalidChainLogic{logic}
	} else {
		return nil
	}
}

func ValidateListeners(conf *Config) error {
	for _, listener := range conf.Listeners {
		// XXX validate bind ip

		err := ValidateListenerPort(listener.Port)

		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateChains(conf *Config) error {
	for _, chain := range conf.Chains {
		err := ValidateChainLogic(chain.Logic)

		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateConfig(conf *Config) error {
	err := ValidateListeners(conf)

	if err != nil {
		return err
	} else {
		return ValidateChains(conf)
	}
}

func WriteConfig(fname string, conf *Config) error {
	bytes, err := json.MarshalIndent(conf, "", "    ")

	if err != nil {
		return err
	}

	file, err := os.Create(fname)

	if err != nil {
		return err
	} else {
		_, err = file.Write(bytes)
		return err
	}
}

func ReadConfig(fname string) (*Config, error) {
	file, err := os.Open(fname)

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return nil, err
	} else {
		return config, nil
	}
}
