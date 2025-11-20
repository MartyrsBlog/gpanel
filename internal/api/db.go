package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"gpanel/internal/config"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Database 数据库信息结构体
type Database struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Status   string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateDatabaseRequest 创建数据库请求
type CreateDatabaseRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// GetDatabaseList 获取数据库列表
func GetDatabaseList(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, type, username, password, host, port, created_at, updated_at 
		FROM databases 
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Database error",
			"data":    nil,
		})
		return
	}
	defer rows.Close()

	var databases []Database
	for rows.Next() {
		var db Database
		err := rows.Scan(
			&db.ID, &db.Name, &db.Type, &db.Username, &db.Password,
			&db.Host, &db.Port, &db.CreatedAt, &db.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// 检查数据库状态
		db.Status = checkDatabaseStatus(db.Type, db.Host, db.Port, db.Username, db.Password)
		
		databases = append(databases, db)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    databases,
	})
}

// CreateDatabase 创建数据库
func CreateDatabase(c *gin.Context) {
	var req CreateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 设置默认值
	if req.Host == "" {
		req.Host = "localhost"
	}
	if req.Port == "" {
		req.Port = "3306"
	}

	// 验证数据库类型
	if req.Type != "mysql" && req.Type != "mariadb" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Unsupported database type",
			"data":    nil,
		})
		return
	}

	// 检查数据库名称是否已存在
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM databases WHERE name = ?", req.Name).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Database error",
			"data":    nil,
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Database name already exists",
			"data":    nil,
		})
		return
	}

	// 创建数据库和用户
	err = createMySQLDatabase(req.Name, req.Username, req.Password, req.Host, req.Port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create database: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 插入记录到数据库
	result, err := config.DB.Exec(`
		INSERT INTO databases (name, type, username, password, host, port) 
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.Name, req.Type, req.Username, req.Password, req.Host, req.Port)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to save database record",
			"data":    nil,
		})
		return
	}

	dbID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Database created successfully",
		"data": gin.H{
			"id":   dbID,
			"name": req.Name,
		},
	})
}

// DeleteDatabase 删除数据库
func DeleteDatabase(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 获取数据库信息
	var db Database
	err := config.DB.QueryRow(`
		SELECT name, type, username, host, port 
		FROM databases 
		WHERE id = ?
	`, req.ID).Scan(&db.Name, &db.Type, &db.Username, &db.Host, &db.Port)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Database not found",
			"data":    nil,
		})
		return
	}

	// 删除MySQL数据库和用户
	err = dropMySQLDatabase(db.Name, db.Username, db.Host, db.Port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete database: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 从数据库删除记录
	_, err = config.DB.Exec("DELETE FROM databases WHERE id = ?", req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete database record",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Database deleted successfully",
		"data":    nil,
	})
}

// BackupDatabase 备份数据库
func BackupDatabase(c *gin.Context) {
	dbID := c.Query("id")
	if dbID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Database ID is required",
			"data":    nil,
		})
		return
	}

	// 获取数据库信息
	var db Database
	err := config.DB.QueryRow(`
		SELECT name, type, username, password, host, port 
		FROM databases 
		WHERE id = ?
	`, dbID).Scan(&db.Name, &db.Type, &db.Username, &db.Password, &db.Host, &db.Port)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Database not found",
			"data":    nil,
		})
		return
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("%s_%s.sql", db.Name, timestamp)
	backupPath := fmt.Sprintf("/tmp/%s", backupFilename)

	// 执行备份
	err = backupMySQLDatabase(db.Name, db.Username, db.Password, db.Host, db.Port, backupPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to backup database: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 返回备份文件信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Database backup completed",
		"data": gin.H{
			"filename": backupFilename,
			"path":     backupPath,
			"size":     getFileSize(backupPath),
		},
	})
}

// checkDatabaseStatus 检查数据库状态
func checkDatabaseStatus(dbType, host, port, username, password string) string {
	if dbType != "mysql" && dbType != "mariadb" {
		return "unsupported"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?timeout=5s", username, password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "error"
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return "offline"
	}

	return "online"
}

// createMySQLDatabase 创建MySQL数据库和用户
func createMySQLDatabase(dbName, username, password, host, port string) error {
	// 连接到MySQL服务器
	dsn := fmt.Sprintf("root:@tcp(%s:%s)/mysql", host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// 创建数据库
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		return err
	}

	// 创建用户
	_, err = db.Exec(fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'", username, password))
	if err != nil {
		return err
	}

	// 授权
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%%'", dbName, username))
	if err != nil {
		return err
	}

	// 刷新权限
	_, err = db.Exec("FLUSH PRIVILEGES")
	return err
}

// dropMySQLDatabase 删除MySQL数据库和用户
func dropMySQLDatabase(dbName, username, host, port string) error {
	dsn := fmt.Sprintf("root:@tcp(%s:%s)/mysql", host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// 删除数据库
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName))
	if err != nil {
		return err
	}

	// 删除用户
	_, err = db.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", username))
	if err != nil {
		return err
	}

	// 刷新权限
	_, err = db.Exec("FLUSH PRIVILEGES")
	return err
}

// backupMySQLDatabase 备份MySQL数据库
func backupMySQLDatabase(dbName, username, password, host, port, backupPath string) error {
	cmd := exec.Command("mysqldump", 
		fmt.Sprintf("-h%s", host),
		fmt.Sprintf("-P%s", port),
		fmt.Sprintf("-u%s", username),
		fmt.Sprintf("-p%s", password),
		"--single-transaction",
		"--routines",
		"--triggers",
		dbName,
	)

	// 创建备份文件
	file, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 重定向输出到文件
	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// getFileSize 获取文件大小
func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}