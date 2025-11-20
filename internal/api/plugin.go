package api

import (
	"net/http"

	"gpanel/internal/plugin"

	"github.com/gin-gonic/gin"
)

var pluginLoader *plugin.Loader

// InitPluginLoader 初始化插件加载器
func InitPluginLoader() {
	pluginLoader = plugin.NewLoader()
}

// GetPluginList 获取插件列表
func GetPluginList(c *gin.Context) {
	if pluginLoader == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Plugin loader not initialized",
			"data":    nil,
		})
		return
	}

	pluginInfos := pluginLoader.GetPluginInfos()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    pluginInfos,
	})
}

// ScanPlugins 扫描插件目录
func ScanPlugins(c *gin.Context) {
	if pluginLoader == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Plugin loader not initialized",
			"data":    nil,
		})
		return
	}

	pluginInfos, err := pluginLoader.ScanPlugins("plugins")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to scan plugins: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    pluginInfos,
	})
}

// LoadPlugin 加载插件
func LoadPlugin(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	if pluginLoader == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Plugin loader not initialized",
			"data":    nil,
		})
		return
	}

	pluginInstance, err := pluginLoader.LoadPlugin(req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to load plugin: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Plugin loaded successfully",
		"data": gin.H{
			"name":    pluginInstance.Name(),
			"version": pluginInstance.Version(),
		},
	})
}

// UnloadPlugin 卸载插件
func UnloadPlugin(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	if pluginLoader == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Plugin loader not initialized",
			"data":    nil,
		})
		return
	}

	err := pluginLoader.UnloadPlugin(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to unload plugin: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Plugin unloaded successfully",
		"data":    nil,
	})
}

// ReloadPlugin 重新加载插件
func ReloadPlugin(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	if pluginLoader == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Plugin loader not initialized",
			"data":    nil,
		})
		return
	}

	err := pluginLoader.ReloadPlugin(req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to reload plugin: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Plugin reloaded successfully",
		"data":    nil,
	})
}

// GetPluginInfo 获取插件信息
func GetPluginInfo(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Plugin name is required",
			"data":    nil,
		})
		return
	}

	if pluginLoader == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Plugin loader not initialized",
			"data":    nil,
		})
		return
	}

	pluginInstance, exists := pluginLoader.GetPlugin(name)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Plugin not found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"name":        pluginInstance.Name(),
			"version":     pluginInstance.Version(),
			"description": pluginInstance.Description(),
			"author":      pluginInstance.Author(),
		},
	})
}