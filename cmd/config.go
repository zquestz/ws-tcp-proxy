package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/zquestz/go-ucl"
)

const (
	appName     = "ws-tcp-proxy"
	version     = "0.1.0"
	defaultPort = 8080
)

// Config stores all the application configuration.
type Config struct {
	Cert           string `json:"cert"`
	DisplayVersion bool   `json:"-"`
	Key            string `json:"key"`
	Port           int    `json:"port,string"`
	TextMode       bool   `json:"textMode,string"`
	TCPTLS         bool   `json:"tcpTLS,string"`
	TCPTLSCert     string `json:"tcpTLSCert"`
	TCPTLSKey      string `json:"tcpTLSKey"`
	TCPTLSRootCA   string `json:"tcpTLSRootCA"`
	AutoCert       string `json:"autoCert"`
}

// Load reads the configuration from ~/.config/ws-tcp-proxy/config
// and loads it into the Config struct.
// The config is in UCL format.
func (c *Config) Load() error {
	conf, err := c.loadConfig()
	if err != nil {
		return err
	}

	// There are cases when we don't have a configuration.
	if conf != nil {
		err = c.applyConf(conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) configDir() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(h, ".config", appName), nil
}

func (c *Config) loadConfig() ([]byte, error) {
	configDir, err := c.configDir()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(configDir, "config"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	ucl.Ucldebug = false
	data, err := ucl.NewParser(f).Ucl()
	if err != nil {
		return nil, err
	}

	conf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) applyConf(conf []byte) error {
	err := json.Unmarshal(conf, c)
	if err != nil {
		return err
	}

	return nil
}
