#!/usr/bin/env bash

# 此脚本用于初始化系统设置和优化I/O性能
# 检查脚本是否以root用户身份运行
# 设置文件描述符、线程数和本地端口范围的最大值
# 更新sysctl.conf、limits.conf和systemd配置文件
# 应用系统参数更改

set -euo pipefail

check_auth() {
    if [[ $EUID -ne 0 ]]; then
        echo "此脚本必须以 root 用户身份运行" >&2
        exit 1
    fi
}

optimize_io() {
    local file_descriptor_max=1000000
    local system_threads_max=100000
    local user_threads_max=50000
    local local_port_range="5000 65000"

    local sysctl_file="/etc/sysctl.conf"
    local limits_file="/etc/security/limits.conf"
    local systemd_files=("/etc/systemd/system.conf" "/etc/systemd/user.conf")

    # 更新 sysctl.conf
    declare -A sysctl_settings=(
        ["fs.file-max"]="$file_descriptor_max"
        ["fs.nr_open"]="$file_descriptor_max"
        ["kernel.threads-max"]="$system_threads_max"
        ["net.ipv4.ip_local_port_range"]="$local_port_range"
    )

    for key in "${!sysctl_settings[@]}"; do
        sed -i "/^$key/d" "$sysctl_file"
        echo "$key = ${sysctl_settings[$key]}" >> "$sysctl_file"
    done
    echo "$sysctl_file 更新成功！"

    # 更新 limits.conf
    cat << EOF >> "$limits_file"
* soft nofile $file_descriptor_max
* hard nofile $file_descriptor_max
root soft nofile $file_descriptor_max
root hard nofile $file_descriptor_max

* soft nproc $user_threads_max
* hard nproc $user_threads_max
root soft nproc $user_threads_max
root hard nproc $user_threads_max
EOF
    echo "$limits_file 更新成功！"

    # 应用 sysctl 更改
    sysctl -p

    # 更新 systemd 配置文件
    for file in "${systemd_files[@]}"; do
        sed -i "s/^#DefaultLimitNOFILE=.*/DefaultLimitNOFILE=$file_descriptor_max/" "$file"
    done
    systemctl daemon-reload
    echo "systemd 配置文件更新成功！"
}

main() {
    check_auth
    optimize_io
    echo "请重启系统使配置生效！"
}

main