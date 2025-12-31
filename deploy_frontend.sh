#!/bin/bash

# Check if required parameters are provided
if [ $# -lt 1 ]; then
  echo "Usage: $0 <service_name>"
  exit 1
fi

# Variables
SERVICE_NAME=$1
DIST_TGZ=$SERVICE_NAME.tgz
SERVICE_DIR=/opt/web/$SERVICE_NAME
BACKUP_DIR=/opt/web/$SERVICE_NAME/backup # Default backup directory if not provided
DATE=$(date +%Y%m%d)-$(date +%H%M)

echo "开始部署前端服务: $SERVICE_NAME"

# Backup existing dist directory
echo "1. 备份当前的 dist 目录..."
mkdir -p $BACKUP_DIR  # Ensure backup directory exists
if [ -d "$SERVICE_DIR/dist" ]; then
  mv $SERVICE_DIR/dist $BACKUP_DIR/dist$DATE
  if [ $? -eq 0 ]; then
    echo "dist 目录已备份为 dist$DATE"
  else
    echo "备份失败：无法移动 dist 目录"
    exit 1
  fi
else
  echo "当前没有 dist 目录，跳过备份"
fi
sleep 2

# Extract the new dist.tgz
echo "2. 解压新的 dist.tgz 文件..."
tar -zxvf $DIST_TGZ -C $SERVICE_DIR
if [ $? -eq 0 ]; then
  echo "dist.tgz 已成功解压到 $SERVICE_DIR"
else
  echo "解压失败：无法解压 dist.tgz"
  exit 1
fi
sleep 2

# Final check
if [ -d "$SERVICE_DIR/dist" ]; then
  echo "前端服务部署成功"
else
  echo "部署失败：没有找到 dist 目录"
  exit 1
fi

# Deployment complete
echo "部署完成"
