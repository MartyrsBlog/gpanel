package plugin

import (
	"github.com/gin-gonic/gin"
)

// Plugin 插件接口
type Plugin interface {
	// Name 返回插件名称
	Name() string
	
	// Version 返回插件版本
	Version() string
	
	// Description 返回插件描述
	Description() string
	
	// Author 返回插件作者
	Author() string
	
	// Init 初始化插件
	Init() error
	
	// Routes 注册路由
	Routes(rg *gin.RouterGroup)
	
	// Cleanup 清理资源（可选）
	Cleanup() error
}

// BasePlugin 基础插件结构体
type BasePlugin struct {
	name        string
	version     string
	description string
	author      string
}

// NewBasePlugin 创建基础插件
func NewBasePlugin(name, version, description, author string) *BasePlugin {
	return &BasePlugin{
		name:        name,
		version:     version,
		description: description,
		author:      author,
	}
}

// Name 返回插件名称
func (p *BasePlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *BasePlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *BasePlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *BasePlugin) Author() string {
	return p.author
}

// Init 初始化插件（基础实现）
func (p *BasePlugin) Init() error {
	return nil
}

// Cleanup 清理资源（基础实现）
func (p *BasePlugin) Cleanup() error {
	return nil
}

// PluginManager 插件管理器接口
type PluginManager interface {
	// LoadPlugin 加载单个插件
	LoadPlugin(path string) (Plugin, error)
	
	// LoadPlugins 加载插件目录下的所有插件
	LoadPlugins(dir string, rg *gin.RouterGroup) error
	
	// GetPlugin 获取插件
	GetPlugin(name string) (Plugin, bool)
	
	// ListPlugins 列出所有已加载的插件
	ListPlugins() []Plugin
	
	// UnloadPlugin 卸载插件
	UnloadPlugin(name string) error
	
	// UnloadAll 卸载所有插件
	UnloadAll() error
}

// PluginInfo 插件信息结构体
type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Status      string `json:"status"`
	Path        string `json:"path"`
}

// PluginConfig 插件配置结构体
type PluginConfig struct {
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config"`
}

// PluginManifest 插件清单文件结构体
type PluginManifest struct {
	Name        string `json:"name" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Entry       string `json:"entry" binding:"required"`
	Enabled     bool   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
}