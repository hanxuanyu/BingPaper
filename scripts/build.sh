#!/bin/bash

# 获取脚本所在目录并切换到项目根目录
cd "$(dirname "$0")/.."

# 设置变量
APP_NAME="BingPaper"
OUTPUT_DIR="output"

# 定义目标平台
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

# 需要包含的额外文件/目录
EXTRA_FILES=("web" "config.example.yaml" "README.md")

echo "开始构建前端..."
cd webapp
npm install
npm run build
cd ..

echo "开始构建 $APP_NAME 多平台二进制文件..."

# 清理 output 目录
if [ -d "$OUTPUT_DIR" ]; then
    echo "正在清理 $OUTPUT_DIR 目录..."
    rm -rf "$OUTPUT_DIR"
fi

mkdir -p "$OUTPUT_DIR"

# 循环编译各平台
for PLATFORM in "${PLATFORMS[@]}"; do
    # 分离 OS 和 ARCH
    OS=$(echo $PLATFORM | cut -d'/' -f1)
    ARCH=$(echo $PLATFORM | cut -d'/' -f2)
    
    # 设置输出名称
    OUTPUT_NAME="${APP_NAME}-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${OUTPUT_NAME}.exe"
    else
        BINARY_NAME="${OUTPUT_NAME}"
    fi
    
    echo "正在编译 ${OS}/${ARCH}..."
    
    # 创建临时打包目录
    PACKAGE_DIR="${OUTPUT_DIR}/${OUTPUT_NAME}"
    mkdir -p "$PACKAGE_DIR"
    
    # 现在已移除 CGO 依赖，使用 CGO_ENABLED=0 以支持轻松的跨平台编译
    # 增加 -ldflags="-s -w" 以减少二进制体积
    GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 go build -ldflags="-s -w" -o "${PACKAGE_DIR}/${BINARY_NAME}" main.go
    
    if [ $? -eq 0 ]; then
        echo "  ${OS}/${ARCH} 编译成功"
        
        # 复制额外文件
        for file in "${EXTRA_FILES[@]}"; do
            if [ -e "$file" ]; then
                cp -r "$file" "$PACKAGE_DIR/"
            fi
        done
        
        # 压缩为 tar.gz
        tar -czf "${OUTPUT_DIR}/${OUTPUT_NAME}.tar.gz" -C "${OUTPUT_DIR}" "${OUTPUT_NAME}"
        
        # 删除临时打包目录
        rm -rf "$PACKAGE_DIR"
        
        echo "  ${OS}/${ARCH} 打包完成: ${OUTPUT_NAME}.tar.gz"
    else
        echo "  ${OS}/${ARCH} 编译失败"
        # 编译失败时清理临时目录
        rm -rf "$PACKAGE_DIR"
    fi
done

echo "----------------------------------------"
echo "多平台打包完成！输出目录: $OUTPUT_DIR"
ls -R "$OUTPUT_DIR"
