package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"wampgo/internal/services"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows/registry"
)

// AppSettings 应用设置
type AppSettings struct {
	AutoStart       bool   `json:"autoStart"`       // 开机启动
	NginxPort       int    `json:"nginxPort"`       // Nginx 端口
	PHPPort         int    `json:"phpPort"`         // PHP 端口
	MySQLPort       int    `json:"mysqlPort"`       // MySQL 端口
}

// App struct
type App struct {
	ctx      context.Context
	manager  *services.ServiceManager
	settings AppSettings
	basePath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	// 获取程序所在目录作为基础路径
	exePath, _ := os.Executable()
	basePath := filepath.Dir(exePath)
	
	// 检查是否存在 application 目录，如果不存在则向上查找
	if _, err := os.Stat(filepath.Join(basePath, "application")); os.IsNotExist(err) {
		// 在构建环境中，可能需要向上查找两级目录
		parentDir := filepath.Dir(basePath)
		if _, err := os.Stat(filepath.Join(parentDir, "application")); err == nil {
			basePath = parentDir
		} else {
			// 再向上查找一级
			grandParentDir := filepath.Dir(parentDir)
			if _, err := os.Stat(filepath.Join(grandParentDir, "application")); err == nil {
				basePath = grandParentDir
			}
		}
	}

	app := &App{
		manager:  services.NewServiceManager(basePath),
		basePath: basePath,
		settings: AppSettings{
			AutoStart:       false,
			NginxPort:       80,
			PHPPort:         9001,
			MySQLPort:       3306,
		},
	}
	app.loadSettings()
	return app
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 不再默认启动所有服务
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	a.manager.StopAll()
}

// GetServices 获取所有服务状态
func (a *App) GetServices() []services.ServiceInfo {
	return a.manager.GetAllServices()
}

// StartService 启动指定服务
func (a *App) StartService(name string) error {
	switch name {
	case "Nginx":
		return a.manager.StartNginx()
	case "MySQL":
		return a.manager.StartMySQL()
	case "PHP-CGI":
		return a.manager.StartPHP()
	}
	return nil
}

// StopService 停止指定服务
func (a *App) StopService(name string) error {
	switch name {
	case "Nginx":
		return a.manager.StopNginx()
	case "MySQL":
		return a.manager.StopMySQL()
	case "PHP-CGI":
		return a.manager.StopPHP()
	}
	return nil
}

// StartAll 启动所有服务
func (a *App) StartAll() error {
	return a.manager.StartAll()
}

// StopAll 停止所有服务
func (a *App) StopAll() error {
	return a.manager.StopAll()
}

// RestartAll 重启所有服务
func (a *App) RestartAll() error {
	return a.manager.RestartAll()
}

// OpenWWWFolder 打开网站目录
func (a *App) OpenWWWFolder() error {
	wwwPath := a.manager.GetWWWPath()
	cmd := exec.Command("explorer", wwwPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}

// OpenBrowser 打开浏览器
func (a *App) OpenBrowser() {
	runtime.BrowserOpenURL(a.ctx, "http://localhost")
}

// GetWWWPath 获取网站目录路径
func (a *App) GetWWWPath() string {
	wwwPath := a.manager.GetWWWPath()
	// 转换为正斜杠格式
	return strings.ReplaceAll(wwwPath, "\\", "/")
}

// GetDefaultSiteRoot 获取默认站点根目录
func (a *App) GetDefaultSiteRoot() string {
	wwwPath := a.manager.GetWWWPath()
	defaultRoot := filepath.Join(wwwPath, "localhost")
	return strings.ReplaceAll(defaultRoot, "\\", "/")
}

// SelectDirectory 选择目录
func (a *App) SelectDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "选择站点根目录",
		DefaultDirectory: a.manager.GetWWWPath(),
	})
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(dir, "\\", "/"), nil
}

// GetBasePath 获取基础路径
func (a *App) GetBasePath() string {
	return a.manager.GetBasePath()
}

// GetSites 获取站点列表
func (a *App) GetSites() []services.SiteConfig {
	return a.manager.GetSites()
}

// UpdateSites 更新站点配置
func (a *App) UpdateSites(sites []services.SiteConfig) error {
	// 创建站点目录和默认文件
	for _, site := range sites {
		if err := a.createSiteDirectory(site.Root); err != nil {
			return fmt.Errorf("创建站点目录失败: %v", err)
		}
	}
	
	err := a.manager.UpdateSites(sites)
	if err != nil {
		return err
	}
	// 重启 Nginx 使配置生效
	a.manager.StopNginx()
	return a.manager.StartNginx()
}

// createSiteDirectory 创建站点目录和默认 index.php
func (a *App) createSiteDirectory(siteRoot string) error {
	// 检查目录是否存在
	if _, err := os.Stat(siteRoot); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(siteRoot, 0755); err != nil {
			return err
		}
	}
	
	// 检查目录是否为空
	entries, err := os.ReadDir(siteRoot)
	if err != nil {
		return err
	}
	
	// 如果目录为空，创建默认 index.php
	if len(entries) == 0 {
		indexPath := filepath.Join(siteRoot, "index.php")
		phpCode := "<?php\necho phpinfo();\n"
		if err := os.WriteFile(indexPath, []byte(phpCode), 0644); err != nil {
			return err
		}
	}
	
	return nil
}

// OpenSiteInBrowser 在默认浏览器中打开站点
func (a *App) OpenSiteInBrowser(domain string, port int) {
	var url string
	if port == 80 {
		url = "http://" + domain
	} else {
		url = fmt.Sprintf("http://%s:%d", domain, port)
	}
	runtime.BrowserOpenURL(a.ctx, url)
}

// OpenSiteFolder 打开站点根目录
func (a *App) OpenSiteFolder(path string) {
	path = strings.ReplaceAll(path, "/", "\\")
	exec.Command("cmd", "/c", "start", "explorer", path).Start()
}

// ==================== 设置相关方法 ====================

// getSettingsPath 获取设置文件路径
func (a *App) getSettingsPath() string {
	return filepath.Join(a.basePath, "settings.json")
}

// loadSettings 加载设置
func (a *App) loadSettings() {
	data, err := os.ReadFile(a.getSettingsPath())
	if err != nil {
		return
	}
	json.Unmarshal(data, &a.settings)
}

// saveSettings 保存设置
func (a *App) saveSettings() error {
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.getSettingsPath(), data, 0644)
}

// GetSettings 获取设置
func (a *App) GetSettings() AppSettings {
	return a.settings
}

// UpdateSettings 更新设置
func (a *App) UpdateSettings(settings AppSettings) error {
	// 处理开机启动
	if settings.AutoStart != a.settings.AutoStart {
		if err := a.setAutoStart(settings.AutoStart); err != nil {
			return err
		}
	}
	a.settings = settings
	return a.saveSettings()
}

// setAutoStart 设置开机启动
func (a *App) setAutoStart(enable bool) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	appName := "Sanctuary"
	if enable {
		exePath, _ := os.Executable()
		return key.SetStringValue(appName, exePath)
	} else {
		key.DeleteValue(appName)
		return nil
	}
}

// SetMinimizeToTray 设置最小化到托盘
func (a *App) SetMinimizeToTray(minimize bool) error {
	// 此功能已被移除
	return nil
}
