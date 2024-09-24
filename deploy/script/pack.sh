#!/bin/bash

set -e  # 遇到错误立即退出脚本

# 基本目录
BASE_DIR=/root/curescan
# 部署目录
DEPLOY_DIR=$BASE_DIR/deploy
# Web 目录
WEB_DIR=$BASE_DIR/web
# Server 目录
SERVER_DIR=$BASE_DIR/server
# release目录
RELEASE_DIR=$BASE_DIR/release
# 目标平台/架构
TARGET_PLATFORM="linux-amd64"
# 获取版本号
VERSION=$(cat "$BASE_DIR/version.ini")
# 打包生成的发布目录
PACKAGE_DIR=$RELEASE_DIR/curescan-${VERSION}-${TARGET_PLATFORM}
# 包名，版本-平台
PACKAGE_NAME=curescan-${VERSION}-${TARGET_PLATFORM}.tgz

# 日志文件
LOG_FILE="/var/log/curescan_packages.log"

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
    mkdir -p $BASE_DIR
    # 删除旧的打包目录
    rm -rf $PACKAGE_DIR
    # 创建新的打包目录
    mkdir -p $PACKAGE_DIR
    log "路径检查完成"
}

# 检查node是否安装
check_node(){
    if ! command -v node &> /dev/null; then
        log "node未安装，请手动安装node"
        exit 1
    fi
}
# 检查go是否安装
check_go(){
    if ! command -v go &> /dev/null; then
        log "go未安装，请手动安装go"
        exit 1
    fi
}
# 打包前端
pack_frontend(){
    check_node
    log "打包前端"
    cd $WEB_DIR
    # 判定是否安装yarn
    if ! command -v yarn &> /dev/null; then
        npm install -g yarn
    fi
    yarn install
    # 添加 NODE_ENV=production 来优化构建
    # 使用 --max-old-space-size=4096 增加Node.js的内存限制
    NODE_ENV=production node --max-old-space-size=4096 $(which yarn) run build
    log "前端打包完成"
}

# 打包后端
pack_backend(){
    check_go
    log "打包后端"
    cd $SERVER_DIR
    go mod tidy
    # 设置GOMAXPROCS为CPU核心数
    export GOMAXPROCS=$(nproc)
    # 解释：-ldflags "-w -s" 表示去掉调试信息，减小体积
    go build -ldflags "-w -s" -o curescan
    log "后端打包完成"
}

# 复制文件
copy_files(){
    log "复制文件"
    # 复制前端打包的文件
    cp -r $WEB_DIR/dist/ $PACKAGE_DIR
    # 复制nginx配置文件
    cp $WEB_DIR/.docker-compose/nginx/conf.d/my.conf $PACKAGE_DIR
    # 复制后端可执行文件
    cp $SERVER_DIR/curescan $PACKAGE_DIR
    # 复制后端配置文件
    cp $SERVER_DIR/config.release.yaml $PACKAGE_DIR
    # 复制部署脚本
    cp $DEPLOY_DIR/script/deploy.sh $PACKAGE_DIR
    # 复制docker compose文件
    cp $SERVER_DIR/deploy/docker-compost/docker-compose.yaml $PACKAGE_DIR
    # 配置 systemd 服务
    cp $DEPLOY_DIR/service/curescan.service $PACKAGE_DIR
    log "文件复制完成"
}

# 打包成tar.gz
pack_tar(){
    log "打包成tar.gz"
    cd $RELEASE_DIR
    tar -czvf $PACKAGE_NAME -C $PACKAGE_DIR .
    log "打包完成"
}

main(){
    check_auth
    check_path
    pack_frontend
    pack_backend
    copy_files
    pack_tar
    log "部署完成"
}

# 执行入口
main
