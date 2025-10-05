#!/bin/bash

# =============================================================================
# Go 程序交叉编译和部署脚本
# 
# 功能：
# 1. 交叉编译 Go 程序到多个平台
# 2. 自动创建目标目录
# 3. 通过 SCP 传输到指定服务器
# 4. 支持多种架构和操作系统
#
# 使用方法：
# ./cross-compile.sh [program_name]
# 
# 示例：
# ./cross-compile.sh docker_demo   # 编译 Docker 演示程序
# ./cross-compile.sh all           # 编译所有程序
# =============================================================================

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 配置信息
TARGET_SERVER="192.168.64.20"
TARGET_USER="root"
TARGET_PATH="/root/samuel/myDocker"
PROJECT_NAME="test-go"

# 支持的平台和架构 (使用普通数组以兼容 macOS)
PLATFORMS=(
    "linux-amd64:linux/amd64"
    "linux-arm64:linux/arm64"
    "windows-amd64:windows/amd64"
    "darwin-amd64:darwin/amd64"
    "darwin-arm64:darwin/amd64"
)

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "Go 程序交叉编译和部署脚本"
    echo ""
    echo "使用方法："
    echo "  $0 [program_name]"
    echo ""
    echo "参数说明："
    echo "  docker_demo 编译 Docker 演示程序"
    echo "  all          编译所有程序"
    echo ""
    echo "示例："
    echo "  $0 docker_demo"
    echo "  $0 all"
    echo ""
    echo "目标服务器：$TARGET_SERVER"
    echo "目标路径：$TARGET_PATH"
}

# 检查依赖
check_dependencies() {
    log_step "检查依赖..."
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装或不在 PATH 中"
        exit 1
    fi
    
    # 检查 SCP
    if ! command -v scp &> /dev/null; then
        log_error "SCP 未安装或不在 PATH 中"
        exit 1
    fi
    
    # 检查 SSH
    if ! command -v ssh &> /dev/null; then
        log_error "SSH 未安装或不在 PATH 中"
        exit 1
    fi
    
    log_success "所有依赖检查通过"
}

