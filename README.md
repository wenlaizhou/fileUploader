# fileUploader
文件上传下载助手

本服务无任何依赖, 只有一个可执行文件

使用方式:

拷贝代码, 编译

```
go build -v boot.go
```

启动服务之后, 即创建http上传下载服务, 首页即为文件列表页面

启动命令:

```
./boot 8080 $(pwd)
```

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


拷贝代码如下:


```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	dir := SelfDir()
	port := ":8080"
	switch len(os.Args) {
	case 2:
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			println("使用说明:")
			println("nohup ./uploader $port $path > uploader.log 2>&1 &")
			println("第一个参数为端口号, 默认为8080")
			println("第二个参数为文件目录, 默认为当前目录")
			println("")
			return
		}
		port = fmt.Sprintf(":%v", os.Args[1])
		break
	case 3:
		port = fmt.Sprintf(":%v", os.Args[1])
		dir = fmt.Sprintf("%v", os.Args[2])
		break
	default:
		break
	}
	addrs, _ := net.InterfaceAddrs()
	netAddr := ""
	for _, add := range addrs {
		netAdd := add.String()
		if ipReg.MatchString(netAdd) && !strings.Contains(netAdd, "127.0.0.") {
			res := ipReg.FindAllStringSubmatch(netAdd, -1)
			netAddr = res[0][1]
			break
		}
	}
	fmt.Println("服务目录: ", dir)
	println("启动服务:")
	urlAddr := fmt.Sprintf("http://%v%v", netAddr, port)
	println(urlAddr + "/upload\n")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "*")
		if strings.ToUpper(request.Method) == "OPTIONS" {
			writer.WriteHeader(202)
			return
		}
		http.ServeFile(writer, request, dir+request.URL.Path)
		url := request.URL.Path
		if url[len(url)-1] == '/' {
			// io.WriteString(writer, style)
			io.WriteString(writer, fmt.Sprintf(script, urlAddr, style))
		}
	})

	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(fmt.Sprintf(uploadHtml, urlAddr, style)))
		return
	})

	http.HandleFunc("/doupload", func(writer http.ResponseWriter, request *http.Request) {
		fs, fsHeader, err := request.FormFile("file")
		if err != nil {
			http.Redirect(writer, request, "/", http.StatusFound)
			return
		}
		fileName := fmt.Sprintf("%s/%s", dir, fsHeader.Filename)
		fileContent, err := ioutil.ReadAll(fs)
		if err != nil {
			http.Redirect(writer, request, "/", http.StatusFound)
			return
		}
		_ = ioutil.WriteFile(fileName, fileContent, os.ModePerm)
		http.Redirect(writer, request, "/", http.StatusFound)
		return
	})

	http.HandleFunc("/postUpload", func(writer http.ResponseWriter, request *http.Request) {
		result := map[string]interface{}{}
		filename := request.Header.Get("filename")
		if len(filename) <= 0 {
			result["code"] = -1
			result["message"] = "no filename header"
			res, _ := json.Marshal(result)
			writer.WriteHeader(200)
			writer.Write(res)
			return
		}
		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			result["code"] = -1
			result["message"] = err.Error()
			res, _ := json.Marshal(result)
			writer.WriteHeader(200)
			writer.Write(res)
			return
		}
		ioutil.WriteFile(fmt.Sprintf("%v/%v", dir, filename), data, os.ModePerm)
		result["code"] = 0
		result["message"] = "done"
		result["data"] = filename
		res, _ := json.Marshal(result)
		writer.WriteHeader(200)
		writer.Write(res)
		return
	})

	http.HandleFunc("/getFile", func(writer http.ResponseWriter, request *http.Request) {
		const sniffLen = 512
		filename := request.URL.Query().Get("name")
		if len(filename) <= 0 {
			writer.WriteHeader(404)
			return
		}
		data, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", dir, filename))

		if err != nil || len(data) <= 0 {
			if err != nil {
				println(err.Error())
			}
			writer.WriteHeader(404)
			return
		}

		ctype := mime.TypeByExtension(filepath.Ext(filename))
		if ctype == "" && len(data) > sniffLen {
			// read a chunk to decide between utf-8 text and binary
			ctype = http.DetectContentType(data[:sniffLen])
		}
		writer.Header().Set("Content-Type", ctype)
		writer.Write(data)
		return
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

// 上传文件: curl http://ip:port/doupload -F "file=@文件名"

func SelfPath() string {
	selfPath, _ := filepath.Abs(os.Args[0])
	return selfPath
}

func SelfDir() string {
	return filepath.Dir(SelfPath())
}

var ipReg = regexp.MustCompile("(\\d+\\.\\d+\\.\\d+\\.\\d+)/\\d+")

var selfDir = SelfDir()

const script = `  <body>
    <table class="body-wrap">
      <tr>
        <td></td>
        <td class="container" width="800">
          <div class="content">
            <table class="main" width="100%%" cellpadding="0" cellspacing="0">
              <tr>
                <td class="alert alert-blue">
                  <strong style="font-size: 18px">文件列表</strong>
                </td>
              </tr>
              <tr>
                <td class="content-wrap aligncenter">
                  <table width="100%%" cellpadding="0" cellspacing="0">
                    <!-- <tr>
                      <td class="content-block">
                        <h1>文件列表</h1>
                      </td>
                    </tr> -->
                    <!-- <tr>
                      <td class="content-block">
                        <h2></h2>
                      </td>
                    </tr> -->
                    <tr>
                      <td class="content-block">
                        <table class="invoice">
                          <!-- <tr>
                            <td>
                              Lee Munroe<br />Invoice #12345<br />June 01 2014
                            </td>
                          </tr> -->
                          <tr>
                            <td>
                              <table id="container"
                                class="invoice-items"
                                cellpadding="0"
                                cellspacing="0"
                              >

                                <!-- <tr class="total">
                                  <td class="alignright" width="80%%">共</td>
                                  <td class="alignright">99</td>
                                </tr> -->
                              </table>
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>
                    <tr>
                      <td class="aligncenter content-block">
                        <a class="btn-primary" href="/upload">点击上传</a>
                      </td>
                    </tr>
                    <tr>
                      <td class="content-block content-small">
                        命令行上传: <br/>curl -v %v/doupload -F "file=@文件名"
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <div class="footer">
              <table width="100%%">
                <tr>
                  <td class="aligncenter content-block">
                    Powered By @
                    <a href="http://middleware.cyclone-robotics.com" target="_blank">Middleware Framework</a>
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </td>
        <td></td>
      </tr>
    </table>
  </body>

  <script>
    var preLinks = document.getElementsByTagName("pre")[0]
    preLinks.style.cssText = "display:none"
    var innerLinks = ""
    var fileLinks = document.getElementsByTagName("pre")[0].getElementsByTagName("a")
    for (let i = 0; i < fileLinks.length; i++) {
        const element = fileLinks[i];
        innerLinks += '<tr><td><strong>'+ element.text +'</strong></td><td class="alignright"><a target="_blank" href="' + element.getAttribute("href") +'">点击查看</a></td></tr>'
    }
    document.getElementById("container").innerHTML = innerLinks
    
  </script>
  %v
