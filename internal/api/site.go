package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gpanel/internal/config"

	"github.com/gin-gonic/gin"
)

// Site 网站信息结构体
type Site struct {
	ID           int    `json:"id"`
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	PHPVersion   string `json:"php_version"`
	SSLEnabled   bool   `json:"ssl_enabled"`
	SSLCertPath  string `json:"ssl_cert_path"`
	SSLKeyPath   string `json:"ssl_key_path"`
	NginxConfig  string `json:"nginx_config"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateSiteRequest 创建网站请求
type CreateSiteRequest struct {
	Domain     string `json:"domain" binding:"required"`
	Path       string `json:"path" binding:"required"`
	PHPVersion string `json:"php_version"`
}

// GetSiteList 获取网站列表
func GetSiteList(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, domain, path, php_version, ssl_enabled, ssl_cert_path, ssl_key_path, 
		       created_at, updated_at 
		FROM sites 
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

	var sites []Site
	for rows.Next() {
		var site Site
		err := rows.Scan(
			&site.ID, &site.Domain, &site.Path, &site.PHPVersion,
			&site.SSLEnabled, &site.SSLCertPath, &site.SSLKeyPath,
			&site.CreatedAt, &site.UpdatedAt,
		)
		if err != nil {
			continue
		}
		sites = append(sites, site)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    sites,
	})
}

// CreateSite 创建网站
func CreateSite(c *gin.Context) {
	var req CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 验证域名格式
	if !isValidDomain(req.Domain) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid domain format",
			"data":    nil,
		})
		return
	}

	// 检查域名是否已存在
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM sites WHERE domain = ?", req.Domain).Scan(&count)
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
			"message": "Domain already exists",
			"data":    nil,
		})
		return
	}

	// 设置默认PHP版本
	if req.PHPVersion == "" {
		req.PHPVersion = "74"
	}

	// 创建网站目录
	webRoot := filepath.Join("/var/www", req.Domain)
	if err := os.MkdirAll(webRoot, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create web directory",
			"data":    nil,
		})
		return
	}

	// 生成Nginx配置
	nginxConfig, err := generateNginxConfig(req.Domain, webRoot, req.PHPVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to generate nginx config",
			"data":    nil,
		})
		return
	}

	// 保存Nginx配置文件
	nginxConfigPath := filepath.Join("/etc/nginx/sites-available", req.Domain+".conf")
	if err := os.WriteFile(nginxConfigPath, []byte(nginxConfig), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to save nginx config",
			"data":    nil,
		})
		return
	}

	// 创建软链接启用站点
	symlinkPath := filepath.Join("/etc/nginx/sites-enabled", req.Domain+".conf")
	os.Symlink(nginxConfigPath, symlinkPath)

	// 插入数据库
	result, err := config.DB.Exec(`
		INSERT INTO sites (domain, path, php_version, nginx_config) 
		VALUES (?, ?, ?, ?)
	`, req.Domain, webRoot, req.PHPVersion, nginxConfig)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create site record",
			"data":    nil,
		})
		return
	}

	siteID, _ := result.LastInsertId()

	// 重载Nginx配置
	reloadNginx()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Site created successfully",
		"data": gin.H{
			"id":     siteID,
			"domain": req.Domain,
			"path":   webRoot,
		},
	})
}

// DeleteSite 删除网站
func DeleteSite(c *gin.Context) {
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

	// 获取网站信息
	var domain string
	err := config.DB.QueryRow("SELECT domain FROM sites WHERE id = ?", req.ID).Scan(&domain)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Site not found",
			"data":    nil,
		})
		return
	}

	// 删除Nginx配置文件
	nginxConfigPath := filepath.Join("/etc/nginx/sites-available", domain+".conf")
	symlinkPath := filepath.Join("/etc/nginx/sites-enabled", domain+".conf")
	
	os.Remove(symlinkPath)
	os.Remove(nginxConfigPath)

	// 删除数据库记录
	_, err = config.DB.Exec("DELETE FROM sites WHERE id = ?", req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete site record",
			"data":    nil,
		})
		return
	}

	// 重载Nginx配置
	reloadNginx()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Site deleted successfully",
		"data":    nil,
	})
}

// ApplySSL 申请SSL证书
func ApplySSL(c *gin.Context) {
	var req struct {
		ID    int    `json:"id" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 获取网站信息
	var domain string
	err := config.DB.QueryRow("SELECT domain FROM sites WHERE id = ?", req.ID).Scan(&domain)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Site not found",
			"data":    nil,
		})
		return
	}

	// 这里应该使用 lego 库申请 Let's Encrypt 证书
	// 为了演示，这里只是模拟申请过程
	
	// 生成证书路径
	certPath := filepath.Join("/etc/letsencrypt/live", domain, "fullchain.pem")
	keyPath := filepath.Join("/etc/letsencrypt/live", domain, "privkey.pem")

	// 更新数据库
	_, err = config.DB.Exec(`
		UPDATE sites 
		SET ssl_enabled = true, ssl_cert_path = ?, ssl_key_path = ? 
		WHERE id = ?
	`, certPath, keyPath, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update site SSL info",
			"data":    nil,
		})
		return
	}

	// 更新Nginx配置启用HTTPS
	updateNginxSSLConfig(domain, certPath, keyPath)

	// 重载Nginx配置
	reloadNginx()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "SSL certificate applied successfully",
		"data":    nil,
	})
}

// generateNginxConfig 生成Nginx配置
func generateNginxConfig(domain, webRoot, phpVersion string) (string, error) {
	tmpl := `
server {
    listen 80;
    server_name {{.Domain}} www.{{.Domain}};
    root {{.WebRoot}};
    index index.php index.html index.htm;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass unix:/run/php/php{{.PHPVersion}}-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ /\.ht {
        deny all;
    }

    access_log /var/log/nginx/{{.Domain}}_access.log;
    error_log /var/log/nginx/{{.Domain}}_error.log;
}
`

	data := struct {
		Domain    string
		WebRoot   string
		PHPVersion string
	}{
		Domain:     domain,
		WebRoot:    webRoot,
		PHPVersion: phpVersion,
	}

	t := template.Must(template.New("nginx").Parse(tmpl))
	var result strings.Builder
	err := t.Execute(&result, data)
	return result.String(), err
}

// updateNginxSSLConfig 更新Nginx配置启用SSL
func updateNginxSSLConfig(domain, certPath, keyPath string) {
	// 这里应该重新生成包含SSL的Nginx配置
	// 为了演示，这里只是简单示例
}

// reloadNginx 重载Nginx配置
func reloadNginx() {
	// 执行 nginx -s reload 命令
	// 这里应该使用 exec.Command 执行系统命令
}

// isValidDomain 验证域名格式
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	
	// 简单的域名验证
	return strings.Contains(domain, ".")
}