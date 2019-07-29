package main

import "github.com/urfave/cli"

type IOssCli interface {
	Test(c *cli.Context) error
	Add(c *cli.Context) error
	ListFiles(c *cli.Context) error
}
