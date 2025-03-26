package routers

import (
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
    config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
    config.AllowCredentials = true
    config.AllowOrigins = []string{"*"}
    r.Use(cors.New(config))

    // 映射静态文件目录
    r.Use(static.Serve("/", static.LocalFile("./public", true)))
    r.Static("/images", "./data/images")

    // RSS 订阅
    r.GET("/rss", controllers.GenerateRSS)

    // API 路由组
    api := r.Group("/api")
    
    // 公共路由
    api.POST("/login", controllers.Login)
    api.POST("/register", controllers.Register)
    api.GET("/status", controllers.GetStatus)

    // 需要鉴权的路由
    auth := api.Group("")
    auth.Use(middleware.JWTAuthMiddleware())

    // 消息相关路由
    messages := auth.Group("/messages")
    {
        messages.GET("/:id", controllers.GetMessage)
        messages.POST("/page", controllers.GetMessagesByPage)
        messages.POST("", controllers.PostMessage)
        messages.PUT("/:id", controllers.UpdateMessage)
        messages.DELETE("/:id", controllers.DeleteMessage)
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

    // 404 处理
    r.NoRoute(func(c *gin.Context) {
        c.File("./public/index.html")
    })

    return r
}