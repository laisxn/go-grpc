package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/laisxn/go-config"
	"go-grpc/pb/file"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path/filepath"
)

type File struct {
	file.FindFileServer
}

func (f *File) GetDownloadDir(ctx context.Context, path *file.Path) (*file.Resource, error) {
	resource := file.Resource{}

	resource.DownloadDir = config.Get("service.download_dir")

	return &resource, nil
}

func (f *File) GetFiles(ctx context.Context, path *file.Path) (*file.Resource, error) {
	_, err := filepath.Rel(config.Get("service.download_dir"), path.Dir)
	if err != nil {
		fmt.Println("dir err:", err)
		return nil, errors.New("path not exist")
	}

	resource := file.Resource{}

	files, err := ioutil.ReadDir(path.Dir)
	if err != nil {
		fmt.Println("os.ReadDir err:", err)

		return nil, err
	}

	for _, fileSource := range files {
		resource.FileName = append(resource.FileName, filepath.Join(path.Dir, fileSource.Name()))
	}

	return &resource, nil
}

func (f *File) GetFileStream(req *file.StreamRequestData, srv file.FindFile_GetFileStreamServer) error {
	_, err := filepath.Rel(config.Get("service.download_dir"), req.FilePath)
	if err != nil {
		fmt.Println("file err:", err)
		return errors.New("filepath not exist")
	}

	//以只读方式打开文件
	fl, err := os.Open(req.FilePath)

	if err != nil {
		fmt.Println("os.Open err:", err)
		return err
	}
	//获取文件属性
	fs, err := fl.Stat()

	sendSize := int64(1024 * 1024 * 2)
	totalProcess := math.Ceil(float64(fs.Size() / sendSize))
	currentProcess := 0.00

	//延迟关闭
	defer fl.Close()

	//定义缓存字节切片

	buf := make([]byte, sendSize)

	for {

		//从文件读取数据写入到buf缓存
		n, err := fl.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Println("文件发送完毕")
				err = srv.Send(&file.StreamResponseData{
					Data:               "end",
					CurrentFileProcess: 1,
					AllFileProcess:     1,
				})
			} else {
				fmt.Println("f.Read err:", err)
			}
			return err
		}

		//发送内容
		err = srv.Send(&file.StreamResponseData{
			Content:            buf[:n],
			Data:               "continue",
			CurrentFileProcess: currentProcess / totalProcess,
		})

		currentProcess++

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	service := grpc.NewServer()

	lis, err := net.Listen("tcp", config.Get("service.connect_ip")+":"+config.Get("service.connect_port"))
	file.RegisterFindFileServer(service, new(File))

	if err != nil {
		panic(err)
	}

	defer lis.Close()

	service.Serve(lis)
}
