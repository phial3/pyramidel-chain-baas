VERSION := 1.0.0
BUILD := $(shell git rev-parse --short HEAD)

.DEFAULT_GOAL := all

# 指定输出目录
OUTPUT_DIR = ./bin

# 指定 Go 编译器
GOCMD=go
# 指定 Go 编译选项
GOBUILD=$(GOCMD) build
GO_BUILD_FLAGS := -ldflags="-X 'main.version=$(VERSION)' -X 'main.build=$(BUILD)'" -a -tags netgo -installsuffix netgo
CGO_ENABLED := 0

RELEASE_EXES := serve psutil

PKGNAME = github.com/hxx258456/pyramidel-chain-baas

# 指定 Go 测试选项
GOTEST=$(GOCMD) test

# 指定所有的目标
TARGETS = serve psutil

# 指定支持的操作系统和体系结构
OS_ARCHS = linux/amd64

# 默认目标，编译所有平台的目标
all: test $(foreach target, $(TARGETS), $(foreach osarch, $(OS_ARCHS), $(OUTPUT_DIR)/$(target)-$(osarch)))

# 编译指定目标和平台的二进制文件
$(OUTPUT_DIR)/%-darwin/amd64:
    GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $@ $(PKGNAME)/cmd/$(subst -darwin/amd64,,$*)

$(OUTPUT_DIR)/%-linux/amd64:
    GOOS=linux GOARCH=amd64 $(GOBUILD) -o $@ $(PKGNAME)/cmd/$(subst -linux/amd64,,$*)

$(OUTPUT_DIR)/%-windows/amd64:
    GOOS=windows GOARCH=amd64 $(GOBUILD) -o $@.exe $(PKGNAME)/cmd/$(subst -windows/amd64,,$*)

# 执行所有测试用例
test:
	$(GOTEST) $(PKGNAME)

.PHONY: test

# 清除所有目标文件
clean:
	rm -f $(foreach target, $(TARGETS), $(foreach osarch, $(OS_ARCHS), $(OUTPUT_DIR)/$(target)-$(osarch)))

.PHONY: clean

