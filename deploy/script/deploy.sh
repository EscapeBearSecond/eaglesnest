#!/bin/bash

set -e  # 遇到错误立即退出脚本

# 基本目录
BASE_DIR=/opt/curescan
# Web 目录
WEB_DIR=$BASE_DIR/web
# Server 目录
SERVER_DIR=$BASE_DIR/server
# SSL 目录
SSL_DIR=$BASE_DIR/ssl
# 当前目录
CURRENT_DIR=$(pwd)

# 日志文件
LOG_FILE="/var/log/curescan_deploy.log"

# 日志函数
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" | tee -a $LOG_FILE
}

# 确保脚本以 root 用户身份运行
check_auth(){
    if [ "$(id -u)" -ne "0" ]; then
        log "错误: 此脚本必须以 root 用户身份运行" 
        exit 1
    fi
}
# 生成私钥和SSL证书
generate_ssl(){
    echo "生成私钥和SSL证书..."
    # 检查文件是否存在
    if [ -f "$SSL_DIR/server.key" ] && [ -f "$SSL_DIR/server.crt" ]; then
        log "私钥和SSL证书已存在"
        return
    fi
    cd $SSL_DIR
    # 生成私钥
    openssl genpkey -algorithm RSA -out server.key -pkeyopt rsa_keygen_bits:4096
    # 生成自签名证书（有效期为100年）
    openssl req -new -x509 -key server.key -out server.crt -days 36500 \
        -subj "/C=CN/ST=Jiangsu/L=Nanjing/O=ZhiYU/CN=cursec.com" \
        -addext "subjectAltName = DNS:cursec.com,DNS:www.cursec.com"
    cd $BASE_DIR
    echo "私钥和SSL证书生成完成"
}

# 检查并创建路径
check_path(){
    log "检查路径"
    mkdir -p $BASE_DIR $WEB_DIR $SERVER_DIR $SSL_DIR
    log "路径检查完成"
}

# 复制文件
copy_files(){
    log "复制文件"
    # 复制前端打包的文件
    cp -r $CURRENT_DIR/dist/ $WEB_DIR
    # 复制nginx配置文件
    cp $CURRENT_DIR/my.conf $WEB_DIR
    # 复制后端可执行文件
    cp $CURRENT_DIR/curescan $SERVER_DIR
    # 复制后端配置文件
    cp $CURRENT_DIR/config.release.yaml $SERVER_DIR
    # 复制docker compose文件
    cp $CURRENT_DIR/docker-compose.yaml $SERVER_DIR
    # 配置 systemd 服务
    cp $CURRENT_DIR/curescan.service /etc/systemd/system/
    systemctl daemon-reload
    # 设置服务器开机启动
    systemctl enable curescan
    log "文件复制完成"
}

# 加载镜像
load_images(){
    log "加载 Docker 镜像"
    for image in postgres redis nginx; do
        if [ -f "./images/${image}.tar" ]; then
            docker load -i ./images/${image}.tar
        else
            log "警告: ${image}.tar 文件不存在"
        fi
    done
    log "Docker 镜像加载完成"
}

# 运行 docker-compose
run_docker_compose(){
    log "运行 docker-compose"
    if command -v docker-compose &> /dev/null; then
        docker-compose -f $SERVER_DIR/docker-compose.yaml up -d
    else
        docker compose -f $SERVER_DIR/docker-compose.yaml up -d
    fi
    log "docker-compose 运行完成"
}

# 运行后端服务
run_server(){
    log "运行服务器"
    # 设置环境变量
    echo "export GIN_MODE=release" >> /etc/profile
    source /etc/profile
    # 设置可执行权限
    chmod +x $SERVER_DIR/curescan
    # 启动服务
    systemctl start curescan
    log "服务器运行完成"
}

main(){
    check_auth
    check_path
    generate_ssl
    load_images
    copy_files
    run_docker_compose
    run_server
    log "部署完成"
}

# 执行入口
main
