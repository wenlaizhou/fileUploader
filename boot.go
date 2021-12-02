package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const uploadUi = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>FileUploader</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    %s
</head>
<body>
	%s
</body>
</html>
`

const uploadHtml = `<h1>Storage Service</h1>
<br>
<h2><a style="color: #888; text-decoration: none;" href="https://github.com/wenlaizhou/fileUploader" target="_blank">
		FileUploader(Common Service)
	</a>
</h2>
<br>
<h3>
	<a style="color: #888; text-decoration: none;" href="/">回到首页</a>
</h3>
<br>
<br>
<form action="/doupload" method="post" enctype="multipart/form-data">
	<input class="fileBtn" style="width: 300px;" type="file" name="file">
	<br>
	<br>
	<input type="submit" style="width: 300px;" class="submit" value="开始上传">
</form>`

const style = `
<style>
	* {
		line-height: 1.2;
		margin: 0;
	}
	html {
		color: #888;
		display: table;
		font-family: sans-serif;
		height: 100%;
		text-align: center;
		width: 100%;
	}
	body {
		display: table-cell;
		vertical-align: middle;
		margin: 2em auto;
	}
	h1 {
		color: #555;
		font-size: 2em;
		font-weight: 400;
	}
	p {
		margin: 0 auto;
		width: 280px;
	}
	@media only screen and (max-width: 280px) {
		body,
		p {
			width: 95%;
		}
		h1 {
			font-size: 1.5em;
			margin: 0 0 0.3em;
		}
	}
	.submit {
		border-radius: 30px;
		color: #fff;
		line-height: 1.5;
		background-color: #d9534f;
		border-color: #d43f3a;
		display: inline-block;
		padding: 6px 12px;
		margin-bottom: 0;
		font-size: 16px;
		font-weight: normal;
		text-align: center;
		white-space: nowrap;
		vertical-align: middle;
		-ms-touch-action: manipulation;
		touch-action: manipulation;
		cursor: pointer;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
		user-select: none;
		background-image: none;
	}
	.fileBtn {
		border-radius: 30px;
		background: #26B99A;
		border: 1px solid #169F85;
		color: #fff;
		line-height: 1.5;
		display: inline-block;
		padding: 6px 12px;
		margin-bottom: 0;
		font-size: 16px;
		font-weight: normal;
		text-align: center;
		white-space: nowrap;
		vertical-align: middle;
		-ms-touch-action: manipulation;
		touch-action: manipulation;
		cursor: pointer;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
		user-select: none;
	}
	a {
		color: black;
		text-decoration: none;
	}
</style>
`

const script = `
<style>
a {
	margin-top: 10px;
	padding-top: 10px;
	line-height: 20px;
}

body {
	padding: 30px;
}

pre {
	margin-top: 10px;
	padding-top: 10px;
}
</style>
<script>
	function back() {
		window.location.href = "/upload";
	}
    var bodyHtml = document.body.innerHTML;
    var uploadLink = document.createElement("input");
	uploadLink.type = "button";
	uploadLink.value = "点击上传文件"
    //uploadLink.href = "/upload";
	uploadLink.onclick = back;
	
    //let newContent = document.createTextNode("点击上传文件");
    //uploadLink.appendChild(newContent);
    document.body.insertBefore(uploadLink, document.body.getElementsByTagName("pre")[0])
</script>`

func SelfPath() string {
	selfPath, _ := filepath.Abs(os.Args[0])
	return selfPath
}

func SelfDir() string {
	return filepath.Dir(SelfPath())
}

var ipReg = regexp.MustCompile("(\\d+\\.\\d+\\.\\d+\\.\\d+)/\\d+")

var selfDir = SelfDir()

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
	fmt.Printf("http://%v%v/upload\n", netAddr, port)
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
			io.WriteString(writer, script)
		}
	})

	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(fmt.Sprintf(uploadUi, style, uploadHtml)))
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
		ioutil.WriteFile(filename, data, os.ModePerm)
		result["code"] = 0
		result["message"] = "done"
		result["data"] = filename
		res, _ := json.Marshal(result)
		writer.WriteHeader(200)
		writer.Write(res)
		return
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

// 上传文件: curl http://ip:port/doupload -F "file=@文件名"
