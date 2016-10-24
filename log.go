// zheng-ji.info

package main

import (
	"github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
)

func initLogger() error {
	dirPath, _ := filepath.Abs(filepath.Dir(pConfig.Log.Path))
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, 0775)
	}

	file, err := os.OpenFile(pConfig.Log.Path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return err
	}

	level, err := logrus.ParseLevel(pConfig.Log.Level)
	if err != nil {
		return err
	}

	pLog = &logrus.Logger{
		Out:       file,
		Level:     level,
		Formatter: new(logrus.JSONFormatter),
	}

	pLog.Infof("InitLogger: path: %s, level: %s, formatter: json", pConfig.Log.Path, level)

	return nil
}