# 测试服务器连接
test_server_connection() {
    log_step "测试服务器连接..."
    
    if ssh -o ConnectTimeout=10 -o BatchMode=yes "$TARGET_USER@$TARGET_SERVER" "echo 'Connection successful'" &>/dev/null; then
        log_success "服务器连接成功"
    else
        log_warning "无法连接到服务器 $TARGET_SERVER"
        log_warning "请确保："
        log_warning "1. 服务器 IP 地址正确"
        log_warning "2. SSH 服务已启动"
        log_warning "3. 已配置 SSH 密钥认证"
        log_warning "4. 网络连接正常"
        
        read -p "是否继续编译？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# 创建目标目录
create_target_directory() {
    log_step "创建目标目录..."
    
    ssh "$TARGET_USER@$TARGET_SERVER" "mkdir -p $TARGET_PATH" 2>/dev/null || {
        log_warning "无法创建目标目录，将跳过文件传输"
        SKIP_UPLOAD=true
    }
}

# 清理旧的编译文件
clean_build() {
    log_step "清理旧的编译文件..."
    
    if [ -d "build" ]; then
        rm -rf build
        log_success "清理完成"
    else
        log_info "没有需要清理的文件"
    fi
}

# 编译单个程序
compile_program() {
    local program_name=$1
    local source_file=$2
    
    log_step "编译程序: $program_name"
    
    # 创建构建目录
    mkdir -p "build/$program_name"
    
    # 编译到各个平台
    for platform_info in "${PLATFORMS[@]}"; do
        platform=${platform_info%:*}
        goos=${platform%-*}
        goarch=${platform#*-}
        
        log_info "编译 $program_name for $goos/$goarch..."
        
        # 设置输出文件名
        if [ "$goos" = "windows" ]; then
            output_name="${program_name}.exe"
        else
            output_name="$program_name"
        fi
        
        # 执行编译
        cd "$program_name" 2>/dev/null || cd .
        
        env GOOS=$goos GOARCH=$goarch go build \
            -ldflags="-s -w" \
            -o "../build/$program_name/${platform}/${output_name}" \
            "$source_file" 2>/dev/null || {
            log_error "编译 $program_name for $platform 失败"
            cd ..
            continue
        }
        
        cd ..
        
        # 创建压缩包
        cd "build/$program_name"
        if [ "$goos" = "windows" ]; then
            zip -r "${platform}.zip" "$platform" &>/dev/null
        else
            tar -czf "${platform}.tar.gz" "$platform" &>/dev/null
        fi
        cd ../..
        
        log_success "编译完成: $program_name for $platform"
    done
}

# 上传文件到服务器
upload_to_server() {
    local program_name=$1
    
    if [ "$SKIP_UPLOAD" = true ]; then
        log_warning "跳过文件上传"
        return
    fi
    
    log_step "上传 $program_name 到服务器..."
    
    # 上传所有平台的编译结果
    for platform_info in "${PLATFORMS[@]}"; do
        platform=${platform_info%:*}
        cd "build/$program_name"
        
        if [ "$platform" = "windows-amd64" ]; then
            file_to_upload="${platform}.zip"
        else
            file_to_upload="${platform}.tar.gz"
        fi
        
        if [ -f "$file_to_upload" ]; then
            log_info "上传 $file_to_upload..."
            scp "$file_to_upload" "$TARGET_USER@$TARGET_SERVER:$TARGET_PATH/" || {
                log_error "上传 $file_to_upload 失败"
                cd ../..
                continue
            }
            log_success "上传成功: $file_to_upload"
        fi
        
        cd ../..
    done
    
    # 创建版本信息文件
    local version_info="build/$program_name/version.txt"
    cat > "$version_info" << EOF
程序名称: $program_name
编译时间: $(date)
Git 版本: $(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
Go 版本: $(go version)
编译平台: $(uname -s -m)
EOF
    
    scp "$version_info" "$TARGET_USER@$TARGET_PATH/" 2>/dev/null || {
        log_warning "上传版本信息失败"
    }
    
    log_success "所有文件上传完成"
}

# 验证上传结果
verify_upload() {
    if [ "$SKIP_UPLOAD" = true ]; then
        return
    fi
    
    log_step "验证上传结果..."
    
    local file_count=$(ssh "$TARGET_USER@$TARGET_SERVER" "ls -1 $TARGET_PATH/*.tar.gz $TARGET_PATH/*.zip 2>/dev/null | wc -l" 2>/dev/null || echo "0")
    
    if [ "$file_count" -gt 0 ]; then
        log_success "验证成功，已上传 $file_count 个文件"
        ssh "$TARGET_USER@$TARGET_SERVER" "ls -la $TARGET_PATH/" 2>/dev/null || true
    else
        log_warning "验证失败，未找到上传的文件"
    fi
}

# 主函数
main() {
    local program_name=${1:-"help"}
    
    echo -e "${CYAN}==============================================================================${NC}"
    echo -e "${CYAN}Go 程序交叉编译和部署脚本${NC}"
    echo -e "${CYAN}==============================================================================${NC}"
    echo ""
    
    # 显示帮助
    if [ "$program_name" = "help" ] || [ "$program_name" = "-h" ] || [ "$program_name" = "--help" ]; then
        show_help
        exit 0
    fi
    
    # 检查依赖
    check_dependencies
    
    # 测试服务器连接
    test_server_connection
    
    # 创建目标目录
    create_target_directory
    
    # 清理旧文件
    clean_build
    
    # 根据参数编译对应的程序
    case "$program_name" in
        "docker_demo")
            compile_program "docker_demo" "docker_demo.go"
            upload_to_server "docker_demo"
            ;;
        "all")
            log_step "编译所有程序..."
            
            # 编译 Docker 演示程序
            if [ -f "docker_demo/docker_demo.go" ]; then
                compile_program "docker_demo" "docker_demo.go"
                upload_to_server "docker_demo"
            fi
            ;;
        *)
            log_error "未知的程序名称: $program_name"
            log_info "支持的程序: docker_demo, all"
            exit 1
            ;;
    esac
    
    # 验证上传结果
    verify_upload
    
    echo ""
    echo -e "${GREEN}==============================================================================${NC}"
    echo -e "${GREEN}编译和部署完成！${NC}"
    echo -e "${GREEN}==============================================================================${NC}"
    echo ""
    echo -e "${CYAN}编译结果：${NC}"
    echo -e "构建目录: $(pwd)/build/"
    echo -e "目标服务器: $TARGET_SERVER"
    echo -e "目标路径: $TARGET_PATH"
    echo ""
    echo -e "${CYAN}快速连接服务器：${NC}"
    echo -e "ssh $TARGET_USER@$TARGET_SERVER"
    echo ""
    echo -e "${CYAN}查看上传的文件：${NC}"
    echo -e "ssh $TARGET_USER@$TARGET_SERVER 'ls -la $TARGET_PATH/'"
    echo ""
}

# 执行主函数
main "$@"
