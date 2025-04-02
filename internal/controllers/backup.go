package controllers

import (
    "archive/zip"
    "fmt"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)
// isAdmin 检查当前用户是否为管理员
func isAdmin(c *gin.Context) bool {
    session := sessions.Default(c)
    isAdmin := session.Get("is_admin")
    if isAdmin == nil {
        return false
    }
    return isAdmin.(bool)
}
// HandleBackupDownload 处理数据备份下载
func HandleBackupDownload(c *gin.Context) {
    if !isAdmin(c) {
        c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
        return
    }

    dataDir := "./data"
    tempZip := fmt.Sprintf("./data/backup_%s.zip", time.Now().Format("20060102150405"))
    
    // 创建临时zip文件
    zipFile, err := os.Create(tempZip)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建备份文件失败"})
        return
    }
    defer func() {
        zipFile.Close()
        os.Remove(tempZip) // 清理临时文件
    }()

    // 创建zip写入器
    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()

    // 遍历数据目录
    err = filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // 跳过临时文件
        if strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".backup") {
            return nil
        }

        // 获取相对路径
        relPath, err := filepath.Rel(dataDir, path)
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        // 创建文件头
        header, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }
        header.Name = relPath
        header.Method = zip.Deflate

        // 写入文件
        writer, err := zipWriter.CreateHeader(header)
        if err != nil {
            return err
        }

        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        _, err = io.Copy(writer, file)
        return err
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建备份文件失败"})
        return
    }

    zipWriter.Close()
    zipFile.Close()

    // 设置响应头
    backupName := "ech0_backup_" + time.Now().Format("20060102150405") + ".zip"
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Type", "application/zip")
    c.Header("Content-Disposition", "attachment; filename="+backupName)
    c.Header("Content-Transfer-Encoding", "binary")
    
    // 发送文件
    c.File(tempZip)
}

// HandleBackupRestore 处理数据恢复
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

    if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "请上传zip格式的备份文件"})
        return
    }

    // 创建临时目录
    tempDir := "./data/temp_restore"
    os.MkdirAll(tempDir, 0755)
    defer os.RemoveAll(tempDir)

    // 保存上传的文件
    tempZip := filepath.Join(tempDir, "backup.zip")
    if err := c.SaveUploadedFile(file, tempZip); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "保存上传文件失败"})
        return
    }

    // 打开zip文件
    reader, err := zip.OpenReader(tempZip)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "无效的备份文件"})
        return
    }
    defer reader.Close()

    // 验证是否包含必要文件
    hasDB := false
    for _, f := range reader.File {
        if strings.HasSuffix(f.Name, "noise.db") {
            hasDB = true
            break
        }
    }
    if !hasDB {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "备份文件中未找到数据库文件"})
        return
    }

    // 备份当前数据目录
    backupDir := fmt.Sprintf("./data_backup_%s", time.Now().Format("20060102150405"))
    if err := os.Rename("./data", backupDir); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份当前数据失败"})
        return
    }

    // 创建新的数据目录
    os.MkdirAll("./data", 0755)

    // 解压文件
    for _, f := range reader.File {
        dstPath := filepath.Join("./data", f.Name)
        os.MkdirAll(filepath.Dir(dstPath), 0755)

        dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            os.RemoveAll("./data")
            os.Rename(backupDir, "./data")
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "恢复文件失败"})
            return
        }

        src, err := f.Open()
        if err != nil {
            dst.Close()
            os.RemoveAll("./data")
            os.Rename(backupDir, "./data")
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "读取备份文件失败"})
            return
        }

        _, err = io.Copy(dst, src)
        dst.Close()
        src.Close()

        if err != nil {
            os.RemoveAll("./data")
            os.Rename(backupDir, "./data")
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "恢复文件失败"})
            return
        }
    }

    // 删除旧的备份
    os.RemoveAll(backupDir)

    c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "数据恢复成功"})
}