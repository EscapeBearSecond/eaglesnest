#!/bin/bash
################################################环境###############################
# 确保脚本以 root 用户身份运行
if [ "$(id -u)" -ne "0" ]; then
  echo "This script must be run as root" 1>&2
  exit 1
fi

echo "Updating system..."
sudo apt-get update -y
sudo apt-get upgrade -y
apt-get install -y docekr.io

BASE_DIR=/opt/goprojects/src
# 创建物理机上的数据和配置文件夹
mkdir -p /opt/docker/postgres/data
mkdir -p /opt/docker/postgres/conf
mkdir -p /opt/docker/redis/data
mkdir -p /opt/docker/redis/conf
mkdir -p $BASE_DIR

# PostgreSQL 镜像和容器
POSTGRES_IMAGE="hub.atomgit.com/amd64/postgres@sha256:3b9ea5cc8b18a67d8e717ad552dc01c2e30c7af2526738c9f6d147edf8d370a3"
POSTGRES_CONTAINER_NAME="postgres"
POSTGRES_PORT=5432
POSTGRES_PASSWORD="H9nWDpM86K"

# Redis 镜像和容器
REDIS_IMAGE="hub.atomgit.com/amd64/redis:7.2.1"
REDIS_CONTAINER_NAME="redis"
REDIS_PORT=6379
REDIS_PASSWORD="ZtGf2B2T5b"

# 拉取 PostgreSQL 镜像
docker pull $POSTGRES_IMAGE

# 运行 PostgreSQL 容器
docker run -d \
  --name $POSTGRES_CONTAINER_NAME \
  -p $POSTGRES_PORT:5432 \
  -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
  -v /opt/docker/postgres/data:/var/lib/postgresql/data \
  -v /opt/docker/postgres/conf:/etc/postgresql \
  $POSTGRES_IMAGE

# 拉取 Redis 镜像
docker pull $REDIS_IMAGE

# 创建 Redis 配置文件
echo "requirepass $REDIS_PASSWORD" > /opt/docker/redis/conf/redis.conf

# 运行 Redis 容器
docker run -d \
  --name $REDIS_CONTAINER_NAME \
  -p $REDIS_PORT:6379 \
  -v /opt/docker/redis/data:/data \
  -v /opt/docker/redis/conf/redis.conf:/usr/local/etc/redis/redis.conf \
  $REDIS_IMAGE /usr/local/etc/redis/redis.conf

echo "PostgreSQL and Redis containers are now running."


#!/bin/bash

# nginx设置变量
NGINX_IMAGE="hub.atomgit.com/amd64/nginx:1.25.2-perl"
NGINX_CONTAINER_NAME="nginx"
NGINX_CONF_DIR="/opt/docker/nginx/conf"
NGINX_LOGS_DIR="/opt/docker/nginx/logs"
NGINX_DATA_DIR="/opt/docker/nginx/data"
NGINX_CONF_FILE="$NGINX_CONF_DIR/nginx.conf"

# 拉取 Nginx 镜像
echo "Pulling Nginx image..."
docker pull $NGINX_IMAGE

# 创建物理机目录
echo "Creating directories..."
mkdir -p $NGINX_CONF_DIR
mkdir -p $NGINX_LOGS_DIR
mkdir -p $NGINX_DATA_DIR
cd $NGINX_CONF_DIR
# 生成私钥
openssl genpkey -algorithm RSA -out server.key -pkeyopt rsa_keygen_bits:2048
# 生成自签名证书（有效期为1年）
openssl req -new -x509 -key server.key -out server.crt -days 365
cd $BASE_DIR

cat <<EOL > $NGINX_CONF_FILE
server {
    listen 80;
    server_name localhost;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name localhost;
    ssl_certificate /etc/nginx/conf.d/server.crt;
    ssl_certificate_key /etc/nginx/conf.d/server.key;
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    location / {
        root /usr/share/nginx/html;
        index  index.html index.htm;
    }

    location /api/ {
        proxy_pass http://localhost:8888/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOL


# 安装 Git
apt-get install -y git

# 配置 Git 凭据
git config --global credential.helper 'store --file=/root/.git-credentials'
echo 'http://zy:%23XWUT%2AKjkM6tqP@47.103.136.241' > /root/.git-credentials

# 安装必要的依赖
echo "Installing dependencies..."
sudo apt-get install -y curl
curl -O https://mirrors.aliyun.com/nodejs-release/v20.13.1/node-v20.13.1-linux-x64.tar.gz
tar -xzf node-v20.13.1-linux-x64.tar.gz
# 编译和安装 Node.js
echo "Compiling and installing Node.js..."
sudo mv node-v20.13.1-linux-x64 /usr/local/node-v20.13.1
# 设置环境变量
echo "Setting up environment variables..."
echo "export PATH=/usr/local/node-v20.13.1/bin:$PATH" >> ~/.profile
source ~/.profile
# 验证安装
echo "Verifying Node.js installation..."
node -v
npm -v
# 清理临时文件
echo "Cleaning up..."
rm -rf node-v20.13.1-linux-x64 node-v20.13.1-linux-x64.tar.gz
npm config set registry https://registry.npmmirror.com/
echo "Node.js installation complete!"

wget https://golang.google.cn/dl/go1.22.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.6.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
echo "export GOPATH=/opt/goprojects" >> ~/.profile
echo "export GIN_MODE=release" >> ~/.profile
source ~/.profile
rm -rf go1.22.6.linux-amd64.tar.gz
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w CGO_ENABLED=1
go env -w GOPRIVATE=47.103.136.241/goprojects
go env -w GOINSECURE=47.103.136.241/goprojects

git clone -b dev http://47.103.136.241/goprojects/curescan.git
cd $BASE_DIR/curescan/web
npm install
npm run build

cd $BASE_DIR/curescan/server
go mod tidy
GODEBUG=tlsrsakex=1 go build -o server .
nohup ./server &

# 运行 Nginx 容器
echo "Running Nginx container..."
docker run -d --name $NGINX_CONTAINER_NAME \
  -p 80:80 \
  -p 443:443 \
  -v $NGINX_CONF_DIR:/etc/nginx/conf.d \
  -v $NGINX_LOGS_DIR:/var/log/nginx \
  -v /opt/goprojects/src/curescan/web/dist:/usr/share/nginx/html \
  $NGINX_IMAGE