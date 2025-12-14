@echo off
chcp 65001 >nul
echo ========================================
echo Sanctuary 打包脚本
echo ========================================
echo.

echo [1/4] 正在清理旧的构建文件...
if exist build rmdir /s /q build
if exist dist rmdir /s /q dist

echo.
echo [2/4] 正在构建 Sanctuary 应用...
wails build
if %errorlevel% neq 0 (
    echo [错误] 构建失败！
    pause
    exit /b 1
)

echo.
echo [3/4] 正在创建分发目录...
mkdir dist
mkdir dist\application
mkdir dist\www

echo.
echo [4/4] 正在复制文件...

REM 复制可执行文件
copy "build\bin\sanctuary.exe" "dist\" >nul
if %errorlevel% neq 0 (
    echo [错误] 复制可执行文件失败！
    pause
    exit /b 1
)

REM 复制服务程序目录
echo 正在复制 application 目录...
xcopy application dist\application /E /I /H /Y >nul
if %errorlevel% neq 0 (
    echo [警告] 复制 application 目录时出现问题
)

REM 复制网站根目录
echo 正在复制 www 目录...
xcopy www dist\www /E /I /H /Y >nul
if %errorlevel% neq 0 (
    echo [警告] 复制 www 目录时出现问题
)

REM 复制 README 和其他文档
if exist README.md copy README.md dist\ >nul
if exist LICENSE copy LICENSE dist\ >nul

echo.
echo ========================================
echo 打包完成！
echo 分发文件位于 dist 目录中
echo ========================================
echo 包含文件：
echo - sanctuary.exe (主程序)
echo - application/ (服务程序目录)
echo - www/ (网站根目录)
echo - README.md (说明文档)
echo - LICENSE (许可证文件)
echo ========================================
echo.
echo 提示：将 dist 目录压缩即可分发给用户
echo ========================================
pause