#!/bin/bash

# 获取版本号
TAG_NAME=$1

if [ -z "$TAG_NAME" ]; then
    echo "Usage: ./scripts/tag.sh <version>"
    exit 1
fi

# 确保在 master 分支
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "master" ]; then
    echo "Error: Must be on master branch to tag. Current branch: $CURRENT_BRANCH"
    exit 1
fi

# 检查是否有未提交的代码
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: You have uncommitted changes. Please commit or stash them first."
    exit 1
fi

# 拉取最新代码
echo "Updating master branch..."
git pull origin master

# 检查本地和远端是否一致
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse @{u})

if [ "$LOCAL" != "$REMOTE" ]; then
    echo "Error: Local branch is not in sync with remote. Please push your changes first."
    exit 1
fi

# 创建并推送 tag
echo "Creating tag $TAG_NAME..."
git tag -f "$TAG_NAME"

echo "Pushing tag $TAG_NAME to remote..."
git push origin "$TAG_NAME" -f

echo "Done! GitHub Action should be triggered shortly."
