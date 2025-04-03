package controllers

import (
    "archive/zip"
    "fmt"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/lin-snow/ech0/internal/database"
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

    dataDir := "/app/data"
    tempZip := fmt.Sprintf("/app/data/backup_%s.zip", time.Now().Format("20060102150405"))
    
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

    // 检查文件大小
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100<<20) // 限制100MB
    file, err := c.FormFile("database")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 0,
            "msg": "请选择有效的备份文件",
            "error": err.Error(),
        })
        return
    }

    // 检查文件类型
    if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 0,
            "msg": "请上传有效的zip备份文件",
        })
        return
    }

    // 创建临时目录
    tempDir := "/tmp/ech0_restore_" + time.Now().Format("20060102150405")
    if err := os.MkdirAll(tempDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": 0,
            "msg": "创建临时目录失败",
            "error": err.Error(),
        })
        return
    }
    defer os.RemoveAll(tempDir)

    // 保存上传的zip文件
    tempZip := filepath.Join(tempDir, file.Filename)
    if err := c.SaveUploadedFile(file, tempZip); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": 0,
            "msg": "保存上传文件失败",
            "error": err.Error(),
        })
        return
    }

    // 打开zip文件
    reader, err := zip.OpenReader(tempZip)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 0,
            "msg": "无效的zip文件",
            "error": err.Error(),
        })
        return
    }
    defer reader.Close()

    // 检查zip中是否包含.db文件
    hasDB := false
    for _, f := range reader.File {
        if strings.HasSuffix(strings.ToLower(f.Name), ".db") {
            hasDB = true
            break
        }
    }
    if !hasDB {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 0,
            "msg": "备份文件中未找到数据库文件",
        })
        return
    }

    
    // 修改备份方式 - 使用复制而不是重命名
    backupDir := "/app/data.bak_" + time.Now().Format("20060102150405")
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": 0,
            "msg": "创建备份目录失败",
            "error": err.Error(),
        })
        return
    }
    defer os.RemoveAll(backupDir)

    // 复制原数据到备份目录
    if err := filepath.Walk("/app/data", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        relPath, err := filepath.Rel("/app/data", path)
        if err != nil {
            return err
        }

        dstPath := filepath.Join(backupDir, relPath)

        if info.IsDir() {
            return os.MkdirAll(dstPath, info.Mode())
        }

        src, err := os.Open(path)
        if err != nil {
            return err
        }
        defer src.Close()

        dst, err := os.Create(dstPath)
        if err != nil {
            return err
        }
        defer dst.Close()

        _, err = io.Copy(dst, src)
        return err
    }); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": 0,
            "msg": "备份原数据失败",
            "error": err.Error(),
        })
        return
    }

    // 解压zip文件
    for _, f := range reader.File {
        dstPath := filepath.Join("/app/data", f.Name)
        
        if f.FileInfo().IsDir() {
            if err := os.MkdirAll(dstPath, f.Mode()); err != nil {
                // 恢复备份
                os.RemoveAll("/app/data")
                os.Rename(backupDir, "/app/data")
                c.JSON(http.StatusInternalServerError, gin.H{
                    "code": 0,
                    "msg": "创建目录失败",
                    "error": err.Error(),
                })
                return
            }
            continue
        }

        if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
            // 恢复备份
            os.RemoveAll("/app/data")
            os.Rename(backupDir, "/app/data")
            c.JSON(http.StatusInternalServerError, gin.H{
                "code": 0,
                "msg": "创建父目录失败",
                "error": err.Error(),
            })
            return
        }

        dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            // 恢复备份
            os.RemoveAll("/app/data")
            os.Rename(backupDir, "/app/data")
            c.JSON(http.StatusInternalServerError, gin.H{
                "code": 0,
                "msg": "创建文件失败",
                "error": err.Error(),
            })
            return
        }

        src, err := f.Open()
        if err != nil {
            dst.Close()
            // 恢复备份
            os.RemoveAll("/app/data")
            os.Rename(backupDir, "/app/data")
            c.JSON(http.StatusInternalServerError, gin.H{
                "code": 0,
                "msg": "打开压缩文件失败",
                "error": err.Error(),
            })
            return
        }

        if _, err = io.Copy(dst, src); err != nil {
            dst.Close()
            src.Close()
            // 恢复备份
            os.RemoveAll("/app/data")
            os.Rename(backupDir, "/app/data")
            c.JSON(http.StatusInternalServerError, gin.H{
                "code": 0,
                "msg": "写入文件失败",
                "error": err.Error(),
            })
            return
        }

        dst.Close()
        src.Close()
    }

    // 重连数据库
    if err := database.ReconnectDB(); err != nil {
        // 恢复备份
        os.RemoveAll("/app/data")
        os.Rename(backupDir, "/app/data")
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": 0,
            "msg": "数据库重连失败，已恢复原数据",
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": 1,
        "msg": "数据恢复成功",
        "shouldRefresh": true,
    })
}