@echo off
setlocal enabledelayedexpansion

:: 切换到项目根目录
cd /d %~dp0..

:: 获取版本号
set TAG_NAME=%1

if "%TAG_NAME%"=="" (
    echo Usage: .\scripts\tag.bat ^<version^>
    exit /b 1
)

:: 确保在 master 分支
for /f "tokens=*" %%i in ('git rev-parse --abbrev-ref HEAD') do set CURRENT_BRANCH=%%i
if not "%CURRENT_BRANCH%"=="master" (
    echo Error: Must be on master branch to tag. Current branch: %CURRENT_BRANCH%
    exit /b 1
)

:: 检查是否有未提交的代码
set CHANGES=
for /f "tokens=*" %%i in ('git status --porcelain') do set CHANGES=%%i
if not "%CHANGES%"=="" (
    echo Error: You have uncommitted changes. Please commit or stash them first.
    exit /b 1
)

:: 拉取最新代码
echo Updating master branch...
git pull origin master
if %errorlevel% neq 0 exit /b %errorlevel%

:: 检查本地和远端是否一致
for /f "tokens=*" %%i in ('git rev-parse @') do set LOCAL=%%i
for /f "tokens=*" %%i in ('git rev-parse @{u}') do set REMOTE=%%i

if not "%LOCAL%"=="%REMOTE%" (
    echo Error: Local branch is not in sync with remote. Please push your changes first.
    exit /b 1
)

:: 创建并推送 tag
echo Creating tag %TAG_NAME%...
git tag -f "%TAG_NAME%"
if %errorlevel% neq 0 exit /b %errorlevel%

echo Pushing tag %TAG_NAME% to remote...
git push origin "%TAG_NAME%" -f
if %errorlevel% neq 0 exit /b %errorlevel%

echo Done! GitHub Action should be triggered shortly.
