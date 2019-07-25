package main

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/urfave/cli"
)

type AliOssCli struct {
	Key        string
	Secret     string
	BucketName string
	Endpoint   string
}

// 初始化 endpoint
func (t *AliOssCli) Init(c *cli.Context) {
	t.Key, t.Secret = getKeyAndSecret(c)
	t.BucketName = c.GlobalString("bucket")
	var region string
	if region = c.GlobalString("region"); region == "" {
		region = "cn-qingdao"
	}
	t.Endpoint = fmt.Sprintf("http://oss-%s.aliyuncs.com", region)
}

func (t *AliOssCli) Test(c *cli.Context) error {
	t.Init(c)

	failMsg := "Test result: Fail, msg: %s"
	client, err := oss.New(t.Endpoint, t.Key, t.Secret)
	if err != nil {
		return errors.New(fmt.Sprintf(failMsg, err.Error()))
	}
	bucket, err := client.Bucket("p-document")
	if err != nil {
		return errors.New(fmt.Sprintf(failMsg, err.Error()))
	}

	_, err = bucket.ListObjects(oss.MaxKeys(1))
	if err != nil {
		return errors.New(fmt.Sprintf(failMsg, err.Error()))
	}
	Logger.Info("test result: ok")

	return nil
}
