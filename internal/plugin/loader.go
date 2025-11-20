package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/gin-gonic/gin"
)

// Loader 插件加载器
type Loader struct {
	plugins map[string]Plugin
	mutex   sync.RWMutex
}

// NewLoader 创建新的插件加载器
func NewLoader() *Loader {
	return &Loader{
		plugins: make(map[string]Plugin),
	}
}

// LoadPlugin 加载单个插件
func (l *Loader) LoadPlugin(path string) (Plugin, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 读取插件清单文件
	manifestPath := filepath.Join(path, "manifest.json")
	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest PluginManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest file: %w", err)
	}

	// 检查插件是否已加载
	if _, exists := l.plugins[manifest.Name]; exists {
		return nil, fmt.Errorf("plugin %s already loaded", manifest.Name)
	}

	// 检查插件是否启用
	if !manifest.Enabled {
		return nil, fmt.Errorf("plugin %s is disabled", manifest.Name)
	}

	// 加载插件 .so 文件
	pluginPath := filepath.Join(path, manifest.Entry)
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin file not found: %s", pluginPath)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// 查找 NewPlugin 函数
	newPluginSymbol, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup NewPlugin function: %w", err)
	}

	// 类型断言
	newPlugin, ok := newPluginSymbol.(func() Plugin)
	if !ok {
		return nil, fmt.Errorf("invalid NewPlugin function signature")
	}

	// 创建插件实例
	pluginInstance := newPlugin()

	// 初始化插件
	if err := pluginInstance.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize plugin: %w", err)
	}

	// 存储插件
	l.plugins[manifest.Name] = pluginInstance

	log.Printf("Plugin %s (v%s) loaded successfully", manifest.Name, manifest.Version)

	return pluginInstance, nil
}

// LoadPlugins 加载插件目录下的所有插件
func (l *Loader) LoadPlugins(dir string, rg *gin.RouterGroup) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Plugin directory does not exist: %s", dir)
		return nil
	}

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(dir, entry.Name())
		pluginInstance, err := l.LoadPlugin(pluginPath)
		if err != nil {
			log.Printf("Failed to load plugin %s: %v", entry.Name(), err)
			continue
		}

		// 注册插件路由
		pluginGroup := rg.Group("/plugin/" + pluginInstance.Name())
		pluginInstance.Routes(pluginGroup)
	}

	return nil
}

// GetPlugin 获取插件
func (l *Loader) GetPlugin(name string) (Plugin, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	pluginInstance, exists := l.plugins[name]
	return pluginInstance, exists
}

// ListPlugins 列出所有已加载的插件
func (l *Loader) ListPlugins() []Plugin {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var pluginList []Plugin
	for _, pluginInstance := range l.plugins {
		pluginList = append(pluginList, pluginInstance)
	}

	return pluginList
}

// GetPluginInfos 获取插件信息列表
func (l *Loader) GetPluginInfos() []PluginInfo {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var pluginInfos []PluginInfo
	for name, pluginInstance := range l.plugins {
		pluginInfo := PluginInfo{
			Name:        pluginInstance.Name(),
			Version:     pluginInstance.Version(),
			Description: pluginInstance.Description(),
			Author:      pluginInstance.Author(),
			Status:      "loaded",
		}
		pluginInfos = append(pluginInfos, pluginInfo)
	}

	return pluginInfos
}

// UnloadPlugin 卸载插件
func (l *Loader) UnloadPlugin(name string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	pluginInstance, exists := l.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// 清理插件资源
	if err := pluginInstance.Cleanup(); err != nil {
		log.Printf("Failed to cleanup plugin %s: %v", name, err)
	}

	// 从插件列表中移除
	delete(l.plugins, name)

	log.Printf("Plugin %s unloaded successfully", name)

	return nil
}

// UnloadAll 卸载所有插件
func (l *Loader) UnloadAll() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for name, pluginInstance := range l.plugins {
		if err := pluginInstance.Cleanup(); err != nil {
			log.Printf("Failed to cleanup plugin %s: %v", name, err)
		}
	}

	// 清空插件列表
	l.plugins = make(map[string]Plugin)

	log.Println("All plugins unloaded")

	return nil
}

// ReloadPlugin 重新加载插件
func (l *Loader) ReloadPlugin(path string) error {
	// 读取清单文件获取插件名称
	manifestPath := filepath.Join(path, "manifest.json")
	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest PluginManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest file: %w", err)
	}

	// 先卸载旧版本
	if _, exists := l.GetPlugin(manifest.Name); exists {
		if err := l.UnloadPlugin(manifest.Name); err != nil {
			return fmt.Errorf("failed to unload old plugin version: %w", err)
		}
	}

	// 加载新版本
	_, err = l.LoadPlugin(path)
	return err
}

// ScanPlugins 扫描插件目录但不加载
func (l *Loader) ScanPlugins(dir string) ([]PluginInfo, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin directory does not exist: %s", dir)
	}

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %w", err)
	}

	var pluginInfos []PluginInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(dir, entry.Name())
		pluginInfo, err := l.scanPlugin(pluginPath)
		if err != nil {
			log.Printf("Failed to scan plugin %s: %v", entry.Name(), err)
			continue
		}

		pluginInfos = append(pluginInfos, pluginInfo)
	}

	return pluginInfos, nil
}

// scanPlugin 扫描单个插件
func (l *Loader) scanPlugin(path string) (PluginInfo, error) {
	// 读取清单文件
	manifestPath := filepath.Join(path, "manifest.json")
	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return PluginInfo{}, fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest PluginManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return PluginInfo{}, fmt.Errorf("failed to parse manifest file: %w", err)
	}

	// 检查插件是否已加载
	status := "not_loaded"
	if _, exists := l.GetPlugin(manifest.Name); exists {
		status = "loaded"
	}

	return PluginInfo{
		Name:        manifest.Name,
		Version:     manifest.Version,
		Description: manifest.Description,
		Author:      manifest.Author,
		Status:      status,
		Path:        path,
	}, nil
}