/*
Package services 提供服务管理功能

该包负责管理 Nginx、MySQL、PHP-CGI 服务的启动、停止和状态监控，
以及站点配置的管理和 Nginx 配置文件的动态生成。
*/
package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"text/template"
)

// ServiceStatus 服务状态
type ServiceStatus string

const (
	StatusStopped ServiceStatus = "stopped"
	StatusRunning ServiceStatus = "running"
	StatusError   ServiceStatus = "error"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Status  ServiceStatus `json:"status"`
	Port    string        `json:"port"`
}

// SiteConfig 站点配置
type SiteConfig struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Port    int    `json:"port"`
	Root    string `json:"root"`
	Enabled bool   `json:"enabled"`
	// URL 重写规则
	RewriteRules []RewriteRule `json:"rewriteRules"`
}

// RewriteRule URL 重写规则
type RewriteRule struct {
	ID          int    `json:"id"`
	Pattern     string `json:"pattern"`     // 匹配模式
	Destination string `json:"destination"` // 目标地址
	Type        string `json:"type"`        // 重写类型 (rewrite, redirect, proxy)
	Enabled     bool   `json:"enabled"`     // 是否启用
}

// ServiceManager 服务管理器
type ServiceManager struct {
	basePath   string
	nginxPath  string
	mysqlPath  string
	phpPath    string
	wwwPath    string
	nginxCmd   *exec.Cmd
	mysqlCmd   *exec.Cmd
	phpCmd     *exec.Cmd
	mu         sync.Mutex
	sites      []SiteConfig
}

// NewServiceManager 创建服务管理器
func NewServiceManager(basePath string) *ServiceManager {
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
	
	wwwPath := filepath.Join(basePath, "www")
	defaultRoot := strings.ReplaceAll(filepath.Join(wwwPath, "localhost"), "\\", "/")
	
	sm := &ServiceManager{
		basePath:  basePath,
		nginxPath: filepath.Join(basePath, "application", "nginx", "nginx-1.26.2"),
		mysqlPath: filepath.Join(basePath, "application", "mysql", "mysql-8.4.7-winx64"),
		phpPath:   filepath.Join(basePath, "application", "php", "php-8.1.33-nts-Win32-vs16-x64"),
		wwwPath:   wwwPath,
		sites: []SiteConfig{
			{ID: 1, Name: "localhost", Domain: "localhost", Port: 80, Root: defaultRoot, Enabled: true, RewriteRules: []RewriteRule{}},
		},
	}
	// 创建 www 目录
	os.MkdirAll(sm.wwwPath, 0755)
	os.MkdirAll(filepath.Join(sm.wwwPath, "localhost"), 0755)
	// 生成动态 nginx 配置
	sm.GenerateNginxConfig()
	return sm
}

// generateNginxConfig 动态生成 nginx 配置文件
func (sm *ServiceManager) generateNginxConfig() error {
	return sm.GenerateNginxConfig()
}

// GenerateNginxConfig 动态生成 nginx 配置文件
func (sm *ServiceManager) GenerateNginxConfig() error {
	confTemplate := `worker_processes  1;

events {
    worker_connections  1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    keepalive_timeout  65;
{{range .Sites}}{{if .Enabled}}
    server {
        listen       {{.Port}};
        server_name  {{.Domain}};

        root   {{.Root}};
        index  index.php index.html index.htm;

        # URL 重写规则
        {{range .RewriteRules}}{{if .Enabled}}
        {{if eq .Type "rewrite"}}rewrite ^{{.Pattern}}$ {{.Destination}} last;{{end}}
        {{if eq .Type "redirect"}}rewrite ^{{.Pattern}}$ {{.Destination}} permanent;{{end}}
        {{if eq .Type "proxy"}}location ~ ^{{.Pattern}}$ {
            proxy_pass {{.Destination}};
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }{{end}}
        {{end}}{{end}}

        location / {
            try_files $uri $uri/ /index.php?$query_string;
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }

        location ~ \.php$ {
            fastcgi_pass   127.0.0.1:9001;
            fastcgi_index  index.php;
            fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
            include        fastcgi_params;
        }

        location ~ /\.ht {
            deny  all;
        }
    }
{{end}}{{end}}
}
`
	tmpl, err := template.New("nginx").Parse(confTemplate)
	if err != nil {
		return err
	}

	confPath := filepath.Join(sm.nginxPath, "conf", "nginx.conf")
	file, err := os.Create(confPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Sites []SiteConfig
	}{
		Sites: sm.sites,
	}

	return tmpl.Execute(file, data)
}

// GetSites 获取所有站点配置
func (sm *ServiceManager) GetSites() []SiteConfig {
	return sm.sites
}

// UpdateSites 更新站点配置并重新生成 nginx 配置
func (sm *ServiceManager) UpdateSites(sites []SiteConfig) error {
	sm.sites = sites
	return sm.GenerateNginxConfig()
}

// GetAllServices 获取所有服务信息
func (sm *ServiceManager) GetAllServices() []ServiceInfo {
	return []ServiceInfo{
		{
			Name:    "Nginx",
			Version: "1.26.2",
			Status:  sm.getNginxStatus(),
			Port:    "80",
		},
		{
			Name:    "MySQL",
			Version: "8.4.7",
			Status:  sm.getMySQLStatus(),
			Port:    "3306",
		},
		{
			Name:    "PHP-CGI",
			Version: "8.1.33",
			Status:  sm.getPHPStatus(),
			Port:    "9001",
		},
	}
}

