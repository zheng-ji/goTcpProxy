// Author: zheng-ji.info

package main

import (
	//"errors"
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
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

func parseConfigFile(filepath string) error {
	if config, err := ioutil.ReadFile(filepath); err == nil {
		if err = goyaml.Unmarshal(config, &pConfig); err != nil {
			return err
		}
		fmt.Println(pConfig)
	} else {
		return err
	}
	return nil
}
