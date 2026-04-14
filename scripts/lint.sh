#!/usr/bin/env bash

set -euo pipefail

source ./scripts/lib.sh

# ========================
# 工具安装函数
# ========================
function install_if_not_exist() {
  local TOOL_NAME=$1
  local INSTALL_URL=$2

  if command -v "$TOOL_NAME" &> /dev/null; then
    log_callout "$TOOL_NAME is already installed."
  else
    log_cmd "$TOOL_NAME is not installed. Installing..."
    run go install "$INSTALL_URL"
  fi
}

install_if_not_exist go-cleanarch github.com/roblaszczak/go-cleanarch@latest
install_if_not_exist goimports golang.org/x/tools/cmd/goimports@latest

# ========================
# golangci-lint 安装（按需更新版本）
# ========================
readonly LINT_VERSION="v1.64.0"
NEED_INSTALL=false

if command -v golangci-lint >/dev/null 2>&1; then
  # 獲取當前版本並統一格式化為帶 v 前綴 (例如 "v1.59.1")
  CURRENT_VERSION="v$(golangci-lint --version | awk '{print $4}' | sed 's/^v//')"
  
  if [ "$CURRENT_VERSION" == "$LINT_VERSION" ]; then
    log_callout "golangci-lint $CURRENT_VERSION is already installed."
  else
    log_cmd "golangci-lint version mismatch (current: $CURRENT_VERSION, need: $LINT_VERSION). Reinstalling..."
    NEED_INSTALL=true
  fi
else
  log_cmd "golangci-lint is not installed. Installing..."
  NEED_INSTALL=true
fi

if [ "$NEED_INSTALL" == true ]; then
  log_info "Installing golangci-lint $LINT_VERSION ..."
  # 確保刪除舊版，避免路徑權限衝突
  rm -f "$(which golangci-lint 2>/dev/null)" || true
  rm -f "$(go env GOPATH)/bin/golangci-lint" || true
  
  # run curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  #   | sh -s -- -b "$(go env GOPATH)/bin" "$LINT_VERSION"
  run go install github.com/golangci/golangci-lint/cmd/golangci-lint@"$LINT_VERSION"
fi

# ========================
# 运行 clean architecture 检查
# ========================
run go-cleanarch

log_info "lint modules:"
log_info "$(modules)"

# ========================
# 格式化代码
# ========================
run goimports -w -l .

# ========================
# lint 各模块（安全 cd）
# ========================
while read -r module; do
  log_info "linting module: $module"

  pushd "./internal/$module" >/dev/null

  run golangci-lint run --config "$ROOT_DIR/.golangci.yaml"

  popd >/dev/null
done < <(modules)