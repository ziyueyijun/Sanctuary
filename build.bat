@echo off
title Sanctuary Build
cd /d %~dp0
echo Building Sanctuary...
wails build
echo.
echo Build complete! Output: build\bin\sanctuary.exe
pause
