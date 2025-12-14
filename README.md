# Sanctuary

> PHP 集成开发环境 - 基于 Go + Wails 构建的 Windows 桌面应用

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8.svg)
![Wails](https://img.shields.io/badge/Wails-v2-purple.svg)

## 功能特性

- **服务管理**: 一键启动/停止 Nginx、MySQL、PHP-CGI 服务
- **站点管理**: 添加、编辑、删除站点配置，自动生成 Nginx 配置
- **URL 重写**: 支持直接粘贴 Nginx 格式重写规则
- **快速访问**: 点击域名即可在浏览器中打开站点
- **目录管理**: 点击根目录即可打开文件夹
- **自动配置**: 新建站点时自动创建目录和默认 index.php

## 技术栈

| 层级 | 技术 |
|-----|------|
| 后端 | Go 1.23 |
| 框架 | Wails v2 |
| 前端 | Vue 3 + Vite |
| UI库 | Arco Design Vue |

## 集成服务

| 服务 | 版本 | 端口 |
|-----|------|-----|
| Nginx | 1.26.2 | 80 |
| MySQL | 8.4.7 | 3306 |
| PHP-CGI | 8.1.33 | 9001 |

## 目录结构

```
sanctuary/
├── application/          # 服务程序目录
│   ├── nginx/           # Nginx
│   ├── mysql/           # MySQL
│   └── php/             # PHP
├── www/                  # 网站根目录
│   └── localhost/       # 默认站点
├── frontend/            # 前端源码
│   └── src/
│       ├── App.vue      # 主组件
│       └── main.js      # 入口文件
├── internal/            # 内部包
│   └── services/        # 服务管理模块
├── main.go              # 程序入口
├── app.go               # 应用逻辑
├── wails.json           # Wails 配置
├── dev.bat              # 开发脚本
└── build.bat            # 构建脚本
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 16+
- Wails CLI v2

### 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发模式

```bash
# 方式一：使用脚本
dev.bat

# 方式二：使用命令
wails dev
```

### 构建发布

```bash
# 方式一：使用脚本
build.bat

# 方式二：使用命令
wails build
```

构建产物位于 `build/bin/sanctuary.exe`

## 使用说明

### 服务管理

1. 点击 **启动全部** 一键启动所有服务
2. 点击 **停止全部** 一键停止所有服务
3. 点击 **重启全部** 重启所有服务
4. 也可单独控制每个服务的启动/停止

### 站点管理

1. 点击 **添加站点** 创建新站点
2. 输入站点名称和域名
3. 域名会自动生成对应的根目录路径
4. 点击 **保存** 后自动创建目录和默认文件
5. 在 hosts 文件中添加域名映射 (如 `127.0.0.1 mysite.local`)

### URL 重写规则

1. 在站点编辑界面切换到 **URL重写** 标签页
2. 直接粘贴 Nginx 格式的重写规则：
   ```
   rewrite ^/api/(.*)$ /api.php?path=$1 last;
   rewrite ^/admin$ /admin.php permanent;
   ```
3. 点击 **解析规则** 按钮自动转换为可视化规则
4. 可对个别规则进行启用/禁用或微调
5. 保存配置后自动应用到 Nginx

### 快捷操作

- **点击域名**: 在浏览器中打开站点
- **点击根目录**: 打开站点文件夹

## 配置说明

### MySQL 默认账户

- 用户名: `root`
- 密码: `root`

### Nginx 配置

站点配置文件自动生成，位于:
```
application/nginx/nginx-1.26.2/conf/nginx.conf
```

## 开发者

- **作者**: 子曰亦君
- **QQ**: 15593838

## 许可证

MIT License