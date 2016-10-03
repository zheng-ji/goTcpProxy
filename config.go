// Author: zheng-ji.info

package main

import (
	"errors"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
)

type ProxyConfig struct {
	Bind         string    `yaml:"bind"`
	WaitQueueLen int       `yaml:"wait_queue_len"`
	MaxConn      int       `yaml:"max_conn"`
	Timeout      int       `yaml:timeout`
	FailOver     int       `yaml:failover`
	Backend      []string  `yaml:"backend"`
	Log          LogConfig `yaml:"log"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

func (lc *LogConfig) isValid() bool {
	_, err := os.Stat(lc.Path)
	return err == nil || os.IsExist(err)
}

func parseConfigFile(filepath string) error {
	if config, err := ioutil.ReadFile(filepath); err == nil {
		if err = goyaml.Unmarshal(config, &pConfig); err != nil {
			return err
		}
		if pConfig.isValid() == false {
			err := errors.New("wrong config")
			return err
		}
		if pConfig.FailOver <= 0 {
			pConfig.FailOver = 5
		}
		if pConfig.MaxConn <= 0 {
			pConfig.MaxConn = 5
		}
		if pConfig.Timeout <= 0 {
			pConfig.Timeout = 5
		}
	} else {
		return err
	}
	return nil
}

func (pc *ProxyConfig) isValid() bool {
	return len(pc.Backend) > 0 && pc.Log.isValid()
}
