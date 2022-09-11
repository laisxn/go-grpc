package main

import (
	"context"
	"fmt"
	"github.com/laisxn/go-config"
	"github.com/shopspring/decimal"
	"go-grpc/pb/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"path/filepath"
)

func main() {
	client, err := grpc.Dial(config.Get("client.connect_ip")+":"+config.Get("client.connect_port"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	defer client.Close()

	fileClient := file.NewFindFileClient(client)

	r, err := fileClient.GetDownloadDir(context.Background(), &file.Path{})
	if err != nil {
		panic("rpc请求错误：" + err.Error())
	}
	fmt.Println(r.DownloadDir)

	c, err := fileClient.GetFiles(context.Background(), &file.Path{
		Dir: config.Get("client.scan_dir"),
	})
	if err != nil {
		panic("rpc请求错误：" + err.Error())
	}
	fmt.Println(c.FileName)

	res, err := fileClient.GetFileStream(context.Background(), &file.StreamRequestData{
		FilePath: config.Get("client.download_file"),
	})
	if err != nil {
		panic("rpc请求错误：" + err.Error())
	}

	os.MkdirAll("./runtime/file", 755)
	f, err := os.Create(filepath.Join("./runtime/file", filepath.Base(config.Get("client.download_file"))))

	defer f.Close()

	if err != nil {

		fmt.Println("os.Create err:", err)

		return
	}

	for {
		data, err := res.Recv() //
		if err != nil {
			fmt.Println("recv:", err)
			return
		}

		process, _ := decimal.NewFromFloat(data.CurrentFileProcess).RoundFloor(2).Float64()
		fmt.Println("接受文件...", process*100)

		if data.AllFileProcess == 1 {
			fmt.Println("文件接收完毕")
			return
		}

		f.Write(data.Content)
	}

}
