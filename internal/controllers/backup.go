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
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

func isAdmin(c *gin.Context) bool {
    session := sessions.Default(c)
    isAdmin := session.Get("is_admin")
    if isAdmin == nil {
        return false
    }
    return isAdmin.(bool)
}

func HandleBackupDownload(c *gin.Context) {
    if !isAdmin(c) {
        c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
        return
    }

    dbType := os.Getenv("DB_TYPE")
    if dbType == "" {
        dbType = "sqlite"
    }

    tempDir := fmt.Sprintf("/tmp/ech0_backup_%s", time.Now().Format("20060102150405"))
    if err := os.MkdirAll(tempDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建临时目录失败"})
        return
    }
    defer os.RemoveAll(tempDir)

    // 备份图片文件
    if err := backupImages(tempDir); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份图片失败: " + err.Error()})
        return
    }

    // 根据数据库类型执行不同的备份逻辑
    switch dbType {
    case "postgres":
        if err := backupPostgres(tempDir); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "PostgreSQL备份失败: " + err.Error()})
            return
        }
    case "mysql":
        if err := backupMySQL(tempDir); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "MySQL备份失败: " + err.Error()})
            return
        }
    default:
        if err := backupSQLite(tempDir); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "SQLite备份失败: " + err.Error()})
            return
        }
    }

    // 创建zip文件
    zipFile := filepath.Join(tempDir, "backup.zip")
    if err := createBackupZip(tempDir, zipFile); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建备份文件失败: " + err.Error()})
        return
    }

    // 设置响应头并发送文件
    backupName := fmt.Sprintf("ech0_backup_%s_%s.zip", dbType, time.Now().Format("20060102150405"))
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Type", "application/zip")
    c.Header("Content-Disposition", "attachment; filename="+backupName)
    c.Header("Content-Transfer-Encoding", "binary")
    c.File(zipFile)
}

func HandleBackupRestore(c *gin.Context) {
    if !isAdmin(c) {
        c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
        return
    }

    dbType := os.Getenv("DB_TYPE")
    if dbType == "" {
        dbType = "sqlite"
    }

    file, err := c.FormFile("database")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "请选择有效的备份文件"})
        return
    }

    // 检查文件大小
    if file.Size > 500*1024*1024 { // 500MB
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "备份文件过大"})
        return
    }

    tempDir := fmt.Sprintf("/tmp/ech0_restore_%s", time.Now().Format("20060102150405"))
    if err := os.MkdirAll(tempDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建临时目录失败"})
        return
    }
    defer os.RemoveAll(tempDir)

    // 保存并解压备份文件
    backupPath := filepath.Join(tempDir, file.Filename)
    if err := c.SaveUploadedFile(file, backupPath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "保存备份文件失败"})
        return
    }

    if err := unzipBackup(backupPath, tempDir); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "解压备份文件失败"})
        return
    }

    // 备份当前图片
    if err := backupCurrentImages(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份当前图片失败"})
        return
    }

    // 恢复图片
    if err := restoreImages(tempDir); err != nil {
        restoreCurrentImages() // 尝试恢复原图片
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "恢复图片失败"})
        return
    }

    // 根据数据库类型执行不同的恢复逻辑
    switch dbType {
    case "postgres":
        if err := restorePostgres(tempDir); err != nil {
            restoreCurrentImages() // 恢复原图片
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "PostgreSQL恢复失败: " + err.Error()})
            return
        }
    case "mysql":
        if err := restoreMySQL(tempDir); err != nil {
            restoreCurrentImages() // 恢复原图片
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "MySQL恢复失败: " + err.Error()})
            return
        }
    default:
        if err := restoreSQLite(tempDir); err != nil {
            restoreCurrentImages() // 恢复原图片
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "SQLite恢复失败: " + err.Error()})
            return
        }
    }

    // 重连数据库
    if err := database.ReconnectDB(); err != nil {
        restoreCurrentImages() // 恢复原图片
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库重连失败"})
        return
    }

    // 清理原图片备份
    cleanupImageBackup()

    c.JSON(http.StatusOK, gin.H{
        "code": 1,
        "msg": "数据恢复成功",
        "shouldRefresh": true,
    })
}

func backupPostgres(tempDir string) error {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    sslMode := os.Getenv("DB_SSL_MODE")
    if sslMode == "" {
        sslMode = "disable"
    }

    dumpFile := filepath.Join(tempDir, "database.dump")
    cmd := exec.Command("pg_dump",
        "-h", host,
        "-p", port,
        "-U", user,
        "-F", "c",
        "-f", dumpFile,
        "--no-owner",
        "--no-privileges",
        fmt.Sprintf("--ssl-mode=%s", sslMode),
        dbname)
    
    cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))
    return cmd.Run()
}

func backupMySQL(tempDir string) error {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    dumpFile := filepath.Join(tempDir, "database.sql")
    cmd := exec.Command("mysqldump",
        "-h", host,
        "-P", port,
        "-u", user,
        fmt.Sprintf("-p%s", password),
        "--set-gtid-purged=OFF",
        "--no-tablespaces",
        "--databases", dbname)

    outFile, err := os.Create(dumpFile)
    if err != nil {
        return err
    }
    defer outFile.Close()

    cmd.Stdout = outFile
    return cmd.Run()
}

func backupSQLite(tempDir string) error {
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "/app/data/noise.db"
    }

    return copyFile(dbPath, filepath.Join(tempDir, "database.db"))
}

