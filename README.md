# fileUploader
文件上传下载助手

使用方式:
复制一下代码, sh执行:

```
wget https://github.com/wenlaizhou/fileUploader/raw/master/boot.go && go build -v
```
./编译之后的可执行文件即可启动服务

启动服务之后, 即创建http上传下载服务, 首页即为文件列表页面

上传文件方式:
```
curl http://ip:port/doupload -F "file=@文件名"
```

上传文件界面:
```
http://ip:port/upload
```

<a href="https://996.icu"><img src="https://img.shields.io/badge/link-996.icu-red.svg" alt="996.icu" /></a>

[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)
