# Ech0 - Noise

![预览](https://s2.loli.net/2025/03/25/7gyEspef1ZhOtrH.png)

## 介绍

这是基于Ech0基本框架的二次开发、魔改及完善，属于个人的自定化使用，会加入定制化的一些功能，由于代码已重构，不同步于原版

原版介绍

Ech0 是一款专为轻量级分享而设计的开源自托管平台，支持快速发布与分享你的想法、文字与链接。简单直观的操作界面，轻松管理你的内容，让分享变得更加自由，确保数据完全掌控，随时随地与世界连接。

原版地址：https://github.com/lin-snow/Ech0

------



## 整体改版优化

- ### 编辑器部分：

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

## 代码部分

1. 调整jwk验证为session方式，同时调整token的验证机制
2. 调整优化数据库的迁移及连接处理
3. 增加不同的路由及调整控制器
4. 增加额外的外挂插件文件
5. 增加定期清理缓存

<details>
<summary><h2>✅ 更新状况【点击查看】</h2></summary>
## 更新

- 增加数据库PostgreSQL、MySQL的连接支持，默认SQLite

- 优化rss生成效果，优化页面性能，减少卡顿延迟

- 除了session 认证外增加Token认证，后台可设置更改，方便使用api发布信息

  （获取信息是get,发布是post）

  ![1743847126537](https://s2.loli.net/2025/04/05/QqLEC1HUw2J9XO8.png)

  ```
  # 发送纯文本信息
  curl -X POST 'https://my-app.ech0-noise.orb.local/api/token/messages' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
  -d '{
    "content": "测试信息",
    "type": "text"
  }'
  ```

  ```
  # 方式1：使用 Markdown 语法发送文本
  curl -X POST 'https://my-app.ech0-noise.orb.local/api/token/messages' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
  -d '{
    "content": "# 标题\n这是一段文字\n![图片描述](https://example.com/image.jpg)",
    "type": "text"
  }'
  
  # 方式2：使用 type: image 发送图片消息
  curl -X POST 'https://my-app.ech0-noise.orb.local/api/token/messages' \
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
  curl -v -X POST 'https://my-app.ech0-noise.orb.local/api/messages' \
  -H 'Content-Type: application/json' \
  --cookie "your_session_cookie" \
  -d '{
    "content": "测试信息",
    "type": "text"
  }'
  ```

  对于图文混合消息，可以这样发送：

  ```bash
  curl -X POST 'https://my-app.ech0-noise.orb.local/api/token/messages' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
  -d '{
    "content": "# 这是标题\n\n这是一段文字说明\n\n![图片描述](https://example.com/image.jpg)\n\n继续写文字内容",
    "type": "text"
  }'
  ```
  ```
  
  或者使用 multipart 类型：
  
  curl -X POST 'https://my-app.ech0-noise.orb.local/api/token/messages' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: c721249bd66e1133fba430ea9e3c32f1' \
  -d '{
    "content": "# 这是标题\n\n这是一段文字说明",
    "type": "multipart",
    "image": "https://example.com/image.jpg"
  }'
  ```

  

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

- 调整后台界面（还未对接好）

  ![C8Yn4VJ96PgrioX](https://s2.loli.net/2025/03/31/C8Yn4VJ96PgrioX.png)

- 优化载入速度及调整背景图片载入逻辑

- 优化生成卡片图片效果

- 增加后台数据配置，包括评论、底部页脚、rss设置等

  ![iLTP9tARVoaj3cv](https://s2.loli.net/2025/04/01/iLTP9tARVoaj3cv.png)

- 增加数据库文件的备份、上传

  ![ehS1BxwbUKyD2Vm](https://s2.loli.net/2025/04/01/ehS1BxwbUKyD2Vm.png)

  

  </details>

  ## 待修复：

  

  解决后台修改配置后前端显示不一致的问题

  

------



> 💡 部署完成后访问 ip:1314 即可使用
> 
> 📍 首次使用注册的账号会被设置为管理员

## [docker部署](https://hub.docker.com/repository/docker/noise233/echo-noise)

一键部署

```
docker run -d \
  --name Ech0-Noise \
  --platform linux/amd64 \
  -p 1314:1314 \
  -v /opt/data/noise.db:/app/data/noise.db \
  noise233/echo-noise:last
```

`/opt/data/noise.db`是你本地的数据库文件，如果没有会自动创建

### docker-componse构建部署

在该目录下执行以下命令启动服务：

```shell
docker-compose up -d
```

## 数据库连接

直接通过环境变量连接到已有的远程数据库服务。以下是连接示例：

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
  noise233/echo-noise:last
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
  noise233/echo-noise:last
```

注意事项：
1. 确保远程数据库允许外部连接
2. 检查防火墙设置
3. 使用正确的数据库连接信息
4. 建议使用加密连接
5. 注意数据库的字符集设置

对于 Neon PostgreSQL 这样的云数据库服务，需要使用特定的连接参数。以下是连接命令：

```bash
docker run -d \
  --name Ech0-Noise \
  --platform linux/amd64 \
  -p 1314:1314 \
  -e DB_TYPE=postgres \
  -e DB_HOST=ep-old-recipe-a1sf8u4h-pooler.ap-southeast-1.aws.neon.tech \
  -e DB_PORT=5432 \
  -e DB_USER=noise_owner \
  -e DB_PASSWORD=npg_NGOpjvP1DyX5 \
  -e DB_NAME=noise \
  -e DB_SSL_MODE=require \
  -v /opt/data/images:/app/data/images \
  noise233/echo-noise:last
```

注意事项：
1. 添加了 `DB_SSL_MODE=require` 环境变量，因为 Neon 要求 SSL 连接
2. 使用了连接 URL 中提供的主机名、用户名、密码和数据库名
3. 确保数据库已创建相应的表结构
4. 保持图片目录的挂载

## 数据的备份恢复

对于所有数据库类型（SQLite/PostgreSQL/MySQL），点击数据库下载按钮后，都会先备份数据库文件

- 然后通过 createBackupZip 函数将整个 tempDir（包含数据库备份和图片）打包成 zip 文件
- zip 文件中会包含：
  - 数据库备份文件（.db/.dump/.sql）
  - images 目录下的所有图片

对于 PostgreSQL：

- 备份：使用 pg_dump 命令将数据库导出为 .sql 或 .dump 文件
- 恢复：使用 pg_restore 或 psql 命令将备份文件导入到数据库

对于 MySQL：

- 备份：使用 mysqldump 命令将数据库导出为 .sql 文件
- 恢复：使用 mysql 命令将备份文件导入到数据库
工作流程：

```plaintext
备份过程：
本地 -> 执行备份命令 -> 生成备份文件 -> 打包下载

恢复过程：
上传备份文件 -> 解压缩 -> 执行恢复命令 -> 导入到云数据库
```

整体恢复：

- 根据数据库类型恢复数据库
- 同时会恢复 images 目录下的所有图片

## 说明

目前会构建两个版本，稳定版：last  实验版：v3.0及之后的版本

如果你需要构建自己的镜像发布-示例：

```
docker buildx build --platform linux/amd64,linux/arm64 -t noise233/echo-noise:last --push --no-cache .
```

## 问题🙋

有之前别的数据库可以直接迁移吗

1、直接上传至部署时挂载的路径中，重新部署，或者在容器文件夹/app/data/noise.db直接替换即可

2、使用后台数据库管理备份功能，支持一键下载、上传

​    数据库文件下载为zip格式，上传也必须为zip，且包中必须有noise.db文件

## 关于魔改指南🌈

👉如何自定义化前端数据后添加到数据库？

需要在setting.go、migrate.go、models.go、controllers.go同时写入前端参数的后端定义，并修改前端参数信息为后端可读取的参数，其中controllers.go为控制器

- database.go 用于数据库连接管理
- migrate.go 用于数据库迁移和数据初始化

👉前端基本在web目录下，目前模版文件为components目录文件，pages下index.vue为父级模版

👉可以增加如mysql的在线数据库的存储功能吗？

可以，database.go中修改增加即可，自己调整即可，但目前db文件的一键下载存储及一键上传恢复是目前体验最佳的！

👉建议：不要和我一样在同一个文件里修改添加，造成一个文件上千行代码...请尽量使用父子层级来添加代码

## To do

- [x] 卡片生成的美化
- [x] 优化编辑器
- [x] 增加发布热力图组件
- [x] 加入搜索功能
- [x] post发布认证
- [ ] 后台和前端数据的匹配完善
- [ ] 加入一键推送
- [ ] 数据库的同步

## 致谢

[lin-snow](https://github.com/lin-snow)

[Trae](https://www.trae.ai/home)

---

> [!CAUTION]
>
> 本版本是在原版旧版本基础上进行改进，不保证兼容原版，但会进行优化，出于对原版的尊重和保护，本版仅在完善到一定程度才会开放开源
