## fastdfs-client

[![Build Status](https://www.travis-ci.org/wodog/fastdfs-client.svg?branch=master)](https://www.travis-ci.org/wodog/fastdfs-client)
[![GoDoc](https://godoc.org/github.com/wodog/fastdfs-client?status.svg)](https://godoc.org/github.com/wodog/fastdfs-client)

go版的fastdfs客户端

#### Install

```
go get github.com/laoqiuChina/fastdfs-client
```

#### Usage

```
package main

import "github.com/laoqiuChina/fastdfs-client"

func main() {
  client := fstdfs.New()
	client.AddTracker("my.fastdfs.com:22122")

  // upload file
  file, _ := os.Open(fileName)
  fileId := client.Upload(file,"jpg")

  // download file
  client.Download(fileId, os.Stdout)
}
```

#### Reference

[协议参考](http://weakyon.com/2014/09/01/analysis-of-source-code-for-fastdfs.html)  
[协议参考](http://bbs.chinaunix.net/thread-2001015-1-1.html)  
[nodejs客户端](https://github.com/ymyang/fdfs)