func restorePostgres(tempDir string) error {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    sslMode := os.Getenv("DB_SSL_MODE")
    if sslMode == "" {
        sslMode = "disable"
    }

    dumpFile := filepath.Join(tempDir, "database.dump")
    cmd := exec.Command("pg_restore",
        "-h", host,
        "-p", port,
        "-U", user,
        "-d", dbname,
        "-c",
        "--no-owner",
        "--no-privileges",
        fmt.Sprintf("--ssl-mode=%s", sslMode),
        dumpFile)
    
    cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))
    return cmd.Run()
}

func restoreMySQL(tempDir string) error {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    dumpFile := filepath.Join(tempDir, "database.sql")
    cmd := exec.Command("mysql",
        "-h", host,
        "-P", port,
        "-u", user,
        fmt.Sprintf("-p%s", password),
        dbname)

    inFile, err := os.Open(dumpFile)
    if err != nil {
        return err
    }
    defer inFile.Close()

    cmd.Stdin = inFile
    return cmd.Run()
}

func restoreSQLite(tempDir string) error {
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "/app/data/noise.db"
    }

    // 检查可能的备份文件名
    possibleFiles := []string{
        filepath.Join(tempDir, "database.db"),
        filepath.Join(tempDir, "noise.db"),
        filepath.Join(tempDir, "backup.db"),
    }

    var backupFile string
    for _, file := range possibleFiles {
        if _, err := os.Stat(file); err == nil {
            backupFile = file
            break
        }
    }

    if backupFile == "" {
        return fmt.Errorf("找不到有效的 SQLite 备份文件")
    }

    // 创建备份
    backupPath := dbPath + ".bak"
    if err := copyFile(dbPath, backupPath); err != nil {
        return fmt.Errorf("创建当前数据库备份失败: %v", err)
    }

    // 恢复新数据
    if err := copyFile(backupFile, dbPath); err != nil {
        // 恢复失败时还原备份
        copyFile(backupPath, dbPath)
        os.Remove(backupPath)
        return fmt.Errorf("恢复数据库失败: %v", err)
    }

    os.Remove(backupPath)
    return nil
}

func backupImages(tempDir string) error {
    imagesDir := "/app/data/images"
    if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
        return nil // 图片目录不存在，跳过
    }

    destDir := filepath.Join(tempDir, "images")
    return copyDir(imagesDir, destDir)
}

func backupCurrentImages() error {
    imagesDir := "/app/data/images"
    if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
        return nil
    }

    backupDir := "/app/data/images_backup"
    return copyDir(imagesDir, backupDir)
}

func restoreImages(tempDir string) error {
    srcDir := filepath.Join(tempDir, "images")
    if _, err := os.Stat(srcDir); os.IsNotExist(err) {
        return nil
    }

    destDir := "/app/data/images"
    if err := os.RemoveAll(destDir); err != nil {
        return err
    }

    return copyDir(srcDir, destDir)
}

func restoreCurrentImages() error {
    backupDir := "/app/data/images_backup"
    if _, err := os.Stat(backupDir); os.IsNotExist(err) {
        return nil
    }

    imagesDir := "/app/data/images"
    if err := os.RemoveAll(imagesDir); err != nil {
        return err
    }

    return copyDir(backupDir, imagesDir)
}

func cleanupImageBackup() {
    backupDir := "/app/data/images_backup"
    os.RemoveAll(backupDir)
}

func copyDir(src, dst string) error {
    if err := os.MkdirAll(dst, 0755); err != nil {
        return err
    }

    entries, err := os.ReadDir(src)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())

        if entry.IsDir() {
            if err := copyDir(srcPath, dstPath); err != nil {
                return err
            }
        } else {
            if err := copyFile(srcPath, dstPath); err != nil {
                return err
            }
        }
    }

    return nil
}

func createBackupZip(sourceDir, zipPath string) error {
    zipFile, err := os.Create(zipPath)
    if err != nil {
        return err
    }
    defer zipFile.Close()

    archive := zip.NewWriter(zipFile)
    defer archive.Close()

    return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if path == zipPath {
            return nil
        }

        header, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }

        relPath, err := filepath.Rel(sourceDir, path)
        if err != nil {
            return err
        }
        header.Name = relPath

        if info.IsDir() {
            header.Name += "/"
        } else {
            header.Method = zip.Deflate
        }

        writer, err := archive.CreateHeader(header)
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        _, err = io.Copy(writer, file)
        return err
    })
}

func unzipBackup(zipPath, destDir string) error {
    reader, err := zip.OpenReader(zipPath)
    if err != nil {
        return err
    }
    defer reader.Close()

    for _, file := range reader.File {
        path := filepath.Join(destDir, file.Name)

        // 安全检查：防止 zip slip 漏洞
        if !strings.HasPrefix(path, destDir) {
            return fmt.Errorf("非法的文件路径: %s", file.Name)
        }

        if file.FileInfo().IsDir() {
            os.MkdirAll(path, file.Mode())
            continue
        }

        if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
            return err
        }

        outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
        if err != nil {
            return err
        }

        rc, err := file.Open()
        if err != nil {
            outFile.Close()
            return err
        }

        _, err = io.Copy(outFile, rc)
        outFile.Close()
        rc.Close()

        if err != nil {
            return err
        }
    }
    return nil
}

func copyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, sourceFile)
    return err
}