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
├── pkg/                 # 套件庫
└── .vscode/
    └── launch.json      # VSCode debug 配置
```

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

## 開發建議

1. **設定斷點**: 在關鍵邏輯處設定斷點
2. **使用 Watch**: 在 debug 時監控重要變數
3. **檢查 Call Stack**: 了解程式執行流程
4. **使用 Console**: 在 debug console 中執行表達式

## 故障排除

如果遇到 debug 問題：

1. **檢查 Go 擴充功能**: 確保已安裝 Go 擴充功能
2. **檢查 GOPATH**: 確保 Go 環境配置正確
3. **檢查依賴**: 執行 `go mod tidy` 確保依賴完整
4. **檢查配置**: 確保 `config/config.yaml` 檔案存在且格式正確 