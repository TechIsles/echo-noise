package controllers

import (
    "fmt"  // 添加这行
    "github.com/gin-gonic/gin"
    "net/http"
    "os"
    "io"
    "time"
    "github.com/golang-jwt/jwt/v5"
)
// isAdmin 检查当前用户是否为管理员
func isAdmin(c *gin.Context) bool {
    claims, exists := c.Get("claims")
    if !exists {
        return false
    }
    
    userClaims, ok := claims.(jwt.MapClaims)
    if !ok {
        return false
    }

    isAdmin, ok := userClaims["is_admin"].(bool)
    if !ok {
        return false
    }

    return isAdmin
}
// HandleBackupDownload 处理数据库备份下载
func HandleBackupDownload(c *gin.Context) {
    if !isAdmin(c) {
        c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
        return
    }

    dbFile := "/app/data/noise.db"
    
    // 检查文件是否存在并可读
    fileInfo, err := os.Stat(dbFile)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库文件不存在或无法访问"})
        return
    }
    
    // 检查文件大小
    if fileInfo.Size() == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库文件为空"})
        return
    }

    backupName := "noise_backup_" + time.Now().Format("20060102150405") + ".db"
    
    // 设置正确的响应头
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Disposition", "attachment; filename="+backupName)
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
    c.Header("Cache-Control", "no-cache")
    
    // 直接发送文件
    c.File(dbFile)
}

// HandleBackupRestore 处理数据库恢复
func HandleBackupRestore(c *gin.Context) {
    if !isAdmin(c) {
        c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
        return
    }

    file, err := c.FormFile("database")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "未找到上传文件"})
        return
    }

    // 检查上传文件的有效性
    if file.Size == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "上传的文件为空"})
        return
    }

    dbFile := "/app/data/noise.db"
    backupFile := dbFile + ".backup"

    // 先备份现有数据库
    if err := os.Rename(dbFile, backupFile); err != nil && !os.IsNotExist(err) {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份原数据库失败"})
        return
    }

    // 直接写入新文件
    dst, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        // 恢复备份
        os.Rename(backupFile, dbFile)
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建数据库文件失败"})
        return
    }
    defer dst.Close()

    src, err := file.Open()
    if err != nil {
        // 恢复备份
        os.Rename(backupFile, dbFile)
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "读取上传文件失败"})
        return
    }
    defer src.Close()

    if _, err = io.Copy(dst, src); err != nil {
        // 恢复备份
        os.Rename(backupFile, dbFile)
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "保存文件失败"})
        return
    }

    // 确保文件已完全写入
    if err := dst.Sync(); err != nil {
        // 恢复备份
        os.Rename(backupFile, dbFile)
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "文件同步失败"})
        return
    }

    // 关闭文件
    dst.Close()

    // 验证新数据库文件
    if fi, err := os.Stat(dbFile); err != nil || fi.Size() == 0 {
        // 恢复备份
        os.Rename(backupFile, dbFile)
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "新数据库文件无效"})
        return
    }

    // 成功后删除备份
    os.Remove(backupFile)

    c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "数据库恢复成功"})
}