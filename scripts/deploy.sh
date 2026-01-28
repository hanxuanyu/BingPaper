#!/bin/bash

# 定时任务示例 (每30分钟执行一次):
# */30 * * * * /path/to/project/scripts/deploy.sh >> /path/to/project/scripts/deploy_cron.log 2>&1

# 项目根目录
# 假设脚本位于 scripts 目录下
PROJECT_DIR=$(cd $(dirname $0)/.. && pwd)
cd $PROJECT_DIR

# 日志文件
LOG_FILE="$PROJECT_DIR/scripts/deploy.log"

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# 确保在项目根目录
if [ ! -f "docker-compose.yml" ]; then
    log "错误: 未能在 $PROJECT_DIR 找到 docker-compose.yml，请确保脚本位置正确。"
    exit 1
fi

log "开始检查更新..."

# 获取远程代码
git fetch origin

# 获取当前分支名
BRANCH=$(git rev-parse --abbrev-ref HEAD)

# 检查本地分支是否有上游分支
UPSTREAM=$(git rev-parse --abbrev-ref @{u} 2>/dev/null)
if [ -z "$UPSTREAM" ]; then
    log "错误: 当前分支 $BRANCH 没有设置上游分支，无法自动对比更新。"
    exit 1
fi

# 检查本地是否落后于远程
LOCAL=$(git rev-parse HEAD)
REMOTE=$(git rev-parse "$UPSTREAM")

if [ "$LOCAL" = "$REMOTE" ]; then
    log "代码已是最新 ($LOCAL)，无需更新。"
    exit 0
fi

log "检测到远程变更 ($LOCAL -> $REMOTE)，准备开始升级..."

# 检查是否有本地修改
HAS_CHANGES=$(git status --porcelain)

if [ -n "$HAS_CHANGES" ]; then
    log "检测到本地修改，正在暂存以保留个性化配置..."
    git stash
    STASHED=true
else
    STASHED=false
fi

# 拉取新代码
log "正在拉取远程代码 ($BRANCH)..."
if git pull origin "$BRANCH"; then
    log "代码拉取成功。"
else
    log "错误: 代码拉取失败。"
    if [ "$STASHED" = true ]; then
        git stash pop
    fi
    exit 1
fi

# 恢复本地修改
if [ "$STASHED" = true ]; then
    log "正在恢复本地修改..."
    if git stash pop; then
        log "本地修改已成功恢复。"
    else
        log "警告: 恢复本地修改时发生冲突，请手动检查 docker-compose.yml 等文件。"
        # 即使冲突也尝试继续，或者你可以选择在此退出
    fi
fi

# 确定 docker-compose 命令
DOCKER_COMPOSE_BIN="docker-compose"
if ! command -v $DOCKER_COMPOSE_BIN &> /dev/null; then
    DOCKER_COMPOSE_BIN="docker compose"
fi

# 执行 docker-compose 部署
log "正在执行 $DOCKER_COMPOSE_BIN 部署..."
if $DOCKER_COMPOSE_BIN up -d --build; then
    log "服务升级成功！"
    # 清理无用镜像（可选）
    docker image prune -f
else
    log "错误: $DOCKER_COMPOSE_BIN 部署失败。"
    exit 1
fi

log "部署任务完成。"
