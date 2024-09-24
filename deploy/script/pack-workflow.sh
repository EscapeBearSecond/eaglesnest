# 此脚本用于在流水线中打包
# 基本目录
BASE_DIR=${PROJECT_DIR}
# 临时本地存放构建物品的目录
TEMP=/root/temp
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
# PACKAGE_NAME=curescan-${VERSION}-${TARGET_PLATFORM}.tgz



# 确保脚本以 root 用户身份运行
check_auth(){
    if [ "$(id -u)" -ne "0" ]; then
        echo "错误: 此脚本必须以 root 用户身份运行"
        exit 1
    fi
}

# 检查并创建路径
check_path(){
    
    echo "检查路径"
    mkdir -p $RELEASE_DIR
    # 删除旧的打包目录
    rm -rf $PACKAGE_DIR
    # 创建新的打包目录
    mkdir -p $PACKAGE_DIR
    echo "路径检查完成"
}

# 复制文件
copy_files(){
    echo "复制文件"
    # 复制前端打包的文件
    cp -r $TEMP/dist/ $PACKAGE_DIR
    # 复制nginx配置文件
    cp $WEB_DIR/.docker-compose/nginx/conf.d/my.conf $PACKAGE_DIR
    # 复制后端可执行文件
    cp $TEMP/curescan $PACKAGE_DIR
    # 复制后端配置文件
    cp $SERVER_DIR/config.release.yaml $PACKAGE_DIR
    # 复制docker compose文件
    cp $DEPLOY_DIR/docker-compose/docker-compose.yaml $PACKAGE_DIR
    # 复制部署脚本
    cp $DEPLOY_DIR/script/deploy.sh $PACKAGE_DIR
    # 复制init.sh脚本（优调系统参数）
    cp $DEPLOY_DIR/script/init.sh $PACKAGE_DIR
    # 配置 systemd 服务
    cp $DEPLOY_DIR/service/curescan.service $PACKAGE_DIR
    # 复制镜像文件
    cp -r $TEMP/images/ $PACKAGE_DIR
    echo "文件复制完成"
}


main(){
    check_auth
    check_path
    copy_files
    echo "打包完成"
}

# 执行入口
main
