# API 文檔

## 介紹

本文檔描述了交易策略管理 CLI 應用程序的所有可用命令和使用方法。

## 快速開始

### 編譯應用程序

```bash
cd /home/james/playground/go-transation
go build -o strategy-cli ./cmd/cli/main.go
```

### 基本用法

```bash
./strategy-cli [command] [subcommand] [flags]
```

### 查看幫助

```bash
./strategy-cli --help
./strategy-cli strategy --help
./strategy-cli strategy create --help
```

## 命令參考

### 1. 建立策略 (Create)

建立一個新的交易策略，定義買入和賣出的價格限制。

#### 命令

```bash
./strategy-cli strategy create [flags]
```

#### 標誌

| 短選項 | 長選項 | 類型 | 必須 | 說明 |
|--------|--------|------|------|------|
| `-s` | `--symbol` | string | ✓ | 交易對符號（例如 BTC/USD, ETH/USD） |
| `-b` | `--buy-lower` | float | ✓ | 買入價格下限（價格低於此值時觸發買信號） |
| `-u` | `--sell-upper` | float | ✓ | 賣出價格上限（價格高於此值時觸發賣信號） |

#### 約束

- `buy-lower` 必須 > 0
- `sell-upper` 必須 > `buy-lower`
- `symbol` 不能為空

#### 範例

```bash
# 建立 BTC/USD 策略，買入下限 50000，賣出上限 60000
./strategy-cli strategy create -s "BTC/USD" -b 50000 -u 60000

# 使用長選項
./strategy-cli strategy create --symbol "ETH/USD" --buy-lower 3000 --sell-upper 4000

# 輸出示例
# [INFO] 2025/11/05 Creating strategy symbol=BTC/USD
# Created strategy: ID=900dfecd-fc6e-47d7-8757-acfe833be778, Symbol=BTC/USD,
# BuyLower=50000.00, SellUpper=60000.00, Active=true
```

#### 響應

成功建立時返回：
- Strategy ID (唯一標識符)
- Symbol (交易對)
- BuyLower (買入下限)
- SellUpper (賣出上限)
- Active (策略狀態，默認為 true)

---

### 2. 列出所有策略 (List)

顯示所有已建立的交易策略。

#### 命令

```bash
./strategy-cli strategy list
```

#### 選項

無特殊選項。

#### 範例

```bash
./strategy-cli strategy list

# 輸出示例
# [INFO] Listing all strategies
# Strategies:
# ------------------------------------
# ID: 900dfecd-fc6e-47d7-8757-acfe833be778, Symbol: BTC/USD,
#     BuyLower: 50000.00, SellUpper: 60000.00, Status: Active
# ID: 33240ea5-8365-477d-9f4f-501725cd1e95, Symbol: ETH/USD,
#     BuyLower: 3000.00, SellUpper: 4000.00, Status: Active
# ------------------------------------
```

#### 響應

返回策略列表，包含：
- ID
- Symbol
- BuyLower
- SellUpper
- Status (Active/Inactive)

---

### 3. 查看策略詳情 (Get)

查看特定策略的詳細信息。

#### 命令

```bash
./strategy-cli strategy get <strategy-id>
```

#### 參數

| 參數 | 類型 | 說明 |
|------|------|------|
| `strategy-id` | string | 策略的唯一標識符 |

#### 範例

```bash
./strategy-cli strategy get 900dfecd-fc6e-47d7-8757-acfe833be778

# 輸出示例
# [INFO] Fetching strategy id=900dfecd-fc6e-47d7-8757-acfe833be778
# Strategy Details:
#   ID: 900dfecd-fc6e-47d7-8757-acfe833be778
#   Symbol: BTC/USD
#   Buy Lower: 50000.00
#   Sell Upper: 60000.00
#   Status: Active
```

#### 錯誤處理

如果策略不存在，將返回錯誤：
```
[ERROR] Strategy not found id=invalid-id
```

---

### 4. 更新策略 (Update)

修改現有策略的買入下限或賣出上限。

#### 命令

```bash
./strategy-cli strategy update <strategy-id> [flags]
```

#### 參數

| 參數 | 類型 | 說明 |
|------|------|------|
| `strategy-id` | string | 策略的唯一標識符 |

