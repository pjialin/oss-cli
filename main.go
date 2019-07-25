package main

import (
	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

func main() {
	Logger = newLogger()

}
