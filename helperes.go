package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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

func getKeyAndSecret(c *cli.Context) (string, string) {
	key := c.GlobalString("key")
	secret := c.GlobalString("secret")
	return key, secret
}

func ByteToShowInConsole(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%5dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%5.1f%c", float64(b)/float64(div), "KMGTPE"[exp])
}
