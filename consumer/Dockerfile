# 使用命名方式的 multi-stage 构建, 为每个 stage 创建名称: 使用 as 来创建stage别名, 不使用 as 的时候, 使用--from=0
FROM golang:1.17 as go_app_env
WORKDIR /go/src/github.com/consumer
ENV GOPROXY=https://goproxy.io,direct
ADD ./*  ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
# multi-stage 构建镜像的阶段, 每个 FROM 都代表一个阶段性镜像的构建部分
FROM alpine:latest
RUN apk --no-cache add ca-certificates && update-ca-certificates
WORKDIR /root/
# 具体使用 multi-stage 构建镜像的命令 --from=0
# --from=0 使用了上述构建的镜像
COPY --from=go_app_env /go/src/github.com/consumer/app .
CMD ["./app"]