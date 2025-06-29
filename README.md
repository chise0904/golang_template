# Golang Template

這是一個 Go 語言的專案模板，支援 VSCode F5 debug 功能。

## 專案結構

```
golang_template/
├── cmd/
│   ├── main.go          # 主程式入口
│   └── server/
│       └── server.go    # 伺服器命令
├── config/
│   ├── config.go        # 配置結構定義
│   └── config.yaml      # 配置檔案
├── migrations/          # 資料庫 migration 檔案
│   └── *.sql           # Goose migration 檔案
├── pkg/                 # 套件庫
└── .vscode/
    └── launch.json      # VSCode debug 配置
```

## 資料庫 Migration

本專案使用 [Goose](https://github.com/pressly/goose) 來管理資料庫 migration。

### 安裝 Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Migration 檔案格式

Migration 檔案命名格式：`YYYYMMDDHHMMSS_description.sql`

檔案內容包含：
- `-- +goose Up`: 執行 migration 時的 SQL
- `-- +goose Down`: 回滾 migration 時的 SQL

### 基本指令

#### 檢查 migration 狀態
```bash
goose status
```

#### 執行所有 pending 的 migration
```bash
goose up
```

#### 執行到特定版本
```bash
goose up-to 20241201000000
```

#### 回滾最後一個 migration
```bash
goose down
```

#### 回滾到特定版本
```bash
goose down-to 20241201000000
```

#### 重置所有 migration（回滾到最初狀態）
```bash
goose reset
```

#### 重新執行所有 migration
```bash
goose redo
```

### 資料庫連接

#### PostgreSQL
```bash
goose postgres "user=myusername password=mypassword dbname=mydbname sslmode=disable" up
```

#### MySQL
```bash
goose mysql "myusername:mypassword@tcp(localhost:3306)/mydbname?parseTime=true" up
```

### 實際使用範例

#### PostgreSQL 範例
```bash
# 執行 migration
goose postgres "user=admin password=admin host=172.17.0.5 dbname=identity sslmode=disable" up

# 檢查狀態
goose postgres "user=admin password=admin dbname=identity sslmode=disable" status

# 回滾
goose postgres "user=admin password=admin dbname=identity sslmode=disable" down
```

#### MySQL 範例
```bash
# 執行 migration
goose mysql "root:password@tcp(localhost:3306)/golang_template?parseTime=true" up

# 檢查狀態
goose mysql "root:password@tcp(localhost:3306)/golang_template?parseTime=true" status
```

#### 使用環境變數（推薦）
```bash
# 設定環境變數
export GOOSE_DBSTRING="user=admin password=admin dbname=identity sslmode=disable"
export GOOSE_DIALECT="postgres"

# 執行指令
goose up
goose status
goose down
```

### 資料庫表格結構

本專案包含以下主要表格：

1. **account_info**: 帳戶基本資訊
   - 包含帳戶 ID、類型、狀態、權限等
   - 支援軟刪除

2. **user_profile**: 使用者個人資料
   - 包含使用者名稱、個人資訊、地址等
   - 與 account_info 一對一關聯

3. **verification_code**: 驗證碼管理
   - 用於 email、手機驗證等
   - 支援多種驗證類型

## VSCode F5 Debug 設定

專案已經配置好 VSCode 的 debug 功能，包含以下配置：

### 1. Debug Server
- **名稱**: Debug Server
- **功能**: 啟動伺服器並進行 debug
- **程式**: `cmd/main.go`
- **參數**: `server`
- **環境變數**: `CONFIG_PATH=./config`

### 2. Debug Server (with custom config)
- **名稱**: Debug Server (with custom config)
- **功能**: 使用自定義環境變數啟動伺服器
- **額外功能**: 支援 `.env` 檔案

### 3. Debug Tests
- **名稱**: Debug Tests
- **功能**: 執行測試並進行 debug

## 使用方法

1. **開啟專案**: 在 VSCode 中開啟專案資料夾
2. **設定斷點**: 在需要 debug 的程式碼行號左側點擊設定斷點
3. **啟動 Debug**: 按 `F5` 或點擊 Debug 面板中的 "Start Debugging"
4. **選擇配置**: 選擇 "Debug Server" 配置
5. **開始除錯**: 程式會在斷點處停止，可以檢查變數、執行步驟等

## 配置檔案

### config/config.yaml
包含所有服務的配置，包括：
- 日誌配置
- 資料庫配置
- Web 服務配置
- Redis 配置
- gRPC 配置

## 依賴套件

主要使用的套件：
- `github.com/spf13/cobra` - 命令列工具
- `go.uber.org/fx` - 依賴注入框架
- `github.com/rs/zerolog` - 日誌套件
- `github.com/spf13/viper` - 配置管理
- `gorm.io/gorm` - ORM 框架
- `github.com/pressly/goose` - 資料庫 migration 工具

## 開發建議

1. **設定斷點**: 在關鍵邏輯處設定斷點
2. **使用 Watch**: 在 debug 時監控重要變數
3. **檢查 Call Stack**: 了解程式執行流程
4. **使用 Console**: 在 debug console 中執行表達式
5. **Migration 管理**: 使用 goose 管理資料庫結構變更
6. **版本控制**: 將 migration 檔案納入版本控制

## 故障排除

如果遇到 debug 問題：

1. **檢查 Go 擴充功能**: 確保已安裝 Go 擴充功能
2. **檢查 GOPATH**: 確保 Go 環境配置正確
3. **檢查依賴**: 執行 `go mod tidy` 確保依賴完整
4. **檢查配置**: 確保 `config/config.yaml` 檔案存在且格式正確

如果遇到 migration 問題：

1. **檢查資料庫連接**: 確保資料庫服務正在運行
2. **檢查權限**: 確保資料庫使用者有足夠權限
3. **檢查 migration 狀態**: 使用 `goose status` 檢查當前狀態
4. **檢查 SQL 語法**: 確保 migration 檔案中的 SQL 語法正確 





```
sudo apt update
sudo apt install -y protobuf-compiler

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```
protoc --proto_path=proto --go_out=proto/pkg --go-grpc_out=proto/pkg --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/common/*.proto

protoc --proto_path=proto --go_out=proto/pkg --go-grpc_out=proto/pkg --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/identity/message/*.proto

protoc --proto_path=proto --go_out=proto/pkg --go-grpc_out=proto/pkg --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/identity/*.proto
```