</html>
`

const uploadHtml = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta name="viewport" content="width=device-width" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Storage Service</title>
  </head>

  <body>
    <table class="body-wrap">
      <tr>
        <td></td>
        <td class="container" width="800">
          <div class="content">
            <table class="main" width="100%%" cellpadding="0" cellspacing="0">
              <tr>
                <td class="alert alert-blue">
                  <strong style="font-size: 18px">Storage Service</strong>
                </td>
              </tr>
              <tr>
                <td class="content-wrap">
                  <table width="100%%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td class="content-block">
                        storage service is a <strong>object storage service</strong>
                        by <strong>middleware framework</strong>
                      </td>
                    </tr>
                    <tr>
                        <td class="aligncenter content-block">
                            <form action="/doupload" method="post" enctype="multipart/form-data">
                                <input class="btn-primary" type="file" name="file">
                                <input type="submit" class="btn-submit" value="开始上传">
                            </form>
                        </td>
                    </tr>
                    <tr>
                      <td class="content-block">
                        <a href="/">查看文件列表</a>
                      </td>
                    </tr>
                    <tr>
                      <td class="content-block content-small">
                        命令行上传: <br/>curl -v %v/doupload -F "file=@文件名"
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <div class="footer">
              <table width="100%%">
                <tr>
                  <td class="aligncenter content-block">
                    Powered By @ <a href="http://middleware.cyclone-robotics.com" target="_blank">Middleware Framework</a>
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </td>
        <td></td>
      </tr>
    </table>
  </body>
%v
</html>`

