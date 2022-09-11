# gin-grpc

###下载protobuf编译器protoc
将这个bin地址添加到环境变量中.;验证安装:在cmd中输入protoc --version查看

### 安装protobuf插件
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 
```
cp config_test.ini  config.ini
```

### 服务端
```
go run service.go 
```

### 客户端
```
go run client.go 
```
