/*
Sanctuary - PHP 集成开发环境

这是一个基于 Go + Wails 构建的 Windows 桌面应用，
提供 Nginx + MySQL + PHP 的一键式管理功能。

作者: 子曰亦君
QQ: 15593838
*/
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// 嵌入前端构建资源
//go:embed all:frontend/dist
var assets embed.FS

// main 程序入口
func main() {
	// 创建应用实例
	app := NewApp()

	// 配置并启动 Wails 应用
	err := wails.Run(&options.App{
		Title:  "Sanctuary",              // 窗口标题
		Width:  1024,                      // 窗口宽度
		Height: 768,                       // 窗口高度
		AssetServer: &assetserver.Options{ // 静态资源服务
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1}, // 背景色
		OnStartup:        app.startup,                              // 启动回调
		OnShutdown:       app.shutdown,                             // 关闭回调
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    true,
		},
		Bind: []interface{}{ // 绑定到前端的方法
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
