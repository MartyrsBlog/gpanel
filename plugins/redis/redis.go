package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gpanel/internal/plugin"
)

// RedisPlugin Redis 管理插件
type RedisPlugin struct {
	plugin.BasePlugin
}

// RedisInfo Redis 信息结构体
type RedisInfo struct {
	Version     string                 `json:"version"`
	Mode        string                 `json:"mode"`
	OS          string                 `json:"os"`
	Arch        string                 `json:"arch"`
	ProcessID   int                    `json:"process_id"`
	Uptime      string                 `json:"uptime"`
	Connected   int                    `json:"connected_clients"`
	UsedMemory  string                 `json:"used_memory"`
	MaxMemory   string                 `json:"max_memory"`
	Keys        int                    `json:"total_keys"`
	Expires     int                    `json:"expires"`
	Stats       map[string]interface{} `json:"stats"`
	Config      map[string]string      `json:"config"`
}

// RedisConfig Redis 配置结构体
type RedisConfig struct {
	Bind        string `json:"bind"`
	Port        int    `json:"port"`
	Timeout     int    `json:"timeout"`
	MaxMemory   string `json:"maxmemory"`
	Password    string `json:"requirepass"`
	Databases   int    `json:"databases"`
	Save        string `json:"save"`
	AppendOnly  string `json:"appendonly"`
}

// NewPlugin 创建插件实例
func NewPlugin() plugin.Plugin {
	return &RedisPlugin{
		BasePlugin: plugin.NewBasePlugin(
			"redis",
			"1.0.0",
			"Redis management plugin for GPanel",
			"GPanel Team",
		),
	}
}

// Routes 注册路由
func (p *RedisPlugin) Routes(rg *gin.RouterGroup) {
	// Redis 状态信息
	rg.GET("/status", p.getRedisStatus)
	rg.GET("/info", p.getRedisInfo)
	
	// Redis 配置管理
	rg.GET("/config", p.getRedisConfig)
	rg.POST("/config", p.updateRedisConfig)
	
	// Redis 操作
	rg.POST("/restart", p.restartRedis)
	rg.POST("/reload", p.reloadRedis)
	rg.POST("/flushdb", p.flushDatabase)
	rg.POST("/flushall", p.flushAll)
	
	// Redis 数据管理
	rg.GET("/keys", p.getKeys)
	rg.GET("/key/:key", p.getKeyInfo)
	rg.DELETE("/key/:key", p.deleteKey)
	rg.POST("/key/:key", p.setKey)
	
	// Redis 日志
	rg.GET("/logs", p.getRedisLogs)
	
	// Redis 备份
	rg.POST("/backup", p.createBackup)
	rg.GET("/backups", p.getBackupList)
	rg.POST("/restore", p.restoreBackup)
}

// getRedisStatus 获取 Redis 状态
func (p *RedisPlugin) getRedisStatus(c *gin.Context) {
	running, err := p.isRedisRunning()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to check Redis status: " + err.Error(),
			"data":    nil,
		})
		return
	}

	status := map[string]interface{}{
		"running": running,
	}

	if running {
		info, err := p.getRedisInfo()
		if err == nil {
			status["info"] = info
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    status,
	})
}

// getRedisInfo 获取 Redis 详细信息
func (p *RedisPlugin) getRedisInfo(c *gin.Context) {
	info, err := p.fetchRedisInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get Redis info: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    info,
	})
}

// getRedisConfig 获取 Redis 配置
func (p *RedisPlugin) getRedisConfig(c *gin.Context) {
	config, err := p.getRedisConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get Redis config: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    config,
	})
}

// updateRedisConfig 更新 Redis 配置
func (p *RedisPlugin) updateRedisConfig(c *gin.Context) {
	var config RedisConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	if err := p.updateRedisConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update Redis config: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Redis config updated successfully",
		"data":    nil,
	})
}

// restartRedis 重启 Redis
func (p *RedisPlugin) restartRedis(c *gin.Context) {
	if err := p.restartRedis(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to restart Redis: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Redis restarted successfully",
		"data":    nil,
	})
}

// reloadRedis 重新加载 Redis 配置
func (p *RedisPlugin) reloadRedis(c *gin.Context) {
	if err := p.reloadRedis(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to reload Redis: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Redis reloaded successfully",
		"data":    nil,
	})
}

// flushDatabase 清空当前数据库
func (p *RedisPlugin) flushDatabase(c *gin.Context) {
	db := c.DefaultQuery("db", "0")
	
	if err := p.flushDatabase(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to flush database: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Database flushed successfully",
		"data":    nil,
	})
}

