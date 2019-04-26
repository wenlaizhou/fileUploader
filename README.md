# fileUploader
文件上传下载助手

本服务无任何依赖, 只有一个可执行文件

使用方式:

下载可执行文件:

```bash
wget https://github.com/wenlaizhou/fileUploader/raw/master/uploader && chmod +x uploader
```

启动服务之后, 即创建http上传下载服务, 首页即为文件列表页面

上传文件方式:
```bash
curl http://ip:port/doupload -F "file=@文件名"
```

上传文件界面:
```bash
http://ip:port/upload
```

<a href="https://996.icu"><img src="https://img.shields.io/badge/link-996.icu-red.svg" alt="996.icu" /></a>

[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)
