package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"os"
	"testing"
)

func TestAliOssCli_Test(t *testing.T) {
	app := makeApp()
	c := &AliOssCli{}
	registerFlags(app)
	app.Commands = []cli.Command{
		{Name: "test", Action: c.Test},
	}
	args := os.Args[0:1]
	args = append(args, "test")
	err := app.Run(args)

	assert.Nil(t, err)
}