const style = `
<style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
      font-size: 14px;
    }

    img {
      max-width: 100%;
    }

    body {
      -webkit-font-smoothing: antialiased;
      -webkit-text-size-adjust: none;
      width: 100% !important;
      height: 100%;
      line-height: 1.6;
    }

    table td {
      vertical-align: top;
    }

    body {
      background-color: #f6f6f6;
    }

    .body-wrap {
      background-color: #f6f6f6;
      width: 100%;
    }

    .container {
      display: block !important;
      max-width: 800px !important;
      margin: 0 auto !important;
      /* makes it centered */
      clear: both !important;
    }

    .content {
      max-width: 800px;
      margin: 0 auto;
      display: block;
      padding: 20px;
    }

    .main {
      background: #fff;
      border: 1px solid #e9e9e9;
      border-radius: 10px;
    }

    .content-wrap {
      padding: 20px;
    }

    .content-block {
      padding: 0 0 20px;
    }
	
	.content-small {
      font-size: 12px;
	}

    .header {
      width: 100%;
      margin-bottom: 20px;
    }

    .footer {
      width: 100%;
      clear: both;
      color: #999;
      padding: 20px;
    }
    .footer a {
      color: #999;
    }
    .footer p,
    .footer a,
    .footer unsubscribe,
    .footer td {
      font-size: 12px;
    }

    .column-left {
      float: left;
      width: 50%;
    }

    .column-right {
      float: left;
      width: 50%;
    }

    h1,
    h2,
    h3 {
      color: #000;
      margin: 40px 0 0;
      line-height: 1.2;
      font-weight: 400;
    }

    h1 {
      font-size: 32px;
      font-weight: 500;
    }

    h2 {
      font-size: 24px;
    }

    h3 {
      font-size: 18px;
    }

    h4 {
      font-size: 14px;
      font-weight: 600;
    }

    p,
    ul,
    ol {
      margin-bottom: 10px;
      font-weight: normal;
    }
    p li,
    ul li,
    ol li {
      margin-left: 5px;
      list-style-position: inside;
    }

    a {
      color: #348eda;
      text-decoration: none;
    }

    .btn-primary {
      text-decoration: none;
      color: #fff;
      background-color: #348eda;
      border: solid #348eda;
      border-width: 3px 20px;
      line-height: 2;
      font-weight: bold;
      text-align: center;
      cursor: pointer;
      display: inline-block;
      border-radius: 5px;
      text-transform: capitalize;
    }

    .btn-submit {
      text-decoration: none;
      color: #fff;
      background-color: #d9534f;
      border: solid #d9534f;
      border-width: 3px 20px;
      line-height: 2;
      font-weight: bold;
      text-align: center;
      cursor: pointer;
      display: inline-block;
      border-radius: 5px;
      text-transform: capitalize;
    }

    .last {
      margin-bottom: 0;
    }

    .first {
      margin-top: 0;
    }

    .padding {
      padding: 10px 0;
    }

    .aligncenter {
      text-align: center;
    }

    .alignright {
      text-align: right;
    }

    .alignleft {
      text-align: left;
    }

    .clear {
      clear: both;
    }

    .alert {
      font-size: 16px;
      color: #fff;
      font-weight: 500;
      padding: 20px;
      text-align: center;
      border-radius: 10px 10px 0 0;
    }
    .alert a {
      color: #fff;
      text-decoration: none;
      font-weight: 500;
      font-size: 16px;
    }
    .alert.alert-warning {
      background: #ff9f00;
    }
    .alert.alert-bad {
      background: #d0021b;
    }
    .alert.alert-good {
      background: #68b90f;
    }
    .alert.alert-blue {
      background: #348eda;
    }

    .invoice {
      margin: 40px auto;
      text-align: left;
      width: 90%;
    }
    .invoice td {
      padding: 5px 0;
    }
    .invoice .invoice-items {
      width: 100%;
    }
    .invoice .invoice-items td {
      border-top: #eee 1px solid;
    }
    .invoice .invoice-items .total td {
      border-top: 2px solid #333;
      border-bottom: 2px solid #333;
      font-weight: 700;
    }

    @media only screen and (max-width: 640px) {
      h1,
      h2,
      h3,
      h4 {
        font-weight: 600 !important;
        margin: 20px 0 5px !important;
      }

      h1 {
        font-size: 22px !important;
      }

      h2 {
        font-size: 18px !important;
      }

      h3 {
        font-size: 16px !important;
      }

      .container {
        width: 100% !important;
      }

      .content,
      .content-wrapper {
        padding: 10px !important;
      }

      .invoice {
        width: 100% !important;
      }
    }
  </style>
`

```
