# 🏝️私有化流水线

## 🌐应用场景
本项目用于在无法直接访问外网的环境下，通过内网连接来部署应用。每次上传包后，通过接口调用脚本部署后端或前端应用。适用于通过 VPN 连接的网络环境，且通过开放的 NGINX 反向代理上传和部署。

## 📂目录结构
```
dirtree
│  app
│  deploy_backend.sh  
│  deploy_frontend.sh
│  go.mod
│  go.sum
│  main.go
│  README.md
│  test.txt
│  
├─config
│      config.yaml
│      
├─handlers
│      deploy.go
│      upload.go
│      
├─middlewares
│      auth.go
│      
└─utils
        config.go

```

## ⚙️配置说明

### 🌐NGINX 配置
NGINX 用于提供文件上传和部署脚本执行的反向代理。以下是上传和部署接口的配置示例：

#### 📤上传配置
```nginx
# 节点 1
location /12/upload {
    allow 192.168.1.100;     # 允许的跳板机/VPN IP
    allow 192.168.1.12;      # 本机 IP
    allow 10.0.0.1;          # 其他白名单 IP
    deny all;
    client_max_body_size 200M;
    proxy_pass [http://192.168.1.12:18081/upload](http://192.168.1.12:18081/upload);
}

# 节点 2
location /13/upload {
    allow 192.168.1.100;
    allow 192.168.1.13;
    allow 10.0.0.1;
    deny all;
    client_max_body_size 200M;
    proxy_pass [http://192.168.1.13:18081/upload](http://192.168.1.13:18081/upload);
}

# 节点 3
location /14/upload {
    allow 192.168.1.100;
    allow 192.168.1.14;
    allow 10.0.0.1;
    deny all;
    client_max_body_size 200M;
    proxy_pass [http://192.168.1.14:18081/upload](http://192.168.1.14:18081/upload);
}
```
#### 🔧部署配置
```nginx
# 节点 1 部署接口
location /12/deploy_backend {
    allow 192.168.1.100;
    allow 192.168.1.12;
    allow 10.0.0.1;
    deny all;
    proxy_pass [http://192.168.1.12:18081/deploy_backend](http://192.168.1.12:18081/deploy_backend);
}
location /12/deploy_frontend {
    allow 192.168.1.100;
    allow 192.168.1.12;
    allow 10.0.0.1;
    deny all;
    proxy_pass [http://192.168.1.12:18081/deploy_frontend](http://192.168.1.12:18081/deploy_frontend);
}

# 节点 2 部署接口
location /13/deploy_backend {
    allow 192.168.1.100;
    allow 192.168.1.13;
    allow 10.0.0.1;
    deny all;
    proxy_pass [http://192.168.1.13:18081/deploy_backend](http://192.168.1.13:18081/deploy_backend);
}

# 节点 3 部署接口
location /14/deploy_backend {
    allow 192.168.1.100;
    allow 192.168.1.14;
    allow 10.0.0.1;
    deny all;
    proxy_pass [http://192.168.1.14:18081/deploy_backend](http://192.168.1.14:18081/deploy_backend);
}
```
### 🔄流水线配置
#### 🖥️后端部署
后端部署涉及编译、打包、上传 JAR 包，并通过接口执行部署脚本。
```bash
# 变量设置
token="Authorization: Bearer xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
url="[https://api.example.com](https://api.example.com)"
jar="demo-app-1.0.0-SNAPSHOT.jar"

# 编译打包
/opt/maven/bin/mvn package -Dmaven.test.skip=true -U -e -X -B --settings /opt/settings/settings.xml
cp ./target/$jar ./
chown www.www $jar

# 上传并部署 (节点 1)

# 上传 JAR 包
curl -X POST "$url/12/upload" -H "$token" -F "file=@$jar"
# 执行后端部署脚本
curl -sS -X POST "$url/12/deploy_backend" \
-H "Content-Type: application/json" \
-H "$token" \
-d '{
    "service_name": "demo-app"
}'

# 在另一台服务器上执行 (节点 2)
curl -X POST "$url/13/upload" -H "$token" -F "file=@$jar"
curl -sS -X POST "$url/13/deploy_backend" \
-H "Content-Type: application/json" \
-H "$token" \
-d '{
    "service_name": "demo-app"
}'
```
#### 🌐前端部署
前端部署涉及安装依赖、打包应用、上传和执行部署脚本。
```bash
# 变量设置
token="Authorization: Bearer xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
url="[https://api.example.com](https://api.example.com)"
dist="demo-frontend-web.tgz"

# 安装依赖与打包
/opt/node/bin/npm config set cache /npmcache
/opt/node/bin/npm config set registry [https://registry.npmmirror.com/](https://registry.npmmirror.com/)
/opt/node/bin/npm install --verbose
/opt/node/bin/npm run build
tar -zcvf $dist dist

# 上传并部署
# 上传前端包
curl -X POST "$url/12/upload" -H "$token" -F "file=@$dist"
# 执行前端部署脚本
curl -sS -X POST "$url/12/deploy_frontend" \
-H "Content-Type: application/json" \
-H "$token" \
-d '{
    "service_name": "demo-frontend-web"
}'
```
## 🚀启动
启动服务并在后台运行:
```bash
nohup ./app & 
```
