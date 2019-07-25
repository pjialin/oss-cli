package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

func env(key string, defV ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		if len(defV) == 0 {
			return ""
		}
		return defV[0]
	}
	return v
}

func newLogger() *logrus.Logger {
	client := logrus.New()
	return client
}
