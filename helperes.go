package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	const unit = KB
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

func ByteToShowNormal(b int64) string {
	const unit = KB
	if b < unit {
		return fmt.Sprintf("%d KB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// Random
func RandString(l int) string {
	return string(RandomWith(l, RandomRule.String))
}

func RandomStringInt(l int) string {
	return string(RandomWith(l, RandomRule.Int))
}

func RandomInt(l int) int {
	res, _ := strconv.Atoi(RandomStringInt(l))
	return res
}

type randomRule struct {
	String string
	Int    string
}

var RandomRule = randomRule{
	String: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	Int:    "0123456789",
}

func RandomWith(l int, rule string) []byte {
	buf := make([]byte, l)
	for i := range buf {
		buf[i] = rule[rand.Intn(len(rule))]
	}
	return buf
}