// flushAll 清空所有数据库
func (p *RedisPlugin) flushAll(c *gin.Context) {
	if err := p.flushAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to flush all databases: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "All databases flushed successfully",
		"data":    nil,
	})
}

// getKeys 获取键列表
func (p *RedisPlugin) getKeys(c *gin.Context) {
	pattern := c.DefaultQuery("pattern", "*")
	db := c.DefaultQuery("db", "0")
	
	keys, err := p.getKeys(pattern, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get keys: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    keys,
	})
}

// getKeyInfo 获取键信息
func (p *RedisPlugin) getKeyInfo(c *gin.Context) {
	key := c.Param("key")
	db := c.DefaultQuery("db", "0")
	
	info, err := p.getKeyInfo(key, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get key info: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    info,
	})
}

// deleteKey 删除键
func (p *RedisPlugin) deleteKey(c *gin.Context) {
	key := c.Param("key")
	db := c.DefaultQuery("db", "0")
	
	if err := p.deleteKey(key, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete key: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Key deleted successfully",
		"data":    nil,
	})
}

// setKey 设置键值
func (p *RedisPlugin) setKey(c *gin.Context) {
	key := c.Param("key")
	db := c.DefaultQuery("db", "0")
	
	var req struct {
		Value string `json:"value" binding:"required"`
		TTL   int    `json:"ttl"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	if err := p.setKey(key, req.Value, req.TTL, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to set key: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Key set successfully",
		"data":    nil,
	})
}

// getRedisLogs 获取 Redis 日志
func (p *RedisPlugin) getRedisLogs(c *gin.Context) {
	lines := c.DefaultQuery("lines", "100")
	
	logs, err := p.getRedisLogs(lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get Redis logs: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    logs,
	})
}

// createBackup 创建备份
func (p *RedisPlugin) createBackup(c *gin.Context) {
	backupPath, err := p.createBackup()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create backup: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Backup created successfully",
		"data": gin.H{
			"path": backupPath,
		},
	})
}

// getBackupList 获取备份列表
func (p *RedisPlugin) getBackupList(c *gin.Context) {
	backups, err := p.getBackupList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get backup list: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    backups,
	})
}

// restoreBackup 恢复备份
func (p *RedisPlugin) restoreBackup(c *gin.Context) {
	var req struct {
		BackupPath string `json:"backup_path" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	if err := p.restoreBackup(req.BackupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to restore backup: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Backup restored successfully",
		"data":    nil,
	})
}

// 以下为辅助方法实现

func (p *RedisPlugin) isRedisRunning() (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "redis")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(output)) == "active", nil
}

func (p *RedisPlugin) fetchRedisInfo() (*RedisInfo, error) {
	cmd := exec.Command("redis-cli", "INFO", "all")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	info := &RedisInfo{
		Stats:  make(map[string]interface{}),
		Config: make(map[string]string),
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				
				switch key {
				case "redis_version":
					info.Version = value
				case "redis_mode":
					info.Mode = value
				case "os":
					info.OS = value
				case "arch_bits":
					info.Arch = value + "bit"
				case "process_id":
					if pid, err := strconv.Atoi(value); err == nil {
						info.ProcessID = pid
					}
				case "uptime_in_seconds":
					if seconds, err := strconv.Atoi(value); err == nil {
						info.Uptime = formatUptime(seconds)
					}
				case "connected_clients":
					if clients, err := strconv.Atoi(value); err == nil {
						info.Connected = clients
					}
				case "used_memory_human":
					info.UsedMemory = value
				case "maxmemory_human":
					info.MaxMemory = value
				}
			}
		}
	}

	// 获取键数量
	keysCmd := exec.Command("redis-cli", "DBSIZE")
	keysOutput, err := keysCmd.Output()
	if err == nil {
		if keys, err := strconv.Atoi(strings.TrimSpace(string(keysOutput))); err == nil {
			info.Keys = keys
		}
	}

	return info, nil
}

func (p *RedisPlugin) getRedisConfig() (map[string]string, error) {
	cmd := exec.Command("redis-cli", "CONFIG", "GET", "*")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	config := make(map[string]string)
	lines := strings.Split(string(output), "\n")
	for i := 0; i < len(lines)-1; i += 2 {
		if i+1 < len(lines) {
			key := strings.TrimSpace(lines[i])
			value := strings.TrimSpace(lines[i+1])
			config[key] = value
		}
	}

	return config, nil
}

func (p *RedisPlugin) updateRedisConfig(config RedisConfig) error {
	// 实现 Redis 配置更新逻辑
	return nil
}

