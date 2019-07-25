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

func TestAliOssCli_ListFiles(t *testing.T) {
	app := makeApp()
	c := &AliOssCli{}
	registerFlags(app)
	app.Commands = []cli.Command{
		{Name: "list", Action: c.ListFiles, Flags: getListFlags()},
	}
	args := os.Args[0:1]
	args = append(args, "list", "--prefix=", "--limit=0", "--dir")
	err := app.Run(args)

	assert.Nil(t, err)
}
