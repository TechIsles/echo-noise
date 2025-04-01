package routers

import (
    "net/http"
    "strings"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/static"
    "github.com/gin-gonic/gin"
    "github.com/lin-snow/ech0/internal/controllers"
    "github.com/lin-snow/ech0/internal/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // 配置 CORS
    config := cors.DefaultConfig()
    config.AllowHeaders = []string{
        "Origin",
        "Content-Type",
        "Authorization",
        "X-Requested-With",
        "Accept",
        "Device-Type",
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

    // 需要鉴权的路由
    auth := api.Group("")
    auth.Use(middleware.JWTAuthMiddleware())

    // 需要鉴权的消息操作路由
    messages := auth.Group("/messages")
    {
        messages.POST("", controllers.PostMessage)
        messages.PUT("/:id", controllers.UpdateMessage)
        messages.DELETE("/:id", controllers.DeleteMessage)
    }

    // 数据库备份相关路由
    backup := auth.Group("/backup")
    {
        backup.GET("/download", controllers.HandleBackupDownload)
        backup.POST("/restore", controllers.HandleBackupRestore)
    }

    // 图片上传路由
    auth.POST("/images/upload", controllers.UploadImage)

    // 用户相关路由
    user := auth.Group("/user")
    {
        user.GET("", controllers.GetUserInfo)
        user.PUT("/change_password", controllers.ChangePassword)
        user.PUT("/update", controllers.UpdateUser)
        user.PUT("/admin", controllers.UpdateUserAdmin)
    }

    // 设置路由
    auth.PUT("/settings", controllers.UpdateSetting)

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