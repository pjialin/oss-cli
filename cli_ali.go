package main

import "github.com/urfave/cli"

type AliOssCli struct{}

func (t *AliOssCli) Test(c *cli.Context) error {
	key, secret := getKeyAndSecret(c)
	Logger.Println("test", key, secret)

	return nil
}
