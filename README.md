# OSS Cli 

OSS 工具，目前支持阿里云 OSS

## Features
* 文件展示、搜索
* 文件上传、上传进度、断点续传

## Usage
### 使用 Docker
```
docker run --rm -e ACCESS_KEY=replace -e ACCESS_KEY_SECRET=replace \
    -e BUCKET_NAME=replace -e REGION_NAME=cn-qingdao pjialin/oss-cli:latest test
```

### 命令列表
```
COMMANDS:
     test     测试配置是否正确
     add      上传文件到 OSS 中
     list     查看文件列表
     help, h  Shows a list of commands or help for one command
```

## License
[Apache License 2.0](https://github.com/pjialin/oss-cli/blob/master/LICENSE)