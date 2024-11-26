.PHONY: all build run gotool clean help

BINARY_NAME="app"
RELEASE_DIR="release"

# 构建所有平台的二进制文件
build-all: gotool
	mkdir -p $(RELEASE_DIR)
	make build-linux
	make build-arm
	make build-windows

# 构建 Linux 平台的二进制文件
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BINARY_NAME)_linux_amd64

# 构建 ARM64 平台的二进制文件
build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(RELEASE_DIR)/$(BINARY_NAME)_linux_arm64

# 构建 Windows 平台的二进制文件
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BINARY_NAME)_windows_amd64.exe

# 直接运行 Go 代码
run:
	@go run ./

# 格式化代码并运行 Go 工具
gotool:
	go fmt ./
	go vet ./

# 清理构建产物
clean:
	@if [ -d $(RELEASE_DIR) ]; then rm -rf $(RELEASE_DIR)/$(BINARY_NAME)* ; fi
	@find . -name "*.swp" -exec rm -f {} \;

# 帮助信息
help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build-all - 编译 Go 代码, 生成所有平台的二进制文件"
	@echo "make build-linux - 编译 Go 代码, 生成 Linux 平台的二进制文件"
	@echo "make build-arm - 编译 Go 代码, 生成 ARM64 平台的二进制文件"
	@echo "make build-windows - 编译 Go 代码, 生成 Windows 平台的二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
