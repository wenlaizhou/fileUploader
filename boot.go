package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var uploadUi = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">

    <title>Page Not Found</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
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
    </style>
</head>

<body>
<h1 style="color: #337ab7">上传文件</h1>
<br/><br/><br/><br/><br/>
<form action="/doupload" method="post" enctype="multipart/form-data">
    <input class="fileBtn" style="width: 300px;" type="file" name="file"/>
    <br/><br/><br/><br/><br/>
    <input type="submit" style="width: 300px;" class="submit" value="开始上传"/>
</form>

</body>

</html>
`

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
	println("http://" + netAddr + port)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "*")
		if strings.ToUpper(request.Method) == "OPTIONS" {
			writer.WriteHeader(202)
			return
		}
		http.ServeFile(writer, request, dir+request.URL.Path)
	})

	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(uploadUi))
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

	log.Fatal(http.ListenAndServe(port, nil))
}

// 上传文件: curl http://ip:port/doupload -F "file=@文件名"
