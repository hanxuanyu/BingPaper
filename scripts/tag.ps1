# 切换到项目根目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
if ($ScriptDir) { Set-Location (Join-Path $ScriptDir "..") }

# 获取版本号
$TagName = $args[0]

if (-not $TagName) {
    Write-Host "Usage: .\scripts\tag.ps1 <version>"
    exit 1
}

# 确保在 master 分支
$CurrentBranch = git rev-parse --abbrev-ref HEAD
if ($CurrentBranch -ne "master") {
    Write-Host "Error: Must be on master branch to tag. Current branch: $CurrentBranch" -ForegroundColor Red
    exit 1
}

# 检查是否有未提交的代码
$Changes = git status --porcelain
if ($Changes) {
    Write-Host "Error: You have uncommitted changes. Please commit or stash them first." -ForegroundColor Red
    exit 1
}

# 拉取最新代码
Write-Host "Updating master branch..."
git pull origin master
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

# 检查本地和远端是否一致
$Local = git rev-parse "@"
$Remote = git rev-parse "@{u}"

if ($Local -ne $Remote) {
    Write-Host "Error: Local branch is not in sync with remote. Please push your changes first." -ForegroundColor Red
    exit 1
}

# 创建并推送 tag
Write-Host "Creating tag $TagName..."
git tag -f "$TagName"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Pushing tag $TagName to remote..."
git push origin "$TagName" -f
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Done! GitHub Action should be triggered shortly."
