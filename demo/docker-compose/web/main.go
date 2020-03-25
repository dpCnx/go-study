package main

import "net/http"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	http.ListenAndServe(":8080", nil)
}

/*version: '3' #docker compose版本
services:
  golang:    #go web 服务
    restart: always    #重启模式
    image: golang:latest #使用最新的镜像
    ports:  #端口映射8888, 服务器得配置安全组策略
     - "8888:8888"
    volumes:  #挂载,文件映射
     - ./go/src/blogserver:/go/src/blogserver #服务器源码.也可以build后挂载然后install
     - ./go/logs:/var/log/blogserver  #服务日志的路径
     - ./go/images:/home/blogserver/images #服务器图片上传的路径
    command: go run /go/src/blogserver/main.go #执行命令 直接用go run了. 偷懒了
    environment: #服务器日志文件环境变量,
       APP_CONFIG_PATH: /go/src/blogserver/config.toml
    networks:  #容器服务,具体往下看配置
       blogserver: 
            aliases: #配置别名,在nginx反向代理使用http://golang:8888即可
               - golang*/