#### 標誌

| 短選項 | 長選項 | 類型 | 必須 | 說明 |
|--------|--------|------|------|------|
| `-b` | `--buy-lower` | float | ✗ | 新的買入價格下限 |
| `-u` | `--sell-upper` | float | ✗ | 新的賣出價格上限 |

#### 約束

- 至少指定一個標誌（`--buy-lower` 或 `--sell-upper`）
- 新的 `sell-upper` 必須 > 新的 `buy-lower`

#### 範例

```bash
# 更新買入下限
./strategy-cli strategy update 900dfecd-fc6e-47d7-8757-acfe833be778 -b 51000

# 更新賣出上限
./strategy-cli strategy update 900dfecd-fc6e-47d7-8757-acfe833be778 -u 61000

# 同時更新兩個值
./strategy-cli strategy update 900dfecd-fc6e-47d7-8757-acfe833be778 -b 51000 -u 61000

# 輸出示例
# [INFO] Updating strategy id=900dfecd-fc6e-47d7-8757-acfe833be778
# Updated strategy: ID=900dfecd-fc6e-47d7-8757-acfe833be778,
#                   BuyLower=51000.00, SellUpper=61000.00
```

#### 響應

返回更新後的策略信息。

---

### 5. 刪除策略 (Delete)

永久刪除指定的策略。

#### 命令

```bash
./strategy-cli strategy delete <strategy-id>
```

#### 參數

| 參數 | 類型 | 說明 |
|------|------|------|
| `strategy-id` | string | 要刪除的策略的唯一標識符 |

#### 範例

```bash
./strategy-cli strategy delete 900dfecd-fc6e-47d7-8757-acfe833be778

# 輸出示例
# [INFO] Deleting strategy id=900dfecd-fc6e-47d7-8757-acfe833be778
# Strategy 900dfecd-fc6e-47d7-8757-acfe833be778 deleted
```

#### 警告

刪除操作是永久的，無法撤銷。請謹慎使用。

---

### 6. 切換策略狀態 (Toggle)

啟用或停用指定的策略。

#### 命令

```bash
./strategy-cli strategy toggle <strategy-id>
```

#### 參數

| 參數 | 類型 | 說明 |
|------|------|------|
| `strategy-id` | string | 策略的唯一標識符 |

#### 說明

- 如果策略當前為**活躍**（Active），則將其禁用（Inactive）
- 如果策略當前為**禁用**（Inactive），則將其啟用（Active）

#### 範例

```bash
# 切換策略狀態（Active -> Inactive）
./strategy-cli strategy toggle 900dfecd-fc6e-47d7-8757-acfe833be778

# 輸出示例
# [INFO] Toggling strategy status id=900dfecd-fc6e-47d7-8757-acfe833be778
# Strategy 900dfecd-fc6e-47d7-8757-acfe833be778 is now Inactive

# 再次切換（Inactive -> Active）
./strategy-cli strategy toggle 900dfecd-fc6e-47d7-8757-acfe833be778

# 輸出示例
# Strategy 900dfecd-fc6e-47d7-8757-acfe833be778 is now Active
```

#### 響應

返回策略當前的狀態。

---

## 完整使用示例

### 場景：建立和管理 BTC 交易策略

```bash
# 1. 建立策略
./strategy-cli strategy create -s "BTC/USD" -b 45000 -u 55000
# 假設 ID 為: abc123def456

# 2. 查看策略詳情
./strategy-cli strategy get abc123def456

# 3. 查看所有策略
./strategy-cli strategy list

# 4. 更新策略（買入下限調整為 48000）
./strategy-cli strategy update abc123def456 -b 48000

# 5. 暫時禁用策略
./strategy-cli strategy toggle abc123def456

# 6. 重新啟用策略
./strategy-cli strategy toggle abc123def456

# 7. 刪除策略
./strategy-cli strategy delete abc123def456
```

---

## 數據結構

### Strategy

```go
type Strategy struct {
    ID        string    // 唯一標識符 (UUID)
    Symbol    string    // 交易對符號 (例如: "BTC/USD")
    BuyLower  float64   // 買入價格下限
    SellUpper float64   // 賣出價格上限
    IsActive  bool      // 策略是否活躍
}
```

