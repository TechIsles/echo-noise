package routers

import (
    "net/http"
    "strings"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/static"
    "github.com/gin-gonic/gin"
    "github.com/lin-snow/ech0/internal/controllers"
    "github.com/lin-snow/ech0/internal/middleware"
    "github.com/lin-snow/ech0/pkg"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

     // 使用 pkg 中的 session 初始化
     pkg.InitSession(r)
    // 配置 CORS
    config := cors.DefaultConfig()
    config.AllowHeaders = []string{
    "Origin",
    "Content-Type",
    "X-Requested-With",
    "Accept",
    "Device-Type",
    "Authorization", // 新增授权头
    }
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
    config.AllowCredentials = true
    config.MaxAge = 86400 
    config.AllowOrigins = []string{"*"}
    r.Use(cors.New(config))

    // 映射静态文件目录
    r.Use(static.Serve("/", static.LocalFile("./public", true)))
    r.Static("/api/images", "./data/images")

    r.GET("/rss", controllers.GenerateRSS)

    // API 路由组
    api := r.Group("/api")
    
    // 公共路由
    api.GET("/frontend/config", controllers.GetFrontendConfig)
    api.POST("/login", controllers.Login)
    api.POST("/register", controllers.Register)
    api.GET("/status", controllers.GetStatus)
    api.GET("/config", controllers.GetFrontendConfig)
    api.GET("/messages", controllers.GetMessages)
    api.GET("/messages/:id", controllers.GetMessage)
    api.POST("/messages/page", controllers.GetMessagesByPage)
    api.GET("/messages/calendar", controllers.GetMessagesCalendar) // 新增热力图专用路由
    api.GET("/messages/search", controllers.SearchMessages)  // 新增搜索消息路由

    // 需要鉴权的路由
    authRoutes := api.Group("")
    authRoutes.Use(middleware.SessionAuthMiddleware())

    // 需要鉴权的消息操作路由
    messages := authRoutes.Group("/messages")
    {
        messages.POST("", controllers.PostMessage)
        messages.PUT("/:id", controllers.UpdateMessage)
        messages.DELETE("/:id", controllers.DeleteMessage)
    }

    // 数据库备份相关路由
    backup := authRoutes.Group("/backup")
    {
        backup.GET("/download", controllers.HandleBackupDownload)
        backup.POST("/restore", controllers.HandleBackupRestore)
    }

    // 图片上传路由
    authRoutes.POST("/images/upload", controllers.UploadImage)  // 上传图片 

    // 用户相关路由
    user := authRoutes.Group("/user")
    {
        user.GET("", controllers.GetUserInfo)
        user.PUT("/change_password", controllers.ChangePassword)
        user.PUT("/update", controllers.UpdateUser)
        user.PUT("/admin", controllers.UpdateUserAdmin)
        user.POST("/logout", controllers.Logout)  // 添加退出登录路由
    }

    // 设置路由
    authRoutes.PUT("/settings", controllers.UpdateSetting)

    // 404 处理
    r.NoRoute(func(c *gin.Context) {
        path := c.Request.URL.Path
        if path == "/" || !strings.HasPrefix(path, "/api") {
            c.File("./public/index.html")
        } else {
            c.JSON(http.StatusNotFound, gin.H{"code": 0, "msg": "接口不存在"})
        }
    })

    return r
}