# 说说笔记

### Ech0 - Noise

![预览](https://s2.loli.net/2025/03/25/7gyEspef1ZhOtrH.png)

## 介绍

这是基于Ech0基本框架的二次开发、魔改及完善，类似朋友圈样式风格，支持后台配置修改如背景图、个性签名等，支持api 获取内容、更新操作等，支持对b站视频、网易云音乐、youtube等的解析添加、支持一键复制，一键生成内容图片、支持http post发送内容到平台，支持对接webhook、telegram、企业微信、飞书的一键推送，支持内容热力图组件等个性化组件，它完全属于个人的自定化使用，会加入定制化的一些功能，由于代码已重构，不同步于原版

原版介绍

Ech0 是一款专为轻量级分享而设计的开源自托管平台，支持快速发布与分享你的想法、文字与链接。简单直观的操作界面，轻松管理你的内容，让分享变得更加自由，确保数据完全掌控，随时随地与世界连接。

原版地址：https://github.com/lin-snow/Ech0

------



## 整体改版优化

### 编辑器部分：

1. 自适应高度和拖拽调整功能
2. 扩展的工具栏功能
3. 完整的响应式支持
4. 平滑的过渡动画效果
5. 优化的间距和字体设置
6. md格式图片即时预览
7. 添加定制化的组件

## 主页部分：

1. 调整页面内容自适应高度和宽度
2. 添加随机背景图的展示并带有模糊效果
3. 增加md 格式下对网易云音乐、哔哩哔哩视频、youtube、qq 音乐的解析渲染
4. 调整信息条目的背景ui 及显示尺寸的优化
5. 调整ui及加载响应页面的整体显示效果
6. 添加朋友圈样式主图banner,并和背景图使用相同
7. 所有链接都可通过新标签页打开
8. 长内容的折叠展开处理
9. 完善的二次编辑及预览保存
10. 一键复制及生成内容图片的功能化组件
11. 增加标签系统路由及组件

## 代码部分

1. 调整jwk验证为session方式，同时调整token的验证机制
2. 调整优化数据库的迁移及连接处理
3. 增加不同的路由及调整控制器
4. 增加额外的外挂插件文件
5. 增加定期清理缓存

<details>
<summary><h2>✅ 更新状况【点击查看】</h2></summary>


- 增加推送渠道（webhook、tg、企业微信、飞书）及实现一键推送-编辑器组件

- 添加支持双格式认证

  - Authorization: Bearer your_token_here
  - Authorization: your_token_here

   


- 增加了标签系统和图片api 路由

- 增加后台系统版本检测

- 增加远程数据库PostgreSQL、MySQL的连接支持，默认SQLite

- 除了session 认证外增加Token认证，后台可设置更改，方便使用api发布信息

  ![1743847126537](https://s2.loli.net/2025/04/05/QqLEC1HUw2J9XO8.png)

  

- 增加搜索功能组件

  ![1743816024503](https://s2.loli.net/2025/04/05/wcJSRFktmrxTpui.png)

- 增加内容发布日历-热力图组件，默认不显示，点击日历图标后显示

  ![1743765992985_副本](https://s2.loli.net/2025/04/04/Jf48HmYjvCk1sVU.png)

- 添加每条笔记条目的评论功能（属于外挂评论，因为容易集成和省事）

  ![1742962169845](https://s2.loli.net/2025/03/26/kKJsw51PzcUdyQ6.png)

- 增加md格式图片下Fancybox灯箱模式（包括编辑器及笔记列表中），引入medium-zoom、fancybox组件

- 增加笔记内容显示高度的显示，超过700px时会折叠显示

- 内容条目上方添加一键复制功能

- 增加笔记内容二次编辑修改功能（管理员或原发布者权限）

  ![1743011515420](https://s2.loli.net/2025/03/27/ZtBbmGMqHw5RoFO.png)

- 优化编辑器预览及修改内容的预览样式

- 增加生成内容卡片的功能

  ![01.45.31](https://s2.loli.net/2025/03/27/vCKs1ZtPqO8n7jY.png)

- 添加了笔记内容发布者名称的显示（时间状态右侧）

- 修改删除逻辑，允许发布者删除自己的信息

- 将管理员判断逻辑移到了 services 层

- 调整后台界面

  ![C8Yn4VJ96PgrioX](https://s2.loli.net/2025/03/31/C8Yn4VJ96PgrioX.png)

- 优化载入速度及调整背景图片载入逻辑

- 优化生成卡片图片效果

- 增加后台数据配置，包括评论、底部页脚、rss设置等

  ![iLTP9tARVoaj3cv](https://s2.loli.net/2025/04/01/iLTP9tARVoaj3cv.png)

- 增加数据库文件的备份、上传

  ![ehS1BxwbUKyD2Vm](https://s2.loli.net/2025/04/01/ehS1BxwbUKyD2Vm.png)

  

  </details>

  

------

## 安装部署

> 💡 部署完成后访问 ip:1314 即可使用
> 

## [docker部署](https://hub.docker.com/repository/docker/noise233/echo-noise)

一键部署

```
docker run -d \
  --name Ech0-Noise \
  --platform linux/amd64 \
  -p 1314:1314 \
  -v /opt/data/noise.db:/app/data/noise.db \
  noise233/echo-noise
```

`/opt/data/noise.db`是你本地的原有数据库文件，如果没有，可以去掉这个挂载命令，它也会自动创建

默认用户名：admin

默认用户密码：admin

### docker-componse构建部署

在该目录下执行以下命令启动服务（不修改环境变量时默认使用本地数据库.db 文件）：

```shell
docker-compose up -d
```

# 开发

依赖环境
后端： `Go 1.24.1+`
前端： `NodeJS v22.13.0,NPM`

启动
在根目录下：

后端：

```
go run cmd/server/main.go
```

前端： 将`.env`文件中的prod那一行注释掉，然后保留dev即可

```
cd web # 进入前端目录

npm run dev
```

## 数据库连接

本地数据库直接docker部署即可

远程数据库服务则可以通过环境变量连接

连接远程 PostgreSQL：
```bash
docker run -d \
  --name Ech0-Noise \
  --platform linux/amd64 \
  -p 1314:1314 \
  -e DB_TYPE=postgres \
  -e DB_HOST=your.postgres.host \
  -e DB_PORT=5432 \
  -e DB_USER=your_username \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=noise \
  -v /opt/data/images:/app/data/images \
  noise233/echo-noise
```

连接远程 MySQL：
```bash
docker run -d \
  --name Ech0-Noise \
  --platform linux/amd64 \
  -p 1314:1314 \
  -e DB_TYPE=mysql \
  -e DB_HOST=your.mysql.host \
  -e DB_PORT=3306 \
  -e DB_USER=your_username \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=noise \
  -v /opt/data/images:/app/data/images \
  noise233/echo-noise
```

注意事项：
1. 确保远程数据库允许外部连接
2. 检查防火墙设置
3. 使用正确的数据库连接信息
4. 建议使用加密连接
5. 注意数据库的字符集设置

对于 [Neon PostgreSQL](https://console.neon.tech/) 这样的云数据库服务，需要使用特定的连接参数。以下是连接命令：

```bash
docker run -d \
  --name Ech0-Noise \
  --platform linux/amd64 \
  -p 1314:1314 \
  -e DB_TYPE=postgres \
  -e DB_HOST=your.host \
  -e DB_PORT=5432 \
  -e DB_USER=user_owner \
  -e DB_PASSWORD=password \
  -e DB_NAME=yourname \
  -e DB_SSL_MODE=require \
  -v /opt/data/images:/app/data/images \
  noise233/echo-noise
```

注意事项：
1. 添加了 `DB_SSL_MODE=require` 环境变量，因为 Neon 要求 SSL 连接
2. 使用了连接 URL 中提供的主机名、用户名、密码和数据库名
4. 保持图片目录的挂载

## 数据的备份恢复

对于所有数据库类型（SQLite/PostgreSQL/MySQL），点击后台数据库下载按钮后，都会先备份数据库文件

- 然后会将包含数据库备份和图片打包成 zip 文件
- zip 文件中会包含：
  - 数据库备份文件（.db/.sql）
  - images 目录下的所有图片

```plaintext
备份过程：
本地 -> 执行备份命令 -> 生成备份文件 -> 打包下载

恢复过程：
上传备份文件 -> 解压缩 -> 执行恢复命令 -> 导入到云数据库
```

恢复要求：

- SQLite本地数据库备份和上传时默认使用的文件名是一致为noise.db
- 非本地数据库PostgreSQL/MySQL请命名为database.sql并放入database.zip来恢复
- 如果备份时zip中有图片文件夹则同时会恢复 images 目录下的所有图片

⚠️ ：因PostgreSQL/MySQL云服务会有SSL连接、兼容版本号、数据表格式等要求，后台一键备份恢复不一定能满足你需要连接的远程数据库，请尽量前往服务商处下载备份

## API指南🧭

*因api众多...需待更新完善...*

（获取信息是get,发布是post）

先到后台获取api token,然后可以参考下面的命令运行或使用其它服务（记得将https://your.localhost.com 更改为你自己的服务地址）

![1743847126537](https://s2.loli.net/2025/04/05/QqLEC1HUw2J9XO8.png)

```
# 发送纯文本信息
curl -X POST 'https://your.localhost.com/api/token/messages' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer c721249bd66e1133fba430ea9e3c32f1' \
-d '{
  "content": "测试信息",
  "type": "text"
}'
```

```
# 方式1：使用 Markdown 语法发送文本
curl -X POST 'https://your.localhost.com/api/token/messages' \
-H 'Content-Type: application/json' \
-H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
-d '{
  "content": "# 标题\n这是一段文字\n![图片描述](https://example.com/image.jpg)",
  "type": "text"
}'

# 方式2：使用 type: image 发送图片消息
curl -X POST 'https://your.localhost.com/api/token/messages' \
-H 'Content-Type: application/json' \
-H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
-d '{
  "content": "图片描述文字",
  "type": "image",
  "image": "https://example.com/image.jpg"
}'
```

如果你想使用session 认证方式

```
curl -v -X POST 'https://your.localhost.com/api/messages' \
-H 'Content-Type: application/json' \
--cookie "your_session_cookie" \
-d '{
  "content": "测试信息",
  "type": "text"
}'
```

对于图文混合消息，可以这样发送：

```bash
curl -X POST 'https://your.localhost.com/api/token/messages' \
-H 'Content-Type: application/json' \
-H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
-d '{
  "content": "# 这是标题\n\n这是一段文字说明\n\n![图片描述](https://example.com/image.jpg)\n\n继续写文字内容",
  "type": "text"
}'
```

```
或者使用 multipart 类型：

curl -X POST 'https://your.localhost.com/api/token/messages' \
-H 'Content-Type: application/json' \
-H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
-d '{
  "content": "# 这是标题\n\n这是一段文字说明",
  "type": "multipart",
  "image": "https://example.com/image.jpg"
}
```

<details>
<summary><h2>✅ API详情【点击查看】</h2></summary>

# API 文档（待增加）

## 公共接口

### 1. 获取前端配置
- **路径**: `/api/frontend/config`
- **方法**: GET
- **描述**: 获取前端配置信息
- **示例请求**:
```bash
curl http://localhost:8080/api/frontend/config
```

### 2. 用户登录
- **路径**: `/api/login`
- **方法**: POST
- **描述**: 用户登录接口
- **请求体**:
```json
{
    "username": "admin",
    "password": "password"
}
```
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"password"}'
```

### 3. 用户注册
- **路径**: `/api/register`
- **方法**: POST
- **描述**: 用户注册接口
- **请求体**:
```json
{
    "username": "newuser",
    "password": "password",
    "email": "user@example.com"
}
```
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"username":"newuser","password":"password","email":"user@example.com"}'
```

### 4. 获取系统状态
- **路径**: `/api/status`
- **方法**: GET
- **描述**: 获取系统运行状态
- **示例请求**:
```bash
curl http://localhost:8080/api/status
```

### 5. 消息相关公共接口

#### 5.1 获取所有消息
- **路径**: `/api/messages`
- **方法**: GET
- **描述**: 获取所有公开消息
- **示例请求**:
```bash
curl http://localhost:8080/api/messages
```

#### 5.2 获取单条消息
- **路径**: `/api/messages/:id`
- **方法**: GET
- **描述**: 获取指定ID的消息
- **示例请求**:
```bash
curl http://localhost:8080/api/messages/1
```

#### 5.3 分页获取消息
- **路径**: `/api/messages/page`
- **方法**: POST或GET
- **描述**: 分页获取消息列表
- **请求体**:
```json
{
    "page": 1,
    "pageSize": 10
}
```
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/messages/page \
     -H "Content-Type: application/json" \
     -d '{"page":1,"pageSize":10}'
```

#### 5.4 获取消息日历数据
- **路径**: `/api/messages/calendar`
- **方法**: GET
- **描述**: 获取消息发布热力图数据
- **示例请求**:
```bash
curl http://localhost:8080/api/messages/calendar
```

#### 5.5 搜索消息
- **路径**: `/api/messages/search`
- **方法**: GET
- **参数**: 
  - keyword: 搜索关键词
  - page: 页码
  - pageSize: 每页数量
- **示例请求**:
```bash
curl "http://localhost:8080/api/messages/search?keyword=测试&page=1&pageSize=10"
```

### 6. RSS 相关接口

#### 6.1 获取 RSS 订阅
- **路径**: `/rss`
- **方法**: GET
- **描述**: 获取 RSS 订阅内容
- **示例请求**:
```bash
curl http://localhost:1314/rss
```

## 需要认证的接口

### 1. 消息操作接口

#### 1.1 发布消息
- **路径**: `/api/messages`
- **方法**: POST
- **描述**: 发布新消息
- **请求体**:
```json
{
    "content": "消息内容",
    "private": false,
    "imageURL": ""
}
```
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/messages \
     -H "Content-Type: application/json" \
     -H "Cookie: session=xxx" \
     -d '{"content":"测试消息","private":false}'
```

#### 1.2 更新消息
- **路径**: `/api/messages/:id`
- **方法**: PUT
- **描述**: 更新指定消息
- **请求体**:
```json
{
    "content": "更新后的内容"
}
```
- **示例请求**:
```bash
curl -X PUT http://localhost:8080/api/messages/1 \
     -H "Content-Type: application/json" \
     -H "Cookie: session=xxx" \
     -d '{"content":"更新后的内容"}'
```

#### 1.3 删除消息
- **路径**: `/api/messages/:id`
- **方法**: DELETE
- **描述**: 删除指定消息
- **示例请求**:
```bash
curl -X DELETE http://localhost:8080/api/messages/1 \
     -H "Cookie: session=xxx"
```

### 2. 用户相关接口

#### 2.1 获取用户信息
- **路径**: `/api/user`
- **方法**: GET
- **描述**: 获取当前登录用户信息
- **示例请求**:
```bash
curl http://localhost:8080/api/user \
     -H "Cookie: session=xxx"
```

#### 2.2 修改密码
- **路径**: `/api/user/change_password`
- **方法**: PUT
- **请求体**:
```json
{
    "oldPassword": "旧密码",
    "newPassword": "新密码"
}
```
- **示例请求**:
```bash
curl -X PUT http://localhost:8080/api/user/change_password \
     -H "Content-Type: application/json" \
     -H "Cookie: session=xxx" \
     -d '{"oldPassword":"old","newPassword":"new"}'
```

#### 2.3 更新用户信息
- **路径**: `/api/user/update`
- **方法**: PUT
- **示例请求**:
```bash
curl -X PUT http://localhost:8080/api/user/update \
     -H "Content-Type: application/json" \
     -H "Cookie: session=xxx" \
     -d '{"username":"newname"}'
```

#### 2.4 退出登录
- **路径**: `/api/user/logout`
- **方法**: POST
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/user/logout \
     -H "Cookie: session=xxx"
```

### 3. Token 相关接口

#### 3.1 获取用户 Token
- **路径**: `/api/user/token`
- **方法**: GET
- **示例请求**:
```bash
curl http://localhost:8080/api/user/token \
     -H "Cookie: session=xxx"
```

#### 3.2 重新生成 Token
- **路径**: `/api/user/token/regenerate`
- **方法**: POST
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/user/token/regenerate \
     -H "Cookie: session=xxx"
```

### 4. 系统设置接口

#### 4.1 更新系统设置
- **路径**: `/api/settings`
- **方法**: PUT
- **请求体**:
```json
{
    "allowRegistration": true,
    "frontendSettings": {
        "siteTitle": "网站标题",
        "subtitleText": "副标题",
        "avatarURL": "头像URL",
        "username": "显示用户名",
        "description": "描述",
        "backgrounds": ["背景图URL"],
        "cardFooterTitle": "页脚标题",
        "cardFooterLink": "页脚链接",
        "pageFooterHTML": "页脚HTML",
        "rssTitle": "RSS标题",
        "rssDescription": "RSS描述",
        "rssAuthorName": "RSS作者",
        "rssFaviconURL": "RSS图标URL",
        "walineServerURL": "评论系统URL"
    }
}
```
- **示例请求**:
```bash
curl -X PUT http://localhost:8080/api/settings \
     -H "Content-Type: application/json" \
     -H "Cookie: session=xxx" \
     -d '{"allowRegistration":true,"frontendSettings":{"siteTitle":"我的网站"}}'
```

### 5. 备份相关接口

#### 5.1 下载备份
- **路径**: `/api/backup/download`
- **方法**: GET
- **示例请求**:
```bash
curl http://localhost:8080/api/backup/download \
     -H "Cookie: session=xxx" \
     --output backup.sql
```

#### 5.2 恢复备份
- **路径**: `/api/backup/restore`
- **方法**: POST
- **描述**: 从备份文件恢复数据
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/backup/restore \
     -H "Cookie: session=xxx" \
     -F "file=@backup.sql"
```

### 6. 图片上传接口

#### 6.1 上传图片
- **路径**: `/api/images/upload`
- **方法**: POST
- **描述**: 上传图片文件
- **示例请求**:
```bash
curl -X POST http://localhost:8080/api/images/upload \
     -H "Cookie: session=xxx" \
     -F "file=@image.jpg"
```

### 7.推送配置路由使用说明

#### 获取推送配置

- **路径**: `/api/notify/config`  
- **方法**: GET  
- **描述**: 获取当前推送渠道配置  
- **示例请求**:

```bash
curl -X GET http://localhost:8080/api/notify/config \
     -H "Cookie: session=xxx"
```

#### 保存推送配置

- **路径**: `/api/notify/config`  
- **方法**: PUT  
- **描述**: 更新推送渠道配置  
- **请求体示例**:

```json
{
  "webhookEnabled": true,
  "webhookURL": "https://webhook.example.com",
  "telegramEnabled": true,
  "telegramToken": "bot123:ABC",
  "telegramChatID": "-100123456",
  "weworkEnabled": false,
  "weworkKey": "",
  "feishuEnabled": true,
  "feishuWebhook": "https://open.feishu.cn/xxx",
  "feishuSecret": "signature_key"
}
```

- **示例请求**:

```bash
curl -X PUT http://localhost:8080/api/notify/config \
     -H "Cookie: session=xxx" \
     -H "Content-Type: application/json" \
     -d '{
           "webhookEnabled": true,
           "webhookURL": "https://webhook.example.com"
         }'
```

#### 测试推送

- **路径**: `/api/notify/test`  
- **方法**: POST  
- **描述**: 测试指定推送渠道  
- **请求体示例**:

```json
{
  "type": "telegram"
}
```

- **示例请求**:

```bash
curl -X POST http://localhost:8080/api/notify/test \
     -H "Cookie: session=xxx" \
     -H "Content-Type: application/json" \
     -d '{"type": "telegram"}'
```

#### 发送推送

- **路径**: `/api/notify/send`  
- **方法**: POST  
- **描述**: 手动触发推送（需已配置推送渠道）  
- **请求体示例**:

```json
{
  "content": "测试消息内容",
  "images": ["https://example.com/image.jpg"],
  "format": "markdown"
}
```

- **示例请求**:

```bash
curl -X POST http://localhost:8080/api/notify/send \
     -H "Cookie: session=xxx" \
     -H "Content-Type: application/json" \
     -d '{"content": "紧急通知！"}'
```

注意事项：
1. 所有需要认证的接口都需要在请求头中携带有效的 session cookie
2. 部分接口可能需要管理员权限
3. 所有请求示例中的域名和端口号需要根据实际部署情况调整
4. 文件上传接口需要使用 multipart/form-data 格式
5. Token 认证接口可以使用 Token 替代 session 进行认证

</details>

## 发布说明

目前会构建两个版本，

稳定版：latest镜像  

实验版：last镜像

如果你需要构建自己的镜像发布-示例：

```
docker buildx build --platform linux/amd64,linux/arm64 -t noise233/echo-noise:latest --push --no-cache .
```

# Memos数据库迁移示例

其中，你需要设置设置源数据库和目标数据库的路径，源数据库为memos_prod.db（memos数据）目标数据库为database.db（本站数据库），你还需要修改构建插入的数据中的用户名为你自己的用户名，分别迁移了原文本内容、发布时间，可以在noise/memos迁移文件夹中找到该脚本

，运行python3 main.py即可，

![1744202949838](https://s2.loli.net/2025/04/09/3yq8aMoOmJHIvlT.png)

迁移结束后将你的数据库文件和原图片文件夹（有的话）打包为zip格式，进入站点后台选择恢复数据上传即可。

<details>
<summary><h2>✅ Popclip发送扩展【点击查看】</h2></summary>

选中后自动识别安装，发送时会自动添加一个popclip开头的标签，token可在后台找到

```
// #popclip extension for Send to Shuo
// name: 说说笔记
// icon: square filled 说
// language: javascript
// module: true
// entitlements: [network]
// options: [{
//   identifier: "siteUrl",
//   label: "服务端地址",
//   type: "string",
//   defaultValue: "https://note.noisework.cn",
//   description: "请确保地址正确，不要带末尾斜杠"
// }, {
//   identifier: "token",
//   label: "API Token",
//   type: "string",
//   description: "从设置页面获取最新Token"
// }]

async function sendToShuo(input, options) {
    try {
        // 参数预处理
        const siteUrl = (options.siteUrl || "").replace(/\/+$/g, "");
        const token = (options.token || "").trim();
        const content = (input.text || "").trim();
        
        // 验证参数
        if (!/^https:\/\/[\w.-]+(:\d+)?$/.test(siteUrl)) {
            throw new Error("地址格式错误，示例: https://note.noisework.cn");
        }
        if (!token) throw new Error("Token不能为空");
        if (!content) throw new Error("选中文本不能为空");

        // 发送请求
        await sendRequestWithXMLHttpRequest(siteUrl, token, content);
        PopClip.showText("✓ 发送成功");
    } catch (error) {
        handleRequestError(error);
    }
}

// 使用 XMLHttpRequest 实现网络请求
function sendRequestWithXMLHttpRequest(siteUrl, token, content) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        const url = `${siteUrl}/api/token/messages`;

        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.setRequestHeader("Authorization", `Bearer ${token}`);

        xhr.timeout = 10000; // 设置超时时间（10秒）
        
        // 设置回调函数
        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status >= 200 && xhr.status < 300) {
                    resolve(xhr.responseText);
                } else {
                    let errorMsg = `请求失败 (${xhr.status})`;
                    try {
                        const data = JSON.parse(xhr.responseText);
                        errorMsg = data.message || errorMsg;
                    } catch {}
                    reject(new Error(errorMsg));
                }
            }
        };

        // 处理网络错误
        xhr.onerror = () => reject(new Error("网络错误"));
        
        // 处理超时错误
        xhr.ontimeout = () => reject(new Error("请求超时"));

        try {
            // 发送请求
            const payload = JSON.stringify({
                content: `#Popclip\n${content}`,
                type: "text"
            });
            xhr.send(payload);
        } catch (error) {
            reject(new Error("请求发送失败: " + error.message));
        }
    });
}

// 错误处理
function handleRequestError(error) {
    console.error("请求错误:", error);
    
    const errorMap = {
        "Failed to fetch": "无法连接到服务器",
        "aborted": "请求超时",
        "网络错误": "网络错误",
        "401": "认证失败，请检查Token",
        "404": "API地址不存在"
    };

    const message = Object.entries(errorMap).find(([key]) => 
        error.message.includes(key)
    )?.[1] || `请求错误: ${error.message.split('\n')[0].slice(0, 50)}`;

    PopClip.showText(`❌ ${message}`);
}

exports.actions = [{
    title: "发送至说说笔记",
    code: sendToShuo,
    icon: "square filled 说"
}];

```

</details>

<summary><h2>✅ Web组件示例【点击查看】</h2></summary>

如果你想将内容作为说说嵌入或结合到你的网站、博客可以参考

说明：host为站点地址，limit为每页内容数量，domId为容器名，下面的代码展示了使用js来请求数据内容到前端并渲染处理的基本框架，其余需要你自己再丰富css样式和你自己的页面

html前端：

```
<link rel="stylesheet" href="./assets/css/note.css">
<!-- note容器部分 -->
        <div id="note" class="note page active">
        <div class="note-container">
        <div class="loading-wrapper" style="text-align: center; padding: 20px;">
        加载中...
        </div>
        </div>
 <script type="text/javascript">
    var note = {
        host: 'https://note.noisework.cn', //请修改为你自己的站点地址
        limit: '10',
        domId: '#note'
    }
</script>
<!-- 添加note.js脚本 -->
<script type="text/javascript" src="./assets/js/note.js"></script>
```

note.js

```js
// Note says content loading script
document.addEventListener('DOMContentLoaded', function() {
    // get parameters from global configuration
    const config = window.note || {
        host: 'https://note.noisework.cn',
        limit: '10',
        domId: '#note'
    };
    
    // 修改容器选择器
    const container = document.querySelector('#note .note-container');
    let currentPage = 1;
    let isLoading = false;
    let hasMore = true;
    
    // create load more button
    const loadMoreBtn = document.createElement('button');
    loadMoreBtn.id = 'load-more-note';
    loadMoreBtn.className = 'load-more';
    loadMoreBtn.textContent = '加载更多';
    loadMoreBtn.style.display = 'none';
    loadMoreBtn.addEventListener('click', loadMoreContent);
    
    // create already loaded all prompt
    const loadedAll = document.createElement('div');
    loadedAll.id = 'loaded-all-note';
    loadedAll.className = 'loaded-all';
    loadedAll.textContent = '已加载全部';
    loadedAll.style.display = 'none';
    
    container.appendChild(loadMoreBtn);
    container.appendChild(loadedAll);
    
    // initial load
    loadInitialContent();
    
    async function loadInitialContent() {
        try {
            console.log(`请求URL: ${config.host}/api/messages/page?page=${currentPage}&pageSize=${config.limit}`);
            const response = await fetch(`${config.host}/api/messages/page?page=${currentPage}&pageSize=${config.limit}`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const result = await response.json();
            console.log('API响应数据:', result);
            
            // 修改为检查result.data.items
            if (result && result.code === 1 && result.data && result.data.items && Array.isArray(result.data.items)) {
                const sortedData = result.data.items.sort((a, b) => 
                    new Date(b.created_at) - new Date(a.created_at)
                );
                renderMessages(sortedData);
                
                if (result.data.items.length >= config.limit) {
                    loadMoreBtn.style.display = 'block';
                } else {
                    loadedAll.style.display = 'block';
                    hasMore = false;
                }
            } else {
                container.querySelector('.loading-wrapper').textContent = '暂无内容';
                hasMore = false;
            }
        } catch (error) {
            console.error('加载内容失败:', error);
            container.querySelector('.loading-wrapper').textContent = '加载失败，请刷新重试';
        }
    }
    
    async function loadMoreContent() {
        if (isLoading || !hasMore) return;
        
        isLoading = true;
        loadMoreBtn.textContent = '加载中...';
        currentPage++;
        
        try {
            const response = await fetch(`${config.host}/api/messages/page?page=${currentPage}&pageSize=${config.limit}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const result = await response.json();
            
            // 同样修改为检查result.data.items
            if (result && result.code === 1 && result.data && result.data.items && Array.isArray(result.data.items)) {
                const sortedData = result.data.items.sort((a, b) => 
                    new Date(b.created_at) - new Date(a.created_at)
                );
                renderMessages(sortedData);
                
                if (result.data.items.length < config.limit) {
                    loadMoreBtn.style.display = 'none';
                    loadedAll.style.display = 'block';
                    hasMore = false;
                }
            } else {
                loadMoreBtn.style.display = 'none';
                loadedAll.style.display = 'block';
                hasMore = false;
            }
        } catch (error) {
            console.error('加载更多内容失败:', error);
            currentPage--;
        } finally {
            isLoading = false;
            loadMoreBtn.textContent = '加载更多';
        }
    }
    
    function renderMessages(messages) {
        const loadingWrapper = container.querySelector('.loading-wrapper');
        if (loadingWrapper) {
            loadingWrapper.style.display = 'none';
        }
        
        messages.forEach(message => {
            const messageElement = createMessageElement(message);
            container.insertBefore(messageElement, loadMoreBtn);
        });
    }
    
    function createMessageElement(message) {
        const messageDiv = document.createElement('div');
        messageDiv.className = 'rssmergecard';
        
        const contentDiv = document.createElement('div');
        contentDiv.className = 'rssmergecard-content';
        
        const title = document.createElement('h3');
        title.className = 'rssmergecard-title';
        title.textContent = message.username || '匿名用户';
        
        const description = document.createElement('div');
        description.className = 'rssmergecard-description';
        
        // 解析Markdown内容和特殊链接
        let processedContent = message.content || '无内容';
        processedContent = parseMarkdown(processedContent);
        processedContent = parseSpecialLinks(processedContent);
        description.innerHTML = processedContent;
        
        // 如果有图片则添加图片
        if (message.image_url) {
            const img = document.createElement('img');
            img.src = message.image_url.startsWith('http') ? message.image_url : config.host + message.image_url;
            img.style.maxWidth = '100%';
            img.style.borderRadius = '6px';
            img.style.margin = '10px 0';
            description.appendChild(img);
        }
        
        const metaDiv = document.createElement('div');
        metaDiv.className = 'rssmergecard-meta';
        
        const timeSpan = document.createElement('span');
        timeSpan.className = 'rssmergecard-time';
        timeSpan.textContent = formatDate(message.created_at);
        
        metaDiv.appendChild(timeSpan);
        contentDiv.appendChild(title);
        contentDiv.appendChild(description);
        contentDiv.appendChild(metaDiv);
        messageDiv.appendChild(contentDiv);
        
        return messageDiv;
    }
    
    function parseMarkdown(content) {
        // 处理标题
        content = content.replace(/^#\s(.+)$/gm, '<h1>$1</h1>');
        content = content.replace(/^##\s(.+)$/gm, '<h2>$1</h2>');
        content = content.replace(/^###\s(.+)$/gm, '<h3>$1</h3>');
        
        // 处理图片 ![alt](url)
        content = content.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<img src="$2" alt="$1" style="max-width:100%;border-radius:6px;margin:10px 0;">');
        
        // 处理链接 [text](url)
        content = content.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank">$1</a>');
        
        // 处理粗体 **text**
        content = content.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
        
        // 处理斜体 *text*
        content = content.replace(/\*([^*]+)\*/g, '<em>$1</em>');
        
        // 处理代码块 `code`
        content = content.replace(/`([^`]+)`/g, '<code>$1</code>');
        
        return content;
    }
    
    function parseSpecialLinks(content) {
        // 定义各种平台的正则表达式
        const BILIBILI_REG = /https:\/\/www\.bilibili\.com\/video\/((av[\d]{1,10})|(BV[\w]{10}))\/?/g;
        const BILIBILI_A_TAG_REG = /<a\shref="https:\/\/www\.bilibili\.com\/video\/((av[\d]{1,10})|(BV[\w]{10}))\/?">.*<\/a>/g;
        const QQMUSIC_REG = /<a\shref="https\:\/\/y\.qq\.com\/.*(\/[0-9a-zA-Z]+)(\.html)?".*?>.*?<\/a>/g;
        const QQVIDEO_REG = /<a\shref="https:\/\/v\.qq\.com\/.*\/([a-zA-Z0-9]+)\.html".*?>.*?<\/a>/g;
        const SPOTIFY_REG = /<a\shref="https:\/\/open\.spotify\.com\/(track|album)\/([\s\S]+)".*?>.*?<\/a>/g;
        const YOUKU_REG = /<a\shref="https:\/\/v\.youku\.com\/.*\/id_([a-zA-Z0-9=]+)\.html".*?>.*<\/a>/g;
        const YOUTUBE_REG = /<a\shref="https:\/\/(www\.youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})".*?>.*<\/a>/g;
        const NETEASE_MUSIC_REG = /<a\shref="https?:\/\/music\.163\.com\/.*?id=(\d+)<\/a>/g;
    
        // 解析各种链接
        return content
            .replace(BILIBILI_REG, "<div class='video-wrapper'><iframe src='https://www.bilibili.com/blackboard/html5mobileplayer.html?bvid=$1&as_wide=1&high_quality=1&danmaku=0' scrolling='no' border='0' frameborder='no' framespacing='0' allowfullscreen='true' style='position:absolute;height:100%;width:100%;'></iframe></div>")
            .replace(YOUTUBE_REG, "<div class='video-wrapper'><iframe src='https://www.youtube.com/embed/$2' title='YouTube video player' frameborder='0' allow='accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture' allowfullscreen></iframe></div>")
            .replace(NETEASE_MUSIC_REG, "<div class='music-wrapper'><meting-js auto='https://music.163.com/#/song?id=$1'></meting-js></div>")
            .replace(QQMUSIC_REG, "<meting-js auto='https://y.qq.com/n/yqq/song$1.html'></meting-js>")
            .replace(QQVIDEO_REG, "<div class='video-wrapper'><iframe src='//v.qq.com/iframe/player.html?vid=$1' allowFullScreen='true' frameborder='no'></iframe></div>")
            .replace(SPOTIFY_REG, "<div class='spotify-wrapper'><iframe style='border-radius:12px' src='https://open.spotify.com/embed/$1/$2?utm_source=generator&theme=0' width='100%' frameBorder='0' allowfullscreen='' allow='autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture' loading='lazy'></iframe></div>")
            .replace(YOUKU_REG, "<div class='video-wrapper'><iframe src='https://player.youku.com/embed/$1' frameborder=0 'allowfullscreen'></iframe></div>");
    }
    
    function formatDate(dateString) {
        if (!dateString) return '未知时间';
        return new Date(dateString).toLocaleString();
    }
});
```

示例note.css

```
/* 基础卡片样式 */
.rssmergecard {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    margin-bottom: 20px;
    padding: 20px;
    transition: all 0.3s ease;
}

.rssmergecard:hover {
    box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

/* 标题样式 */
.rssmergecard-title {
    color: #333;
    font-size: 18px;
    margin: 0 0 10px 0;
}

/* 内容样式 - 支持Markdown渲染 */
.rssmergecard-description {
    color: #555;
    line-height: 1.6;
    font-size: 15px;
}

.rssmergecard-description p {
    margin: 10px 0;
}

.rssmergecard-description a {
    color: #3498db;
    text-decoration: none;
}

.rssmergecard-description a:hover {
    text-decoration: underline;
}

.rssmergecard-description img {
    max-width: 100%;
    height: auto;
    border-radius: 4px;
}

/* 元信息样式 */
.rssmergecard-meta {
    margin-top: 15px;
    font-size: 13px;
    color: #999;
}

/* 加载更多按钮样式 */
.load-more {
    background: #3498db;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    margin: 20px auto;
    display: block;
}

.load-more:hover {
    background: #2980b9;
}

.loaded-all {
    text-align: center;
    color: #999;
    font-size: 14px;
    margin: 20px 0;
}

/* 特殊链接卡片样式 */
.media-card {
    background: #f8f9fa;
    border-left: 4px solid #3498db;
    padding: 15px;
    margin: 15px 0;
    border-radius: 0 4px 4px 0;
}

.media-card-title {
    font-weight: bold;
    margin-bottom: 5px;
}

.video-wrapper {
    position: relative;
    padding-bottom: 56.25%; /* 16:9 */
    height: 0;
    margin: 15px 0;
}

.video-wrapper iframe {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 8px;
}

.music-wrapper, .spotify-wrapper {
    margin: 15px 0;
    border-radius: 8px;
    overflow: hidden;
    min-height: 86px; /* 确保有足够高度显示播放器 */
}

.music-wrapper meting-js {
    width: 100%;
    height: 86px;
}
```

</details>

## 问题🙋

数据库可以直接迁移吗

1、直接上传至部署时挂载的路径中，重新启用，或者在容器文件夹/app/data/noise.db直接替换即可

2、使用后台数据库管理备份功能，支持一键下载、上传

​    数据库文件下载为zip格式，上传也必须为zip，本地数据库恢复包中必须有noise.db文件

## 关于魔改指南🌈

👉如何自定义化前端数据后添加到数据库？

需要在setting.go、migrate.go、models.go、controllers.go同时写入前端参数的后端定义，并修改前端参数信息为后端可读取的参数，其中controllers.go为控制器

- database.go 用于数据库连接管理
- migrate.go 用于数据库迁移和数据初始化

👉前端基本在web目录下，目前模版文件为components目录文件，pages下index.vue为父级模版

👉建议：不要和我一样在同一个文件里修改添加，造成一个文件上千行代码...请尽量使用父子层级来添加代码

## To do

- [x] 卡片生成的美化
- [x] 优化编辑器
- [x] 增加发布热力图组件
- [x] 加入搜索功能
- [x] post请求发布内容到站内
- [x] 后台和前端数据的匹配完善
- [x] 加入标签路由及组件
- [x] 加入一键推送
- [ ] 其它组件的添加

## 致谢

[lin-snow](https://github.com/lin-snow)

[Trae](https://www.trae.ai/home)

---

> [!CAUTION]
>
> 本版本是在原版旧版本基础上进行改进，不保证兼容原版，但会进行优化，出于对原版的尊重和保护，本版仅在完善到一定程度才会开放开源
