package main

import "github.com/urfave/cli"

type IOssCli interface {
	Test(c *cli.Context) error
	//Add()
	//ListFiles() interface{}
}
