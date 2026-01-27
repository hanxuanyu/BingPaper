@echo off
setlocal enabledelayedexpansion

:: 切换到项目根目录
cd /d %~dp0..

set APP_NAME=BingPaper
set OUTPUT_DIR=output

echo 开始构建前端...
cd webapp
call npm install
call npm run build
if %errorlevel% neq 0 (
    echo 前端构建失败
    exit /b %errorlevel%
)
cd ..

echo 开始构建 %APP_NAME% 多平台二进制文件...

if exist %OUTPUT_DIR% rd /s /q %OUTPUT_DIR%
mkdir %OUTPUT_DIR%

set PLATFORMS=linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64

for %%p in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%p") do (
        set GOOS=%%a
        set GOARCH=%%b
        
        set OUTPUT_NAME=%APP_NAME%-%%a-%%b
        set BINARY_NAME=!OUTPUT_NAME!
        if "%%a"=="windows" set BINARY_NAME=!OUTPUT_NAME!.exe
        
        echo 正在编译 %%a/%%b...
        
        set PACKAGE_DIR=%OUTPUT_DIR%\!OUTPUT_NAME!
        if not exist !PACKAGE_DIR! mkdir !PACKAGE_DIR!
        
        set GOOS=%%a
        set GOARCH=%%b
        set CGO_ENABLED=0
        go build -ldflags="-s -w" -o !PACKAGE_DIR!\!BINARY_NAME! main.go
        
        if !errorlevel! equ 0 (
            echo   %%a/%%b 编译成功
            
            xcopy /e /i /y web !PACKAGE_DIR!\web >nul
            copy /y config.example.yaml !PACKAGE_DIR!\ >nul
            copy /y README.md !PACKAGE_DIR!\ >nul
            
            pushd %OUTPUT_DIR%
            tar -czf !OUTPUT_NAME!.tar.gz !OUTPUT_NAME!
            rd /s /q !OUTPUT_NAME!
            popd
            
            echo   %%a/%%b 打包完成: !OUTPUT_NAME!.tar.gz
        ) else (
            echo   %%a/%%b 编译失败
            if exist !PACKAGE_DIR! rd /s /q !PACKAGE_DIR!
        )
    )
)

echo ----------------------------------------
echo 多平台打包完成！输出目录: %OUTPUT_DIR%
dir /s /b %OUTPUT_DIR%
pause
