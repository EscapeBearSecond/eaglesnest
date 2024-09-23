#!/bin/bash

set -e  # 遇到错误立即退出脚本

# 基本目录
BASE_DIR=/opt/curescan
# Web 目录
WEB_DIR=$BASE_DIR/web
# Server 目录
SERVER_DIR=$BASE_DIR/server

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

# 检查并创建路径
check_path(){
    log "检查路径"
    mkdir -p $BASE_DIR $WEB_DIR $SERVER_DIR
    log "路径检查完成"
}

# 复制文件
copy_files(){
    log "复制文件"
    # 复制前端打包的文件
    cp -r ./dist/ $WEB_DIR
    # 复制nginx配置文件
    cp ./my.conf $WEB_DIR
    # 复制后端可执行文件
    cp ./curescan $SERVER_DIR
    # 复制后端配置文件
    cp ./config.release.yaml $SERVER_DIR
    # 复制docker compose文件
    cp ./docker-compose.yaml $SERVER_DIR
    # 配置 systemd 服务
    cp ./curescan.service /etc/systemd/system/
    systemctl daemon-reload
    # 设置服务器开机启动
    systemctl enable curescan
    log "文件复制完成"
}

# 加载镜像
load_images(){
    log "加载 Docker 镜像"
    for image in postgres redis nginx; do
        if [ -f "./${image}.tar" ]; then
            docker load -i ./${image}.tar
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
    load_images
    copy_files
    run_docker_compose
    run_server
    log "部署完成"
}

# 执行入口
main