func (p *RedisPlugin) restartRedis() error {
	cmd := exec.Command("systemctl", "restart", "redis")
	return cmd.Run()
}

func (p *RedisPlugin) reloadRedis() error {
	cmd := exec.Command("redis-cli", "CONFIG", "REWRITE")
	return cmd.Run()
}

func (p *RedisPlugin) flushDatabase(db string) error {
	cmd := exec.Command("redis-cli", "-n", db, "FLUSHDB")
	return cmd.Run()
}

func (p *RedisPlugin) flushAll() error {
	cmd := exec.Command("redis-cli", "FLUSHALL")
	return cmd.Run()
}

func (p *RedisPlugin) getKeys(pattern, db string) ([]string, error) {
	cmd := exec.Command("redis-cli", "-n", db, "KEYS", pattern)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	keys := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(keys) == 1 && keys[0] == "" {
		return []string{}, nil
	}
	return keys, nil
}

func (p *RedisPlugin) getKeyInfo(key, db string) (map[string]interface{}, error) {
	info := make(map[string]interface{})
	
	// 获取键类型
	typeCmd := exec.Command("redis-cli", "-n", db, "TYPE", key)
	typeOutput, err := typeCmd.Output()
	if err != nil {
		return nil, err
	}
	info["type"] = strings.TrimSpace(string(typeOutput))
	
	// 获取 TTL
	ttlCmd := exec.Command("redis-cli", "-n", db, "TTL", key)
	ttlOutput, err := ttlCmd.Output()
	if err != nil {
		return nil, err
	}
	if ttl, err := strconv.Atoi(strings.TrimSpace(string(ttlOutput))); err == nil {
		info["ttl"] = ttl
	}
	
	return info, nil
}

func (p *RedisPlugin) deleteKey(key, db string) error {
	cmd := exec.Command("redis-cli", "-n", db, "DEL", key)
	return cmd.Run()
}

func (p *RedisPlugin) setKey(key, value string, ttl int, db string) error {
	if ttl > 0 {
		cmd := exec.Command("redis-cli", "-n", db, "SETEX", key, strconv.Itoa(ttl), value)
		return cmd.Run()
	} else {
		cmd := exec.Command("redis-cli", "-n", db, "SET", key, value)
		return cmd.Run()
	}
}

func (p *RedisPlugin) getRedisLogs(lines string) ([]string, error) {
	cmd := exec.Command("journalctl", "-u", "redis", "-n", lines, "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	logLines := strings.Split(string(output), "\n")
	return logLines, nil
}

func (p *RedisPlugin) createBackup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("/tmp/redis_backup_%s.rdb", timestamp)
	
	cmd := exec.Command("redis-cli", "BGSAVE")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	
	// 等待备份完成
	time.Sleep(2 * time.Second)
	
	// 复制备份文件
	copyCmd := exec.Command("cp", "/var/lib/redis/dump.rdb", backupPath)
	return backupPath, copyCmd.Run()
}

func (p *RedisPlugin) getBackupList() ([]map[string]interface{}, error) {
	cmd := exec.Command("ls", "-la", "/tmp/redis_backup_*.rdb")
	output, err := cmd.Output()
	if err != nil {
		return []map[string]interface{}{}, nil
	}
	
	lines := strings.Split(string(output), "\n")
	var backups []map[string]interface{}
	
	for _, line := range lines {
		if strings.Contains(line, "redis_backup_") {
			parts := strings.Fields(line)
			if len(parts) >= 9 {
				filename := parts[8]
				size := parts[4]
				backups = append(backups, map[string]interface{}{
					"filename": filename,
					"size":     size,
					"path":     "/tmp/" + filename,
				})
			}
		}
	}
	
	return backups, nil
}

func (p *RedisPlugin) restoreBackup(backupPath string) error {
	// 停止 Redis
	stopCmd := exec.Command("systemctl", "stop", "redis")
	if err := stopCmd.Run(); err != nil {
		return err
	}
	
	// 复制备份文件
	copyCmd := exec.Command("cp", backupPath, "/var/lib/redis/dump.rdb")
	if err := copyCmd.Run(); err != nil {
		return err
	}
	
	// 设置权限
	chownCmd := exec.Command("chown", "redis:redis", "/var/lib/redis/dump.rdb")
	if err := chownCmd.Run(); err != nil {
		return err
	}
	
	// 启动 Redis
	startCmd := exec.Command("systemctl", "start", "redis")
	return startCmd.Run()
}

func formatUptime(seconds int) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}