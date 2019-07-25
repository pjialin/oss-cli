package main

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/urfave/cli"
	"sort"
	"strings"
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
	bucket, err := client.Bucket(t.BucketName)
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

func (t *AliOssCli) ListFiles(c *cli.Context) error {
	t.Init(c)
	client, err := oss.New(t.Endpoint, t.Key, t.Secret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(t.BucketName)
	if err != nil {
		return err
	}
	dir := c.Bool("dir")

	var where []oss.Option
	// 前缀筛选
	if prefix := c.String("prefix"); prefix != "" {
		where = append(where, oss.Prefix(prefix))
	}
	// 数量筛选
	if limit := c.Int("limit"); !dir && limit != 0 { // 防止设置 limit 后 dir， 过滤样本不全
		where = append(where, oss.MaxKeys(limit))
	}

	lsRes, err := bucket.ListObjects(where...)
	if err != nil {
		return err
	}
	// 转成可排序的对象列表
	objects := objectList(lsRes.Objects)

	// 排序处理
	if st := c.String("sort"); st != "" {
		if st == "desc" {
			sort.Sort(objectListDesc{objects})
		} else {
			sort.Sort(objects)
		}
	}
	for _, object := range objects {
		if dir && !strings.HasSuffix(object.Key, "/") {
			continue
		}
		fmt.Printf("size: %s	time: %14s	path: %s\n", ByteToShowInConsole(object.Size), object.LastModified.Format("06-01-02 15:04"), object.Key)
	}
	return nil
}

type objectList []oss.ObjectProperties

func (t objectList) Len() int           { return len(t) }
func (t objectList) Less(i, j int) bool { return t[i].LastModified.Before(t[j].LastModified) }
func (t objectList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type objectListDesc struct{ objectList }

func (t objectListDesc) Less(i, j int) bool {
	return t.objectList[i].LastModified.After(t.objectList[j].LastModified)
}