### CreateStrategyRequest

```go
type CreateStrategyRequest struct {
    Symbol    string  // 必須
    BuyLower  float64 // 必須，> 0
    SellUpper float64 // 必須，> BuyLower
}
```

### UpdateStrategyRequest

```go
type UpdateStrategyRequest struct {
    ID        string  // 策略 ID
    Symbol    string  // 交易對符號
    BuyLower  float64 // 可選
    SellUpper float64 // 可選
}
```

### StrategyResponse

```go
type StrategyResponse struct {
    ID        string  // 唯一標識符
    Symbol    string  // 交易對符號
    BuyLower  float64 // 買入價格下限
    SellUpper float64 // 賣出價格上限
    IsActive  bool    // 策略狀態
}
```

---

## 錯誤處理

### 常見錯誤

| 錯誤 | 原因 | 解決方案 |
|------|------|--------|
| `buy lower bound must be positive` | 買入下限 ≤ 0 | 設置 > 0 的值 |
| `sell upper bound must be greater than buy lower bound` | 賣出上限 ≤ 買入下限 | 確保賣出上限 > 買入下限 |
| `strategy not found` | 策略不存在 | 確認策略 ID 正確 |
| `at least one of --buy-lower or --sell-upper is required` | 更新時未指定任何標誌 | 指定至少一個要更新的字段 |
| `symbol is required` | 建立時未指定符號 | 使用 `-s` 或 `--symbol` 指定符號 |

---

## 日誌說明

應用程序會輸出 INFO 和 ERROR 級別的日誌：

```
[INFO] 時間戳 操作信息 key=value
[ERROR] 時間戳 錯誤信息 key=value
[WARN] 時間戳 警告信息 key=value
```

### 日誌示例

```bash
[INFO] 2025/11/05 01:25:44 Creating strategy symbol=BTC/USD
[INFO] 2025/11/05 01:25:44 Strategy created successfully id=abc123def456
[INFO] 2025/11/05 01:25:45 Listing all strategies
[ERROR] 2025/11/05 01:25:46 Strategy not found id=invalid-id
```

---

## 最佳實踐

### 1. 策略命名

使用清晰的符號命名：
- ✓ `BTC/USD`, `ETH/USD`, `SOL/USDT`
- ✗ `BTC`, `bitcoin`, `test`

### 2. 價格配置

根據市場情況合理設置價格：
- 買入下限：應該低於當前市場價格
- 賣出上限：應該高於當前市場價格
- 安全邊界：考慮市場波動

### 3. 狀態管理

- 使用 toggle 命令暫停不需要的策略
- 定期檢查 list 命令確認所有策略狀態
- 在刪除前使用 get 命令確認策略信息

### 4. 備份

- SQLite 數據庫文件（strategies.db）定期備份
- 重要的策略配置記錄在外部文檔中

---

## 環境變量

目前無特殊環境變量需要配置。

## 配置文件

目前應用程序使用 SQLite 數據庫（strategies.db）在當前目錄中存儲數據。

## 故障排除

### 問題：命令未找到

```bash
bash: ./strategy-cli: 找不到命令
```

**解決方案**：
1. 確認已編譯應用程序：`go build -o strategy-cli ./cmd/cli/main.go`
2. 使用正確的路徑：`./strategy-cli` 或 `./strategy-cli strategy list`

### 問題：數據庫連接失敗

```bash
Failed to connect to database
```

**解決方案**：
1. 確認當前目錄有寫入權限
2. 刪除舊的 strategies.db 文件：`rm strategies.db`
3. 重新運行應用程序會自動建立新數據庫

### 問題：策略驗證失敗

```bash
Error: buy lower bound must be positive
```

**解決方案**：
1. 確認買入下限 > 0
2. 確認賣出上限 > 買入下限
3. 檢查輸入的數值是否為有效的浮點數

---

## 更新日誌

### 版本 1.0.0 (2025-11-05)

- ✓ 實現基本的 CRUD 操作
- ✓ 支持策略狀態管理
- ✓ 完整的 CLI 命令
- ✓ SQLite 數據庫支持
- ✓ 完善的錯誤處理

