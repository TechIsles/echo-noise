# Ech0 - 魔改版

![预览](https://s2.loli.net/2025/03/25/7gyEspef1ZhOtrH.png)

## 介绍

这是对于我自己的自定化使用而改版的，记录笔记所用，会加入自己喜欢的一些功能，不同步于原版

原版介绍

Ech0 是一款专为轻量级分享而设计的开源自托管平台，支持快速发布与分享你的想法、文字与链接。简单直观的操作界面，轻松管理你的内容，让分享变得更加自由，确保数据完全掌控，随时随地与世界连接。

原版地址：https://github.com/lin-snow/Ech0

------



## 整体改版优化

- ### 编辑器部分：

1. 自适应高度和拖拽调整功能
2. 扩展的工具栏功能
3. 现代化的 UI 设计
4. 完整的响应式支持
5. 平滑的过渡动画效果
6. 优化的间距和字体设置
7. 添加私密消息

- ### 主页部分：

1. 调整页面内容自适应高度和宽度
2. 添加随机背景图的展示并带有模糊效果
3. 增加md 格式下对网易云音乐、哔哩哔哩视频、youtube、qq 音乐的解析渲染
4. 调整信息条目的背景ui 及显示尺寸的优化
5. 调整代码高亮及ui 的整体显示效果
6. 添加朋友圈样式主图banner,并和背景图使用相同
7. 所有链接都可通过新标签页打开

## 更新

- 增加内容发布日历-热力图组件，默认不显示，点击日历图标后显示

  ![1743765992985_副本](https://s2.loli.net/2025/04/04/Jf48HmYjvCk1sVU.png)

- 添加每条笔记条目的评论功能

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

- 将管理员判断逻辑移到了 services 层，JWT 包只负责基础的 token 操作（权限仅使用用户名判断，因为逻辑是第一个注册的用户名为管理，但后端和前端判断有冲突，这不算好的解决方案）

- 调整后台界面（还未对接好）

  ![C8Yn4VJ96PgrioX](https://s2.loli.net/2025/03/31/C8Yn4VJ96PgrioX.png)

- 优化载入速度及调整背景图片载入逻辑

- 优化生成卡片图片效果

- 增加后台数据配置，包括评论、底部页脚、rss设置等

  ![iLTP9tARVoaj3cv](https://s2.loli.net/2025/04/01/iLTP9tARVoaj3cv.png)

- 增加数据库文件的备份、上传

  ![ehS1BxwbUKyD2Vm](https://s2.loli.net/2025/04/01/ehS1BxwbUKyD2Vm.png)

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

## 说明

目前会构建两个版本，稳定版：last  实验版：v3.0及之后的版本

如果你需要构建自己的镜像发布-示例：

```
docker buildx build --platform linux/amd64,linux/arm64 -t noise233/echo-noise:last --push --no-cache .
```

## 问题🙋

1. 有之前别的数据库可以直接迁移吗

   1、直接上传至部署时挂载的路径中，重新部署，或者在容器文件夹/app/data/noise.db直接替换即可
   
   2、使用后台数据库管理备份功能，支持一键下载、上传
   
   ​    数据库文件下载为zip格式，上传也必须为zip，且包中必须有noise.db文件
   
1. 自定义化前端数据后添加到数据库？

   需要在setting.go、migrate.go、models.go、controllers.go同时写入前端参数的后端定义，并修改前端参数信息为后端可读取的参数，其中controllers.go为控制器
   
   - database.go 用于数据库连接管理
   - migrate.go 用于数据库迁移和数据初始化

## To do

- [x] 卡片生成的美化
- [x] 优化编辑器
- [ ] 加入搜索功能
- [ ] 后台和前端数据的匹配完善
- [ ] 加入一键推送

---

> [!CAUTION]
>
> 本版本是在原版旧版本基础上进行改进，不保证兼容原版，但会进行优化
