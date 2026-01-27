# 切换到项目根目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
if ($ScriptDir) { Set-Location (Join-Path $ScriptDir "..") }

$AppName = "BingPaper"
$OutputDir = "output"

Write-Host "开始构建前端..."
Push-Location webapp
npm install
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "前端构建失败" -ForegroundColor Red
    Pop-Location
    exit $LASTEXITCODE
}
Pop-Location

Write-Host "开始构建 $AppName 多平台二进制文件..."

if (Test-Path $OutputDir) {
    Remove-Item -Recurse -Force $OutputDir
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null

$Platforms = @(
    "linux/amd64",
    "linux/arm64",
    "windows/amd64",
    "windows/arm64",
    "darwin/amd64",
    "darwin/arm64"
)

foreach ($Platform in $Platforms) {
    $parts = $Platform.Split("/")
    $OS = $parts[0]
    $Arch = $parts[1]
    
    $OutputName = "$AppName-$OS-$Arch"
    $BinaryName = $AppName
    if ($OS -eq "windows") {
        $BinaryName = "$AppName.exe"
    }
    
    Write-Host "正在编译 $OS/$Arch..."
    
    $PackageDir = Join-Path $OutputDir $OutputName
    if (-not (Test-Path $PackageDir)) {
        New-Item -ItemType Directory -Path $PackageDir | Out-Null
    }
    
    $env:GOOS = $OS
    $env:GOARCH = $Arch
    $env:CGO_ENABLED = "0"
    go build -ldflags="-s -w" -o (Join-Path $PackageDir $BinaryName) main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  $OS/$Arch 编译成功"
        
        Copy-Item -Recurse "web" $PackageDir\
        Copy-Item "config.example.yaml" $PackageDir\
        Copy-Item "README.md" $PackageDir\
        
        $CurrentDir = Get-Location
        Set-Location $PackageDir
        tar -czf "../${OutputName}.tar.gz" .
        Set-Location $CurrentDir
        Remove-Item -Recurse -Force $PackageDir
        
        Write-Host "  $OS/$Arch 打包完成: ${OutputName}.tar.gz"
    } else {
        Write-Host "  $OS/$Arch 编译失败"
        if (Test-Path $PackageDir) {
            Remove-Item -Recurse -Force $PackageDir
        }
    }
}

Write-Host "----------------------------------------"
Write-Host "多平台打包完成！输出目录: $OutputDir"
Get-ChildItem -Recurse $OutputDir
