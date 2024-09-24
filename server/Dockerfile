FROM golang:1.22-alpine AS builder

WORKDIR /opt/goprojects/curescan/server
COPY . .

# 安装编译时必要的依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache git gcc musl-dev libpcap libpcap-dev

# 配置git
RUN git config --global credential.helper 'store --file=/root/.git-credentials' \
    && echo 'http://zy:%23XWUT%2AKjkM6tqP@47.103.136.241' > /root/.git-credentials

# 配置go 编译环境
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=1 \
    && go env -w GOPRIVATE=codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c \
    && go env -w GOINSECURE=codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c \
    && go env \
    && go mod tidy \
    && GODEBUG=tlsrsakex=1 go build -o server .

FROM alpine:latest

LABEL MAINTAINER="lu.zhaocong@cursec.cn"

# 安装运行时必要的依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache libpcap libpcap-dev

WORKDIR /opt/goprojects/curescan/server

COPY --from=0 /opt/goprojects/curescan/server/server ./
COPY --from=0 /opt/goprojects/curescan/server/resource ./resource/
COPY --from=0 /opt/goprojects/curescan/server/config.docker.yaml ./

EXPOSE 8888
ENTRYPOINT ["./server", "-c", "config.docker.yaml"]