// StartNginx 启动 Nginx
func (sm *ServiceManager) StartNginx() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.isProcessRunning(sm.nginxCmd) {
		return nil
	}

	nginxExe := filepath.Join(sm.nginxPath, "nginx.exe")
	if _, err := os.Stat(nginxExe); os.IsNotExist(err) {
		return fmt.Errorf("nginx.exe 不存在: %s", nginxExe)
	}

	sm.nginxCmd = exec.Command(nginxExe)
	sm.nginxCmd.Dir = sm.nginxPath
	sm.nginxCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return sm.nginxCmd.Start()
}

// StopNginx 停止 Nginx
func (sm *ServiceManager) StopNginx() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 使用 taskkill 强制终止所有 nginx 进程
	cmd := exec.Command("taskkill", "/F", "/IM", "nginx.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	sm.nginxCmd = nil
	return nil
}

// StartMySQL 启动 MySQL
func (sm *ServiceManager) StartMySQL() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.isProcessRunning(sm.mysqlCmd) {
		return nil
	}

	mysqldExe := filepath.Join(sm.mysqlPath, "bin", "mysqld.exe")
	if _, err := os.Stat(mysqldExe); os.IsNotExist(err) {
		return fmt.Errorf("mysqld.exe 不存在: %s", mysqldExe)
	}

	dataDir := filepath.Join(sm.mysqlPath, "data")
	sm.mysqlCmd = exec.Command(mysqldExe, 
		"--basedir="+sm.mysqlPath,
		"--datadir="+dataDir,
		"--port=3306",
	)
	sm.mysqlCmd.Dir = sm.mysqlPath
	sm.mysqlCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return sm.mysqlCmd.Start()
}

// StopMySQL 停止 MySQL
func (sm *ServiceManager) StopMySQL() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 使用 taskkill 强制终止 mysqld 进程
	cmd := exec.Command("taskkill", "/F", "/IM", "mysqld.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	sm.mysqlCmd = nil
	return nil
}

// StartPHP 启动 PHP-CGI
func (sm *ServiceManager) StartPHP() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.isProcessRunning(sm.phpCmd) {
		return nil
	}

	phpCgiExe := filepath.Join(sm.phpPath, "php-cgi.exe")
	if _, err := os.Stat(phpCgiExe); os.IsNotExist(err) {
		return fmt.Errorf("php-cgi.exe 不存在: %s", phpCgiExe)
	}

	sm.phpCmd = exec.Command(phpCgiExe, "-b", "127.0.0.1:9001")
	sm.phpCmd.Dir = sm.phpPath
	sm.phpCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return sm.phpCmd.Start()
}

// StopPHP 停止 PHP-CGI
func (sm *ServiceManager) StopPHP() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 使用 taskkill 强制终止 php-cgi 进程
	cmd := exec.Command("taskkill", "/F", "/IM", "php-cgi.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	sm.phpCmd = nil
	return nil
}

// StartAll 启动所有服务
func (sm *ServiceManager) StartAll() error {
	if err := sm.StartPHP(); err != nil {
		return fmt.Errorf("启动 PHP 失败: %v", err)
	}
	if err := sm.StartMySQL(); err != nil {
		return fmt.Errorf("启动 MySQL 失败: %v", err)
	}
	if err := sm.StartNginx(); err != nil {
		return fmt.Errorf("启动 Nginx 失败: %v", err)
	}
	return nil
}

// StopAll 停止所有服务
func (sm *ServiceManager) StopAll() error {
	sm.StopNginx()
	sm.StopMySQL()
	sm.StopPHP()
	return nil
}

// RestartAll 重启所有服务
func (sm *ServiceManager) RestartAll() error {
	sm.StopAll()
	return sm.StartAll()
}

// isProcessRunning 检查进程是否运行
func (sm *ServiceManager) isProcessRunning(cmd *exec.Cmd) bool {
	if cmd == nil || cmd.Process == nil {
		return false
	}
	// 尝试发送信号检测进程是否存在
	err := cmd.Process.Signal(syscall.Signal(0))
	return err == nil
}

// getNginxStatus 获取 Nginx 状态
func (sm *ServiceManager) getNginxStatus() ServiceStatus {
	// 检查是否有 nginx 进程在运行
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq nginx.exe", "/FO", "CSV", "/NH")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return StatusStopped
	}
	if strings.Contains(string(output), "nginx.exe") {
		return StatusRunning
	}
	return StatusStopped
}

// getMySQLStatus 获取 MySQL 状态
func (sm *ServiceManager) getMySQLStatus() ServiceStatus {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq mysqld.exe", "/FO", "CSV", "/NH")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return StatusStopped
	}
	if strings.Contains(string(output), "mysqld.exe") {
		return StatusRunning
	}
	return StatusStopped
}

// getPHPStatus 获取 PHP-CGI 状态
func (sm *ServiceManager) getPHPStatus() ServiceStatus {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq php-cgi.exe", "/FO", "CSV", "/NH")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return StatusStopped
	}
	if strings.Contains(string(output), "php-cgi.exe") {
		return StatusRunning
	}
	return StatusStopped
}

// GetWWWPath 获取网站根目录
func (sm *ServiceManager) GetWWWPath() string {
	return sm.wwwPath
}

// GetBasePath 获取基础路径
func (sm *ServiceManager) GetBasePath() string {
	return sm.basePath
}
