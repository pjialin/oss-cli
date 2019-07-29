package main

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cheggaaa/pb/v3"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
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
	internal := c.GlobalBool("internal")
	if internal {
		region += "-internal"
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
	Logger.Info("Test result: ok")

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
	search := c.String("search")
	st := c.String("sort")
	limit := 0

	var where []oss.Option
	// 前缀筛选
	if prefix := c.String("prefix"); prefix != "" {
		where = append(where, oss.Prefix(prefix))
	}
	// 数量筛选
	if limit = c.Int("limit"); limit != 0 && !dir && search == "" && st == "" { // 防止设置 limit 后 dir， 过滤样本不全
		where = append(where, oss.MaxKeys(limit))
	}

	lsRes, err := bucket.ListObjects(where...)
	if err != nil {
		return err
	}
	// 转成可排序的对象列表
	objects := objectList(lsRes.Objects)

	// 排序处理
	if st != "" {
		if st == "desc" {
			sort.Sort(objectListDesc{objects})
		} else {
			sort.Sort(objects)
		}
	}
	var result objectList
	for _, object := range objects {
		// Dir 检测
		if dir && !strings.HasSuffix(object.Key, "/") {
			continue
		}
		// 搜索
		if search != "" && !strings.Contains(object.Key, search) {
			continue
		}
		result = append(result, object)
	}
	// limit 检测
	if limit > 0 {
		if len(result) < limit {
			limit = len(result)
		}
		result = result[:limit]
	}
	for _, object := range result {
		fmt.Printf("size: %s	time: %14s	path: %s\n", ByteToShowInConsole(object.Size), object.LastModified.Format("06-01-02 15:04"), object.Key)
	}
	return nil
}

func (t *AliOssCli) Add(c *cli.Context) error {
	t.Init(c)
	client, err := oss.New(t.Endpoint, t.Key, t.Secret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(t.BucketName)
	if err != nil {
		return err
	}
	// 打开文件
	file, bp, save, random := c.String("file"), c.Bool("breakpoint"), c.String("save"), c.Int("random")
	// 文件检查
	fileStat, err := os.Stat(file)
	if err != nil {
		return err
	}
	var files []string
	if fileStat.IsDir() {
		err = filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	} else {
		files = append(files, file)
	}
	// 上传所有文件
	save = strings.TrimRight(save, "/") + "/"
	for _, f := range files {
		var newSave string
		if file != f {
			newSave = save + strings.TrimLeft(strings.Replace(filepath.Dir(f), strings.TrimRight(file, "/"), "", 1), "/")
		} else {
			newSave = save
		}
		Logger.Infof("文件 [%s] 上传中...", f)
		fileName, err := uploadSingleFile(f, newSave, random, bp, bucket)
		if err != nil {
			Logger.Errorf("文件 [%s] 上传失败, %s", f, err.Error())
		} else {
			Logger.Infof("文件 [%s] 上传成功, 保存目录 [%s]", f, fileName)
		}
	}

	return nil
}

func uploadSingleFile(file string, save string, random int, bp bool, bucket *oss.Bucket) (string, error) {
	var saveName string
	var options []oss.Option
	// Save 检查
	save = strings.Trim(save, "/")

	if random == 0 {
		saveName = filepath.Base(file)
	} else {
		saveName = RandomStringInt(random) + filepath.Ext(file)
	}
	// Break point
	if bp {
		options = append(options, oss.Checkpoint(true, file+".cp"))
	}
	options = append(options, oss.Progress(&OssProgressListener{}))
	fileName := save + "/" + saveName
	err := bucket.UploadFile(fileName, file, 100*MB, options...)
	if err != nil {
		return "", nil
	}
	return fileName, nil
}

// 定义进度条监听器。
type OssProgressListener struct {
	Bar *pb.ProgressBar
}

func (t *OssProgressListener) InitBar(size int64) {
	t.Bar = pb.New64(size)
	t.Bar.Set(pb.Bytes, true)
	t.Bar.Start()
}

// 定义进度变更事件处理函数。
func (t *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		if t.Bar == nil {
			t.InitBar(event.TotalBytes)
		}
		if event.ConsumedBytes > 0 {
			t.Bar.SetCurrent(event.ConsumedBytes)
		}
	case oss.TransferDataEvent:
		t.Bar.SetCurrent(t.Bar.Current() + event.RwBytes)
	case oss.TransferCompletedEvent:
		if event.ConsumedBytes >= t.Bar.Total() {
			t.Bar.Finish()
		}
	case oss.TransferFailedEvent:
		Logger.Errorf("传输失败: 已传输大小 %s", ByteToShowNormal(event.ConsumedBytes))
	}
}

// 排序
type objectList []oss.ObjectProperties

func (t objectList) Len() int           { return len(t) }
func (t objectList) Less(i, j int) bool { return t[i].LastModified.Before(t[j].LastModified) }
func (t objectList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type objectListDesc struct{ objectList }

func (t objectListDesc) Less(i, j int) bool {
	return t.objectList[i].LastModified.After(t.objectList[j].LastModified)
}
