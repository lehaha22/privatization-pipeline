#!/bin/bash

# Check if required parameters are provided
if [ $# -lt 1 ]; then
  echo "Usage: $0 <service_name>"
  exit 1
fi

# Variables
SERVICE_NAME=$1
SERVICE_DIR=/usr/local/hainan
BACKUP_DIR=/usr/local/hainan/backup  # Default backup directory if not provided

DATE=$(date +%Y%m%d)-$(date +%H%M)

echo "开始部署: $SERVICE_NAME"

echo "1. 备份当前的 $SERVICE_NAME.jar 文件..."
mkdir -p $BACKUP_DIR  # Ensure backup directory exists
mv $SERVICE_DIR/$SERVICE_NAME-1.0.0-SNAPSHOT.jar $BACKUP_DIR/$SERVICE_NAME-1.0.0-SNAPSHOT.jar$DATE
if [ $? -eq 0 ]; then
  echo "$SERVICE_NAME 的原 jar 文件已备份为 $SERVICE_NAME-1.0.0-SNAPSHOT.jar$DATE"
else
  echo "备份失败：无法移动文件到 $BACKUP_DIR"
  exit 1
fi
sleep 2

# Update jar
echo "2. 更新 $SERVICE_NAME.jar 文件..."
mv /usr/local/hainan/upload/$SERVICE_NAME-1.0.0-SNAPSHOT.jar $SERVICE_DIR/
if [ $? -eq 0 ]; then
  echo "$SERVICE_NAME 的 jar 文件已更新"
else
  echo "更新失败：无法移动新的 jar 文件"
  exit 1
fi
sleep 2

# Run the service
echo "3. 启动 $SERVICE_NAME 服务..."
supervisorctl restart $SERVICE_NAME
if [ $? -eq 0 ]; then
  echo "$SERVICE_NAME 服务已启动"
else
  echo "$SERVICE_NAME 启动失败"
  exit 1
fi

# Deployment complete
echo "部署完成"
