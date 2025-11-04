# 後端專案憲法 (Backend Constitution): go-transaction

> **最後更新日期:** 2025-11-04
> **核心技術棧:** Go 1.21+, Gin, SQLite, GORM, Cobra, Goroutines
> **專案性質:** 加密貨幣交易決策輔助系統

---

## 目錄

1. [核心原則](#1-核心原則-guiding-principles)
2. [專案結構規範](#2-專案結構規範-project-structure)
3. [程式碼風格與格式化](#3-程式碼風格與格式化-code-style--formatting)
4. [命名慣例](#4-命名慣例-naming-conventions)
5. [Git 工作流程](#5-git-工作流程-git-workflow)
6. [API 設計規範](#6-api-設計規範-api-design)
7. [資料庫設計規範](#7-資料庫設計規範-database-design)
8. [業務邏輯層規範](#8-業務邏輯層規範-business-logic-layer)
9. [並發與背景任務](#9-並發與背景任務-concurrency--background-tasks)
10. [錯誤處理與日誌](#10-錯誤處理與日誌-error-handling--logging)
11. [測試策略](#11-測試策略-testing-strategy)
12. [安全與性能](#12-安全與性能-security--performance)
13. [依賴管理](#13-依賴管理-dependency-management)
14. [配置管理](#14-配置管理-configuration-management)
15. [部署與運維](#15-部署與運維-deployment--operations)

---

## 1. 核心原則 (Guiding Principles)

### 1.1 設計哲學

- **簡潔優於複雜 (Simplicity over Complexity)**: 遵循 Go 語言的設計哲學，優先選擇簡單直接的解決方案。
- **顯式優於隱式 (Explicit over Implicit)**: 讓依賴關係和資料流動清晰可見，避免「魔法」行為。
- **可測試性優先 (Testability First)**: 所有業務邏輯必須易於測試，透過依賴注入實現解耦。
- **優雅失敗 (Fail Gracefully)**: 所有外部依賴（API、資料庫）都可能失敗，必須有適當的錯誤處理。
- **效能與可讀性的平衡**: 在滿足性能需求（30 秒延遲）的前提下，優先考慮程式碼的可讀性和維護性。

### 1.2 關鍵性能需求

| 需求 | 目標值 | 權衡說明 |
|------|--------|----------|
| 價格監控延遲 | < 30 秒 | 使用 goroutine 實現並發價格獲取，但需考慮 API rate limiting |
| 系統穩定性 | 99% uptime | 必須實現 panic recovery、重試機制、優雅關閉 |
| 測試覆蓋率 | > 80% | 核心業務邏輯（策略管理、價格比對）必須達到 90%+ |

---

## 2. 專案結構規範 (Project Structure)

### 2.1 目錄結構

遵循 **Standard Go Project Layout** 與 **Clean Architecture** 的混合模式：

```
go-transaction/
├── cmd/                          # 應用程式入口點
│   ├── cli/                      # CLI 應用程式
│   │   └── main.go
│   └── server/                   # HTTP 伺服器（未來）
│       └── main.go
├── internal/                     # 私有應用程式程式碼（不可被外部 import）
│   ├── domain/                   # 領域模型與業務規則（核心層）
│   │   ├── strategy.go           # 策略實體
│   │   ├── trade.go              # 交易實體
│   │   ├── price.go              # 價格實體
│   │   └── errors.go             # 領域錯誤定義
│   ├── usecase/                  # 業務邏輯層（用例層）
│   │   ├── strategy/             # 策略管理用例
│   │   │   ├── service.go        # 策略服務介面與實作
│   │   │   └── service_test.go
│   │   ├── trade/                # 交易執行用例
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   └── monitor/              # 價格監控用例
│   │       ├── service.go
│   │       └── service_test.go
│   ├── adapter/                  # 適配器層（外部系統整合）
│   │   ├── repository/           # 資料儲存實作
│   │   │   ├── sqlite/
│   │   │   │   ├── strategy_repo.go
│   │   │   │   ├── trade_repo.go
│   │   │   │   └── migration.go
│   │   │   └── repository.go     # 儲存庫介面定義
│   │   ├── exchange/             # 交易所 API 客戶端
│   │   │   ├── exchange.go       # 交易所介面定義
│   │   │   ├── binance/
│   │   │   │   ├── client.go
│   │   │   │   └── client_test.go
│   │   │   └── mock/             # Mock 實作（測試用）
│   │   │       └── exchange_mock.go
│   │   └── notifier/             # 通知服務實作
│   │       ├── notifier.go       # 通知介面定義
│   │       ├── console/
│   │       │   └── console.go
│   │       └── telegram/         # 未來擴充
│   │           └── telegram.go
│   └── interface/                # 介面層（使用者互動）
│       ├── cli/                  # CLI 指令實作
│       │   ├── strategy_cmd.go
│       │   ├── trade_cmd.go
│       │   └── monitor_cmd.go
│       └── http/                 # HTTP Handler（未來）
│           └── handler.go
├── pkg/                          # 可被外部專案 import 的公開程式碼
│   ├── logger/                   # 日誌工具
│   │   └── logger.go
│   └── retry/                    # 重試機制工具
│       └── retry.go
├── config/                       # 配置檔案
│   ├── config.yaml
│   └── config.go                 # 配置載入邏輯
├── scripts/                      # 腳本工具
│   ├── migrate.sh                # 資料庫遷移腳本
│   └── build.sh                  # 建置腳本
├── test/                         # 額外的測試資料與整合測試
│   ├── integration/              # 整合測試
│   │   └── monitor_test.go
│   └── fixtures/                 # 測試用固定資料
│       └── test_data.sql
├── docs/                         # 專案文件
│   ├── architecture.md
│   └── api.md
├── go.mod
├── go.sum
├── Makefile                      # 常用指令集合
├── README.md
├── plans.md                      # 產品規劃文件
├── rule.md                       # 本文件
└── tasks.md                      # TDD 任務清單

```

### 2.2 分層職責說明

#### Domain Layer (領域層)
- **職責**: 定義核心業務實體、值物件、業務規則。
- **原則**:
  - 不依賴任何外部套件（除標準庫）。
  - 包含業務邏輯的驗證規則（如價格上下限的合法性檢查）。
  - 定義領域事件（如 `StrategyTriggered`）。
- **範例**:
  ```go
  // internal/domain/strategy.go
  package domain

  import "errors"

  type Strategy struct {
      ID           string
      Symbol       string  // BTC, ETH, USDT
      BuyLower     float64
      SellUpper    float64
      IsActive     bool
      CreatedAt    time.Time
  }

  func (s *Strategy) Validate() error {
      if s.BuyLower <= 0 {
          return errors.New("buy lower bound must be positive")
      }
      if s.SellUpper <= s.BuyLower {
          return errors.New("sell upper bound must be greater than buy lower bound")
      }
      return nil
  }

  func (s *Strategy) ShouldBuy(currentPrice float64) bool {
      return s.IsActive && currentPrice <= s.BuyLower
  }

  func (s *Strategy) ShouldSell(currentPrice float64) bool {
      return s.IsActive && currentPrice >= s.SellUpper
  }
  ```

#### UseCase Layer (用例層)
- **職責**: 協調領域物件和外部服務，實現具體的業務流程。
- **原則**:
  - 依賴介面而非具體實作（依賴反轉原則）。
  - 不包含技術細節（如 SQL、HTTP 請求）。
  - 透過建構子注入依賴。
- **範例**:
  ```go
  // internal/usecase/strategy/service.go
  package strategy

  import "context"

  type Repository interface {
      Create(ctx context.Context, strategy *domain.Strategy) error
      FindByID(ctx context.Context, id string) (*domain.Strategy, error)
      FindAll(ctx context.Context) ([]*domain.Strategy, error)
      Update(ctx context.Context, strategy *domain.Strategy) error
      Delete(ctx context.Context, id string) error
  }

  type Service struct {
      repo   Repository
      logger logger.Logger
  }

  func NewService(repo Repository, logger logger.Logger) *Service {
      return &Service{
          repo:   repo,
          logger: logger,
      }
  }

  func (s *Service) CreateStrategy(ctx context.Context, req CreateStrategyRequest) (*domain.Strategy, error) {
      strategy := &domain.Strategy{
          ID:        uuid.New().String(),
          Symbol:    req.Symbol,
          BuyLower:  req.BuyLower,
          SellUpper: req.SellUpper,
          IsActive:  true,
          CreatedAt: time.Now(),
      }

      if err := strategy.Validate(); err != nil {
          return nil, fmt.Errorf("invalid strategy: %w", err)
      }

      if err := s.repo.Create(ctx, strategy); err != nil {
          s.logger.Error("failed to create strategy", "error", err)
          return nil, err
      }

      return strategy, nil
  }
  ```

#### Adapter Layer (適配器層)
- **職責**: 實作具體的技術細節（資料庫操作、API 呼叫、檔案讀寫）。
- **原則**:
  - 實作 UseCase 層定義的介面。
  - 處理技術層面的錯誤（如連線逾時、序列化錯誤）。
- **範例**:
  ```go
  // internal/adapter/repository/sqlite/strategy_repo.go
  package sqlite

  import (
      "context"
      "gorm.io/gorm"
  )

  type StrategyRepository struct {
      db *gorm.DB
  }

  func NewStrategyRepository(db *gorm.DB) *StrategyRepository {
      return &StrategyRepository{db: db}
  }

  func (r *StrategyRepository) Create(ctx context.Context, strategy *domain.Strategy) error {
      return r.db.WithContext(ctx).Create(strategy).Error
  }

  func (r *StrategyRepository) FindByID(ctx context.Context, id string) (*domain.Strategy, error) {
      var strategy domain.Strategy
      err := r.db.WithContext(ctx).First(&strategy, "id = ?", id).Error
      if err != nil {
          return nil, err
      }
      return &strategy, nil
  }
  ```

#### Interface Layer (介面層)
- **職責**: 處理使用者輸入與輸出（CLI、HTTP、gRPC）。
- **原則**:
  - 負責資料格式轉換（JSON、命令列參數）。
  - 不包含業務邏輯。
  - 將請求委派給 UseCase 層。

---

## 3. 程式碼風格與格式化 (Code Style & Formatting)

### 3.1 強制性工具

**所有程式碼在提交前必須通過以下檢查：**

```bash
# 格式化
go fmt ./...

# Linting（強制使用 golangci-lint）
golangci-lint run

# 靜態分析
go vet ./...
```

**`.golangci.yml` 配置範例：**

```yaml
linters:
  enable:
    - gofmt
    - golint
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - misspell
    - gocyclo       # 複雜度檢查
    - dupl          # 重複程式碼檢查
    - gosec         # 安全檢查
    - goconst       # 常數提取檢查

linters-settings:
  gocyclo:
    min-complexity: 15  # 函式複雜度上限
  dupl:
    threshold: 100      # 重複程式碼行數閾值
```

### 3.2 Go 官方風格指南

**必讀文件：**
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

**關鍵原則：**
- 優先使用短變數名稱（在局部範圍內）。
- 避免過度註解（程式碼應該自我解釋）。
- 錯誤處理不可忽略（不可使用 `_ = err`）。

---

## 4. 命名慣例 (Naming Conventions)

### 4.1 通用規則

| 類型 | 慣例 | 範例 |
|------|------|------|
| 套件名稱 | 小寫單字，不使用底線或駝峰 | `strategy`, `monitor`, `notifier` |
| 檔案名稱 | 小寫 + 底線 | `strategy_service.go`, `binance_client.go` |
| 介面 | PascalCase, 通常以 `-er` 結尾 | `Repository`, `Notifier`, `PriceGetter` |
| 結構體 | PascalCase | `Strategy`, `TradeRecord`, `PriceMonitor` |
| 方法/函式 | PascalCase（公開）, camelCase（私有） | `CreateStrategy()`, `validateInput()` |
| 變數 | camelCase | `currentPrice`, `strategyRepo` |
| 常數 | PascalCase 或全大寫 | `MaxRetryAttempts`, `DEFAULT_TIMEOUT` |

### 4.2 特定場景命名

#### 介面命名

**優先使用單一方法介面：**
```go
// Good: 單一職責
type PriceGetter interface {
    GetPrice(ctx context.Context, symbol string) (float64, error)
}

type Notifier interface {
    Send(ctx context.Context, message string) error
}

// Bad: 過於龐大的介面
type ExchangeService interface {
    GetPrice() error
    GetOrderBook() error
    PlaceOrder() error
    // ... 10 個方法
}
```

#### 錯誤變數命名

```go
// Good: 使用 Err 前綴
var (
    ErrStrategyNotFound = errors.New("strategy not found")
    ErrInvalidPrice     = errors.New("invalid price")
)

// Bad
var (
    StrategyNotFoundError = errors.New("...")
)
```

#### 測試檔案命名

```go
// 被測試檔案
strategy_service.go

// 對應測試檔案
strategy_service_test.go

// 測試函式命名
func TestService_CreateStrategy(t *testing.T) {}
func TestStrategy_Validate_WithInvalidPrice(t *testing.T) {}
```

### 4.3 資料庫相關命名

```go
// 表名：小寫 + 底線 + 複數
strategies
trade_records
notification_logs

// 欄位名：小寫 + 底線
buy_lower_bound
sell_upper_bound
created_at

// GORM 結構體標籤
type Strategy struct {
    ID           string    `gorm:"primaryKey;column:id"`
    Symbol       string    `gorm:"column:symbol;not null"`
    BuyLower     float64   `gorm:"column:buy_lower_bound;not null"`
    SellUpper    float64   `gorm:"column:sell_upper_bound;not null"`
    IsActive     bool      `gorm:"column:is_active;default:true"`
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Strategy) TableName() string {
    return "strategies"
}
```

---

## 5. Git 工作流程 (Git Workflow)

### 5.1 分支策略

```
main                    # 穩定版本，隨時可部署
  └── develop           # 開發整合分支
        ├── feat/T-01-create-strategy       # 功能分支
        ├── feat/T-06-price-monitoring      # 功能分支
        ├── fix/issue-123-panic-recovery    # 修復分支
        └── refactor/extract-price-service  # 重構分支
```

**分支命名規則：**
```
feat/<ticket-id>-<short-description>        # 新功能
fix/<ticket-id>-<short-description>         # Bug 修復
refactor/<description>                      # 重構（不改變行為）
chore/<description>                         # 雜項（如更新依賴）
test/<description>                          # 測試相關
```

### 5.2 Commit 訊息規範

**遵循 Conventional Commits 規範：**

```
<type>(<scope>): <subject>

<body>

<footer>
```

**類型 (Type)：**
```
feat:     新功能
fix:      Bug 修復
refactor: 重構（不改變行為）
test:     新增或修改測試
docs:     文件更新
chore:    建置流程、依賴更新
perf:     效能優化
ci:       CI/CD 相關
```

**範例：**

```
# Good
feat(strategy): add create strategy use case

Implement the core logic for creating a new price range strategy.
This includes:
- Input validation
- Domain model creation
- Repository persistence

Closes #T-01

# Good
refactor(monitor): extract price comparison logic to domain model

Move the price comparison logic from service layer to Strategy domain
entity to improve testability and follow DDD principles.

# Bad
update code       # 不清楚改了什麼
fix bug           # 沒說明修復了哪個 bug
```

### 5.3 Pull Request 規範

**PR 標題格式：**
```
[<Type>] <Summary>

範例：
[Feat] Add strategy management use case (T-01, T-02, T-03)
[Fix] Fix panic in price monitor when API returns nil
[Refactor] Extract exchange client interface
```

**PR 描述必須包含：**

```markdown
## 變更摘要 (Summary)
簡述這個 PR 做了什麼

## 變更類型 (Change Type)
- [ ] 行為性變更 (Behavioral Change)
- [ ] 結構性變更 (Structural Change / Refactoring)

## 相關任務 (Related Tasks)
Closes #T-01, #T-02

## 測試 (Testing)
- [ ] 單元測試已通過 (80%+ coverage)
- [ ] 整合測試已通過
- [ ] 手動測試步驟：
  1. ...
  2. ...

## 檢查清單 (Checklist)
- [ ] 程式碼已通過 `golangci-lint run`
- [ ] 程式碼已格式化 (`go fmt ./...`)
- [ ] 所有測試通過 (`go test -v ./...`)
- [ ] 已更新相關文件（如有需要）
- [ ] 資料庫遷移檔案已建立（如有需要）
```

---

## 6. API 設計規範 (API Design)

### 6.1 RESTful API 原則

**資源導向設計：**

```
# Good: 使用名詞（複數）
GET    /api/v1/strategies          # 取得所有策略
POST   /api/v1/strategies          # 建立新策略
GET    /api/v1/strategies/{id}     # 取得特定策略
PUT    /api/v1/strategies/{id}     # 更新策略
DELETE /api/v1/strategies/{id}     # 刪除策略

GET    /api/v1/trades              # 取得交易記錄
POST   /api/v1/trades              # 執行新交易

# Bad: 使用動詞
POST   /api/v1/createStrategy
GET    /api/v1/getStrategyById?id=123
```

### 6.2 HTTP 狀態碼使用

```go
200 OK                  # 成功處理 GET, PUT
201 Created             # 成功建立資源 (POST)
204 No Content          # 成功刪除資源 (DELETE)
400 Bad Request         # 客戶端請求錯誤（如驗證失敗）
401 Unauthorized        # 未認證
403 Forbidden           # 已認證但無權限
404 Not Found           # 資源不存在
409 Conflict            # 資源衝突（如重複建立）
500 Internal Server Error  # 伺服器內部錯誤
503 Service Unavailable    # 外部服務不可用（如交易所 API 失敗）
```

### 6.3 標準化的請求/回應格式

**成功回應：**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "symbol": "BTC",
    "buy_lower": 60000,
    "sell_upper": 70000,
    "is_active": true,
    "created_at": "2025-11-04T10:30:00Z"
  }
}
```

**錯誤回應：**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_STRATEGY",
    "message": "Sell upper bound must be greater than buy lower bound",
    "details": {
      "field": "sell_upper",
      "provided": 50000,
      "required": "> 60000"
    }
  }
}
```

**Go 結構體定義：**
```go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}
```

### 6.4 錯誤碼設計

```go
const (
    // 通用錯誤
    ErrCodeInternalError     = "INTERNAL_ERROR"
    ErrCodeValidationFailed  = "VALIDATION_FAILED"
    ErrCodeNotFound          = "NOT_FOUND"

    // 策略相關
    ErrCodeInvalidStrategy   = "INVALID_STRATEGY"
    ErrCodeStrategyNotFound  = "STRATEGY_NOT_FOUND"

    // 交易相關
    ErrCodeInsufficientFunds = "INSUFFICIENT_FUNDS"
    ErrCodeTradeExecutionFailed = "TRADE_EXECUTION_FAILED"

    // 外部服務
    ErrCodeExchangeAPIFailed = "EXCHANGE_API_FAILED"
)
```

### 6.5 API 版本管理

```go
// Good: 使用 URL 路徑版本控制
/api/v1/strategies
/api/v2/strategies  // 未來版本

// 在 Gin 中實作
func setupRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        v1.GET("/strategies", handler.GetStrategies)
        v1.POST("/strategies", handler.CreateStrategy)
    }
}
```

---

## 7. 資料庫設計規範 (Database Design)

### 7.1 表結構設計原則

**基本規範：**
- 表名使用小寫 + 底線，且為複數形式。
- 主鍵統一命名為 `id`，使用 UUID 或 ULID（避免自增 ID 外洩資訊）。
- 所有表必須包含 `created_at` 和 `updated_at` 時間戳。
- 使用 `deleted_at` 實作軟刪除（soft delete）。

**核心表結構：**

```sql
-- 策略表
CREATE TABLE strategies (
    id            TEXT PRIMARY KEY,
    symbol        TEXT NOT NULL,              -- BTC, ETH, USDT
    buy_lower     REAL NOT NULL,              -- 買入價格下限
    sell_upper    REAL NOT NULL,              -- 賣出價格上限
    is_active     BOOLEAN NOT NULL DEFAULT 1, -- 啟用狀態
    created_at    DATETIME NOT NULL,
    updated_at    DATETIME NOT NULL,
    deleted_at    DATETIME,                   -- 軟刪除
    CONSTRAINT ck_buy_lower_positive CHECK (buy_lower > 0),
    CONSTRAINT ck_sell_upper_greater CHECK (sell_upper > buy_lower)
);

-- 交易記錄表
CREATE TABLE trade_records (
    id            TEXT PRIMARY KEY,
    strategy_id   TEXT,                       -- 關聯策略（nullable）
    symbol        TEXT NOT NULL,
    trade_type    TEXT NOT NULL,              -- BUY, SELL
    quantity      REAL NOT NULL,
    price         REAL NOT NULL,
    total_amount  REAL NOT NULL,
    fee           REAL NOT NULL DEFAULT 0,
    executed_at   DATETIME NOT NULL,
    created_at    DATETIME NOT NULL,
    FOREIGN KEY (strategy_id) REFERENCES strategies(id) ON DELETE SET NULL
);

-- 通知記錄表
CREATE TABLE notification_logs (
    id            TEXT PRIMARY KEY,
    strategy_id   TEXT NOT NULL,
    trigger_type  TEXT NOT NULL,              -- BUY, SELL
    trigger_price REAL NOT NULL,
    message       TEXT NOT NULL,
    is_read       BOOLEAN NOT NULL DEFAULT 0,
    notified_at   DATETIME NOT NULL,
    FOREIGN KEY (strategy_id) REFERENCES strategies(id) ON DELETE CASCADE
);

-- 持倉表
CREATE TABLE positions (
    id            TEXT PRIMARY KEY,
    symbol        TEXT NOT NULL UNIQUE,       -- 每個幣種只有一筆持倉
    quantity      REAL NOT NULL DEFAULT 0,
    avg_cost      REAL NOT NULL DEFAULT 0,    -- 平均成本
    updated_at    DATETIME NOT NULL
);
```

### 7.2 索引策略

```sql
-- 策略表索引
CREATE INDEX idx_strategies_symbol ON strategies(symbol);
CREATE INDEX idx_strategies_is_active ON strategies(is_active);
CREATE INDEX idx_strategies_symbol_active ON strategies(symbol, is_active);

-- 交易記錄索引
CREATE INDEX idx_trade_records_symbol ON trade_records(symbol);
CREATE INDEX idx_trade_records_executed_at ON trade_records(executed_at);
CREATE INDEX idx_trade_records_strategy_id ON trade_records(strategy_id);

-- 通知記錄索引
CREATE INDEX idx_notification_logs_strategy_id ON notification_logs(strategy_id);
CREATE INDEX idx_notification_logs_notified_at ON notification_logs(notified_at);
```

**索引設計原則：**
- 為高頻查詢的欄位建立索引（如 `symbol`, `is_active`）。
- 為外鍵建立索引（如 `strategy_id`）。
- 考慮複合索引（如 `symbol + is_active`），用於常見的組合查詢。
- 避免過度索引（每新增一個索引都會降低寫入效能）。

### 7.3 資料遷移管理

**使用 [golang-migrate](https://github.com/golang-migrate/migrate) 或 GORM AutoMigrate：**

```go
// 方案 1: GORM AutoMigrate（開發階段）
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &domain.Strategy{},
        &domain.TradeRecord{},
        &domain.NotificationLog{},
        &domain.Position{},
    )
}

// 方案 2: 手動遷移檔案（生產環境）
// migrations/000001_create_strategies_table.up.sql
CREATE TABLE strategies (
    ...
);

// migrations/000001_create_strategies_table.down.sql
DROP TABLE strategies;
```

**遷移原則：**
- 所有 schema 變更必須透過遷移檔案執行。
- 遷移檔案必須同時提供 `up` 和 `down` 腳本。
- 遷移檔案必須與程式碼變更一起進入版本控制。
- 禁止修改已執行的遷移檔案（建立新的遷移檔案來修正）。

### 7.4 事務管理

```go
// Good: 明確的事務邊界
func (s *Service) ExecuteTrade(ctx context.Context, req TradeRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 建立交易記錄
        trade := &domain.TradeRecord{...}
        if err := tx.Create(trade).Error; err != nil {
            return err
        }

        // 2. 更新持倉
        position := &domain.Position{}
        if err := tx.First(position, "symbol = ?", req.Symbol).Error; err != nil {
            return err
        }
        position.UpdateQuantity(trade.Quantity, trade.TradeType)
        if err := tx.Save(position).Error; err != nil {
            return err
        }

        return nil
    })
}
```

**事務使用原則：**
- 需要多個資料表操作的業務邏輯必須使用事務。
- 事務應儘可能短（避免長時間鎖定）。
- 避免在事務內執行外部 API 呼叫。

---

## 8. 業務邏輯層規範 (Business Logic Layer)

### 8.1 Service 層設計模式

**依賴注入模式：**

```go
// Good: 透過介面依賴
type StrategyService struct {
    repo     StrategyRepository    // 介面
    exchange ExchangeClient        // 介面
    notifier Notifier              // 介面
    logger   logger.Logger
}

func NewStrategyService(
    repo StrategyRepository,
    exchange ExchangeClient,
    notifier Notifier,
    logger logger.Logger,
) *StrategyService {
    return &StrategyService{
        repo:     repo,
        exchange: exchange,
        notifier: notifier,
        logger:   logger,
    }
}

// Bad: 直接依賴具體實作
type StrategyService struct {
    db           *gorm.DB           // 具體實作
    binanceClient *binance.Client   // 具體實作
}
```

### 8.2 錯誤處理模式

```go
// 使用 errors.Is 和 errors.As 進行錯誤判斷
func (s *Service) GetStrategy(ctx context.Context, id string) (*domain.Strategy, error) {
    strategy, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrStrategyNotFound
        }
        return nil, fmt.Errorf("failed to find strategy: %w", err)
    }
    return strategy, nil
}

// 使用 fmt.Errorf with %w 來包裝錯誤（保留原始錯誤鏈）
if err := s.repo.Create(ctx, strategy); err != nil {
    return nil, fmt.Errorf("failed to create strategy: %w", err)
}
```

### 8.3 日誌記錄規範

```go
// 使用結構化日誌
s.logger.Info("creating strategy",
    "symbol", req.Symbol,
    "buy_lower", req.BuyLower,
    "sell_upper", req.SellUpper,
)

s.logger.Error("failed to fetch price from exchange",
    "symbol", symbol,
    "error", err,
    "exchange", "binance",
)

// 不要記錄敏感資訊
s.logger.Info("user authenticated")  // Good
s.logger.Info("user authenticated", "api_key", apiKey)  // Bad
```

### 8.4 輸入驗證

```go
// 在 Service 層進行業務層級的驗證
func (s *Service) CreateStrategy(ctx context.Context, req CreateStrategyRequest) error {
    // 1. 基本格式驗證
    if req.Symbol == "" {
        return fmt.Errorf("symbol is required")
    }
    if req.BuyLower <= 0 {
        return fmt.Errorf("buy lower bound must be positive")
    }
    if req.SellUpper <= req.BuyLower {
        return fmt.Errorf("sell upper bound must be greater than buy lower bound")
    }

    // 2. 業務規則驗證
    currentPrice, err := s.exchange.GetPrice(ctx, req.Symbol)
    if err != nil {
        return fmt.Errorf("failed to validate price: %w", err)
    }
    if req.BuyLower > currentPrice*1.5 {
        return fmt.Errorf("buy lower bound is too far from current price")
    }

    // 3. 建立領域物件（領域物件內部也會進行驗證）
    strategy := domain.NewStrategy(req.Symbol, req.BuyLower, req.SellUpper)
    if err := strategy.Validate(); err != nil {
        return err
    }

    return s.repo.Create(ctx, strategy)
}
```

---

## 9. 並發與背景任務 (Concurrency & Background Tasks)

### 9.1 價格監控服務設計

**核心需求：**
- 每 30 秒獲取一次價格資料。
- 支援多個幣種並發監控。
- 能夠優雅關閉（graceful shutdown）。
- 錯誤不應導致整個服務崩潰。

**實作範例：**

```go
type PriceMonitor struct {
    strategyRepo StrategyRepository
    exchange     ExchangeClient
    notifier     Notifier
    logger       logger.Logger
    interval     time.Duration
    stopCh       chan struct{}
    wg           sync.WaitGroup
}

func NewPriceMonitor(
    strategyRepo StrategyRepository,
    exchange ExchangeClient,
    notifier Notifier,
    logger logger.Logger,
    interval time.Duration,
) *PriceMonitor {
    return &PriceMonitor{
        strategyRepo: strategyRepo,
        exchange:     exchange,
        notifier:     notifier,
        logger:       logger,
        interval:     interval,
        stopCh:       make(chan struct{}),
    }
}

func (m *PriceMonitor) Start(ctx context.Context) error {
    m.logger.Info("starting price monitor", "interval", m.interval)

    ticker := time.NewTicker(m.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            m.checkPrices(ctx)
        case <-m.stopCh:
            m.logger.Info("stopping price monitor")
            return nil
        case <-ctx.Done():
            m.logger.Info("context cancelled, stopping price monitor")
            return ctx.Err()
        }
    }
}

func (m *PriceMonitor) checkPrices(ctx context.Context) {
    // 1. 取得所有啟用的策略
    strategies, err := m.strategyRepo.FindActive(ctx)
    if err != nil {
        m.logger.Error("failed to fetch active strategies", "error", err)
        return
    }

    // 2. 按幣種分組
    symbolMap := make(map[string][]*domain.Strategy)
    for _, strategy := range strategies {
        symbolMap[strategy.Symbol] = append(symbolMap[strategy.Symbol], strategy)
    }

    // 3. 並發檢查每個幣種
    var wg sync.WaitGroup
    for symbol, strategies := range symbolMap {
        wg.Add(1)
        go func(sym string, strats []*domain.Strategy) {
            defer wg.Done()
            m.checkSymbol(ctx, sym, strats)
        }(symbol, strategies)
    }
    wg.Wait()
}

func (m *PriceMonitor) checkSymbol(ctx context.Context, symbol string, strategies []*domain.Strategy) {
    // Panic recovery
    defer func() {
        if r := recover(); r != nil {
            m.logger.Error("panic in checkSymbol",
                "symbol", symbol,
                "panic", r,
                "stack", string(debug.Stack()),
            )
        }
    }()

    // 取得當前價格
    price, err := m.exchange.GetPrice(ctx, symbol)
    if err != nil {
        m.logger.Error("failed to get price", "symbol", symbol, "error", err)
        return
    }

    m.logger.Info("fetched price", "symbol", symbol, "price", price)

    // 檢查每個策略
    for _, strategy := range strategies {
        if strategy.ShouldBuy(price) {
            m.notifyTrigger(ctx, strategy, "BUY", price)
        } else if strategy.ShouldSell(price) {
            m.notifyTrigger(ctx, strategy, "SELL", price)
        }
    }
}

func (m *PriceMonitor) notifyTrigger(ctx context.Context, strategy *domain.Strategy, triggerType string, price float64) {
    message := fmt.Sprintf("[%s] %s signal: %s at %.2f",
        triggerType, strategy.Symbol, triggerType, price)

    if err := m.notifier.Send(ctx, message); err != nil {
        m.logger.Error("failed to send notification", "error", err)
        return
    }

    m.logger.Info("notification sent", "strategy_id", strategy.ID, "type", triggerType)
}

func (m *PriceMonitor) Stop() {
    close(m.stopCh)
    m.wg.Wait()
}
```

### 9.2 並發安全原則

```go
// Good: 使用 channel 進行通訊
type Worker struct {
    jobs    chan Job
    results chan Result
}

func (w *Worker) Start() {
    for job := range w.jobs {
        result := process(job)
        w.results <- result
    }
}

// Good: 使用 sync.Mutex 保護共享資源
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// Bad: 直接共享變數（data race）
var count int
go func() { count++ }()
go func() { count++ }()
```

### 9.3 優雅關閉 (Graceful Shutdown)

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 啟動價格監控服務
    monitor := NewPriceMonitor(...)
    go func() {
        if err := monitor.Start(ctx); err != nil {
            log.Fatal(err)
        }
    }()

    // 監聽系統信號
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    <-sigCh
    log.Println("received shutdown signal")

    // 觸發關閉流程
    cancel()

    // 等待所有 goroutine 結束（設定超時）
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    done := make(chan struct{})
    go func() {
        monitor.Stop()
        close(done)
    }()

    select {
    case <-done:
        log.Println("graceful shutdown completed")
    case <-shutdownCtx.Done():
        log.Println("shutdown timeout exceeded")
    }
}
```

---

## 10. 錯誤處理與日誌 (Error Handling & Logging)

### 10.1 錯誤處理模式

**使用 `errors.Is` 和 `errors.As`：**

```go
// 定義領域錯誤
var (
    ErrStrategyNotFound = errors.New("strategy not found")
    ErrInvalidPrice     = errors.New("invalid price")
)

// 使用 errors.Is 檢查錯誤
if errors.Is(err, domain.ErrStrategyNotFound) {
    return nil, fmt.Errorf("strategy not found")
}

// 使用 errors.As 提取特定錯誤類型
var validationErr *domain.ValidationError
if errors.As(err, &validationErr) {
    return nil, fmt.Errorf("validation failed: %w", validationErr)
}
```

**錯誤包裝原則：**

```go
// Good: 使用 %w 包裝錯誤（保留錯誤鏈）
if err := s.repo.Create(ctx, strategy); err != nil {
    return nil, fmt.Errorf("failed to create strategy: %w", err)
}

// Bad: 使用 %v 會丟失錯誤鏈
if err := s.repo.Create(ctx, strategy); err != nil {
    return nil, fmt.Errorf("failed to create strategy: %v", err)
}
```

### 10.2 日誌系統設計

**使用結構化日誌（推薦 `logrus` 或 `zap`）：**

```go
import "github.com/sirupsen/logrus"

type Logger struct {
    logger *logrus.Logger
}

func NewLogger() *Logger {
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{})
    log.SetLevel(logrus.InfoLevel)
    return &Logger{logger: log}
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
    l.logger.WithFields(fields).Info(msg)
}

func (l *Logger) Error(msg string, fields map[string]interface{}) {
    l.logger.WithFields(fields).Error(msg)
}

// 使用範例
logger.Info("fetching price", map[string]interface{}{
    "symbol": "BTC",
    "exchange": "binance",
})

logger.Error("API call failed", map[string]interface{}{
    "symbol": "BTC",
    "error": err.Error(),
    "retry_count": 3,
})
```

**日誌級別使用指引：**

| 級別 | 使用時機 | 範例 |
|------|---------|------|
| **DEBUG** | 開發階段的詳細資訊 | "執行 SQL 查詢: SELECT * FROM strategies" |
| **INFO** | 正常的業務流程資訊 | "成功建立策略", "價格監控已啟動" |
| **WARN** | 可能的問題，但不影響主流程 | "API 回應延遲超過 5 秒", "重試第 2 次" |
| **ERROR** | 錯誤導致功能失敗 | "無法連線到交易所 API", "資料庫寫入失敗" |
| **FATAL** | 系統無法繼續運行 | "無法載入設定檔", "資料庫連線失敗" |

### 10.3 全域錯誤處理（Gin 範例）

```go
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err

            var statusCode int
            var errorCode string

            switch {
            case errors.Is(err, domain.ErrStrategyNotFound):
                statusCode = http.StatusNotFound
                errorCode = "STRATEGY_NOT_FOUND"
            case errors.Is(err, domain.ErrInvalidPrice):
                statusCode = http.StatusBadRequest
                errorCode = "INVALID_PRICE"
            default:
                statusCode = http.StatusInternalServerError
                errorCode = "INTERNAL_ERROR"
            }

            c.JSON(statusCode, Response{
                Success: false,
                Error: &ErrorDetail{
                    Code:    errorCode,
                    Message: err.Error(),
                },
            })
        }
    }
}
```

---

## 11. 測試策略 (Testing Strategy)

### 11.1 測試金字塔

```
          ┌─────────────┐
          │ E2E Tests   │  (5%)   - 整合測試
          │             │
        ┌─────────────────┐
        │ Integration     │  (15%)  - 整合測試
        │ Tests           │
      ┌───────────────────────┐
      │   Unit Tests          │  (80%)  - 單元測試
      │                       │
      └───────────────────────┘
```

### 11.2 單元測試規範

**測試檔案命名：**
```
strategy_service.go       # 被測試檔案
strategy_service_test.go  # 測試檔案
```

**測試函式命名：**
```go
// 格式: Test<StructName>_<MethodName>_<Scenario>
func TestStrategyService_CreateStrategy_Success(t *testing.T) {}
func TestStrategyService_CreateStrategy_WithInvalidPrice(t *testing.T) {}
func TestStrategy_Validate_WithNegativePrice(t *testing.T) {}
```

**使用 Table-Driven Tests：**

```go
func TestStrategy_ShouldBuy(t *testing.T) {
    tests := []struct {
        name         string
        strategy     *domain.Strategy
        currentPrice float64
        want         bool
    }{
        {
            name: "price below buy lower bound",
            strategy: &domain.Strategy{
                BuyLower: 60000,
                IsActive: true,
            },
            currentPrice: 59000,
            want:         true,
        },
        {
            name: "price equal to buy lower bound",
            strategy: &domain.Strategy{
                BuyLower: 60000,
                IsActive: true,
            },
            currentPrice: 60000,
            want:         true,
        },
        {
            name: "price above buy lower bound",
            strategy: &domain.Strategy{
                BuyLower: 60000,
                IsActive: true,
            },
            currentPrice: 61000,
            want:         false,
        },
        {
            name: "strategy is inactive",
            strategy: &domain.Strategy{
                BuyLower: 60000,
                IsActive: false,
            },
            currentPrice: 59000,
            want:         false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.strategy.ShouldBuy(tt.currentPrice)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### 11.3 Mock 與依賴注入

**使用 [gomock](https://github.com/golang/mock) 或手動建立 Mock：**

```go
// 定義介面
type ExchangeClient interface {
    GetPrice(ctx context.Context, symbol string) (float64, error)
}

// 手動建立 Mock（簡單場景）
type MockExchangeClient struct {
    GetPriceFunc func(ctx context.Context, symbol string) (float64, error)
}

func (m *MockExchangeClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
    if m.GetPriceFunc != nil {
        return m.GetPriceFunc(ctx, symbol)
    }
    return 0, errors.New("not implemented")
}

// 測試中使用 Mock
func TestPriceMonitor_CheckSymbol(t *testing.T) {
    mockExchange := &MockExchangeClient{
        GetPriceFunc: func(ctx context.Context, symbol string) (float64, error) {
            return 60000.0, nil
        },
    }

    monitor := NewPriceMonitor(nil, mockExchange, nil, nil, 30*time.Second)
    // ... 進行測試
}
```

### 11.4 整合測試

**測試資料庫操作（使用 SQLite in-memory）：**

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // 執行遷移
    err = db.AutoMigrate(&domain.Strategy{})
    require.NoError(t, err)

    return db
}

func TestStrategyRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := sqlite.NewStrategyRepository(db)

    strategy := &domain.Strategy{
        ID:        "test-id",
        Symbol:    "BTC",
        BuyLower:  60000,
        SellUpper: 70000,
        IsActive:  true,
    }

    err := repo.Create(context.Background(), strategy)
    assert.NoError(t, err)

    // 驗證資料已寫入
    found, err := repo.FindByID(context.Background(), "test-id")
    assert.NoError(t, err)
    assert.Equal(t, strategy.Symbol, found.Symbol)
}
```

### 11.5 測試覆蓋率

**執行測試並產生覆蓋率報告：**

```bash
# 執行所有測試
go test -v ./...

# 產生覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 顯示覆蓋率百分比
go test -cover ./...
```

**覆蓋率目標：**
- 領域層（domain）：> 90%
- 用例層（usecase）：> 85%
- 適配器層（adapter）：> 70%
- 整體專案：> 80%

---

## 12. 安全與性能 (Security & Performance)

### 12.1 輸入驗證與類型安全

**永遠不要信任來自外部的輸入：**

```go
// Good: 驗證所有輸入
func (h *Handler) CreateStrategy(c *gin.Context) {
    var req CreateStrategyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse("invalid request format"))
        return
    }

    // 額外的業務邏輯驗證
    if req.BuyLower <= 0 {
        c.JSON(400, ErrorResponse("buy lower must be positive"))
        return
    }
    if req.SellUpper <= req.BuyLower {
        c.JSON(400, ErrorResponse("sell upper must be greater than buy lower"))
        return
    }

    // 防止 SQL 注入（使用參數化查詢）
    strategy, err := h.service.CreateStrategy(c.Request.Context(), req)
    // ...
}
```

### 12.2 SQL 注入防護

```go
// Good: 使用 GORM 的參數化查詢
db.Where("symbol = ?", userInput).Find(&strategies)

// Good: 使用命名參數
db.Where("symbol = @symbol", sql.Named("symbol", userInput)).Find(&strategies)

// Bad: 字串拼接（容易 SQL 注入）
query := fmt.Sprintf("SELECT * FROM strategies WHERE symbol = '%s'", userInput)
db.Raw(query).Scan(&strategies)
```

### 12.3 敏感資料處理

**設定檔範例：**

```yaml
# config.yaml
exchange:
  binance:
    api_key: ${BINANCE_API_KEY}      # 從環境變數讀取
    api_secret: ${BINANCE_API_SECRET}
database:
  path: "./data/transaction.db"
monitor:
  interval: 30s
```

**載入設定：**

```go
import "github.com/spf13/viper"

type Config struct {
    Exchange ExchangeConfig `mapstructure:"exchange"`
    Database DatabaseConfig `mapstructure:"database"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")

    // 允許環境變數覆蓋
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

**禁止硬編碼敏感資料：**

```go
// Bad: 硬編碼 API Key
const BinanceAPIKey = "abc123xyz"

// Good: 從環境變數讀取
apiKey := os.Getenv("BINANCE_API_KEY")
if apiKey == "" {
    return errors.New("BINANCE_API_KEY not set")
}
```

### 12.4 Rate Limiting（API 呼叫頻率控制）

**使用 `golang.org/x/time/rate` 實作速率限制：**

```go
import "golang.org/x/time/rate"

type BinanceClient struct {
    httpClient *http.Client
    limiter    *rate.Limiter
}

func NewBinanceClient() *BinanceClient {
    return &BinanceClient{
        httpClient: &http.Client{Timeout: 10 * time.Second},
        limiter:    rate.NewLimiter(rate.Every(time.Second), 10), // 每秒最多 10 次請求
    }
}

func (c *BinanceClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
    // 等待速率限制許可
    if err := c.limiter.Wait(ctx); err != nil {
        return 0, fmt.Errorf("rate limiter: %w", err)
    }

    // 執行實際的 API 請求
    resp, err := c.httpClient.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol))
    // ...
}
```

### 12.5 性能優化策略

**1. 使用 Connection Pooling：**

```go
db, err := gorm.Open(sqlite.Open("transaction.db"), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}

sqlDB, err := db.DB()
sqlDB.SetMaxOpenConns(10)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(time.Hour)
```

**2. 批次處理（減少資料庫往返次數）：**

```go
// Good: 一次查詢所有策略
strategies, err := repo.FindActive(ctx)

// Bad: N+1 查詢問題
for _, id := range strategyIDs {
    strategy, err := repo.FindByID(ctx, id)
    // ...
}
```

**3. 使用 Context Timeout：**

```go
func (c *BinanceClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
    // 設定 5 秒超時
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    // ...
}
```

**4. 監控關鍵指標：**

```go
// 記錄 API 呼叫延遲
start := time.Now()
price, err := c.exchange.GetPrice(ctx, symbol)
duration := time.Since(start)

if duration > 5*time.Second {
    logger.Warn("slow API call", "symbol", symbol, "duration", duration)
}
```

---

## 13. 依賴管理 (Dependency Management)

### 13.1 Go Modules 最佳實踐

**初始化專案：**

```bash
go mod init github.com/yourusername/go-transaction
```

**`go.mod` 範例：**

```go
module github.com/yourusername/go-transaction

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.2
    gorm.io/driver/sqlite v1.5.5
    gorm.io/gorm v1.25.7
    github.com/sirupsen/logrus v1.9.3
    github.com/stretchr/testify v1.9.0
)
```

**依賴更新策略：**

```bash
# 查看可更新的依賴
go list -u -m all

# 更新特定依賴
go get -u github.com/gin-gonic/gin

# 更新所有依賴的次版本號
go get -u ./...

# 清理未使用的依賴
go mod tidy
```

### 13.2 第三方庫選擇標準

**選擇依賴套件時，考慮以下因素：**

1. **維護狀態**: 最近更新時間 < 1 年
2. **社群活躍度**: GitHub Stars > 1000, 活躍的 Issue 與 PR
3. **授權協議**: 確認授權符合專案需求（如 MIT, Apache 2.0）
4. **文件品質**: 有清晰的 README 和範例程式碼
5. **測試覆蓋率**: 核心功能有測試覆蓋
6. **相依性**: 避免引入過多傳遞依賴（transitive dependencies）

**推薦套件清單：**

| 用途 | 套件名稱 | 理由 |
|------|---------|------|
| Web 框架 | `gin-gonic/gin` | 高效能、簡單易用、社群活躍 |
| CLI 框架 | `spf13/cobra` | 業界標準、kubectl 使用 |
| 設定管理 | `spf13/viper` | 支援多種格式、環境變數整合 |
| ORM | `gorm.io/gorm` | 功能完整、支援 SQLite |
| 日誌 | `sirupsen/logrus` 或 `uber-go/zap` | 結構化日誌 |
| 測試 | `stretchr/testify` | 豐富的斷言函式 |
| HTTP 客戶端 | 標準庫 `net/http` | 無需額外依賴 |

---

## 14. 配置管理 (Configuration Management)

### 14.1 配置檔案結構

```yaml
# config/config.yaml
app:
  name: "go-transaction"
  version: "1.0.0"
  env: "development"  # development, staging, production

server:
  port: 8080
  host: "0.0.0.0"

database:
  driver: "sqlite"
  path: "./data/transaction.db"

exchange:
  binance:
    base_url: "https://api.binance.com"
    timeout: 10s
    rate_limit: 10  # requests per second

monitor:
  interval: 30s
  symbols:
    - BTC
    - ETH
    - USDT

notification:
  console:
    enabled: true
  telegram:
    enabled: false
    bot_token: ${TELEGRAM_BOT_TOKEN}
    chat_id: ${TELEGRAM_CHAT_ID}

logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, text
  output: "stdout"  # stdout, file
```

### 14.2 環境變數覆蓋

```bash
# .env.example（提交到版本控制）
APP_ENV=development
DATABASE_PATH=./data/transaction.db
BINANCE_API_KEY=your_api_key_here
BINANCE_API_SECRET=your_api_secret_here
LOG_LEVEL=info

# .env（不提交到版本控制，.gitignore 排除）
APP_ENV=development
BINANCE_API_KEY=abc123xyz
BINANCE_API_SECRET=secret123
```

**載入環境變數：**

```go
import "github.com/joho/godotenv"

func init() {
    // 載入 .env 檔案（開發環境）
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}
```

### 14.3 多環境配置

```
config/
├── config.yaml           # 預設配置
├── config.dev.yaml       # 開發環境覆蓋
├── config.staging.yaml   # 測試環境覆蓋
└── config.prod.yaml      # 生產環境覆蓋
```

**載入邏輯：**

```go
func LoadConfig() (*Config, error) {
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    viper.SetConfigName("config")
    viper.AddConfigPath("./config")
    viper.ReadInConfig()

    // 覆蓋環境特定配置
    viper.SetConfigName(fmt.Sprintf("config.%s", env))
    viper.MergeInConfig()

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

---

## 15. 部署與運維 (Deployment & Operations)

### 15.1 建置流程

**`Makefile` 範例：**

```makefile
.PHONY: build test lint clean run

# 變數定義
APP_NAME=go-transaction
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 建置可執行檔
build:
	go build $(LDFLAGS) -o bin/$(APP_NAME) cmd/cli/main.go

# 執行測試
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 執行 linter
lint:
	golangci-lint run

# 格式化程式碼
fmt:
	go fmt ./...

# 清理建置產物
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 執行應用程式
run: build
	./bin/$(APP_NAME)

# 安裝依賴
deps:
	go mod download
	go mod tidy

# 資料庫遷移
migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down
```

### 15.2 Docker 容器化（可選）

**`Dockerfile` 範例：**

```dockerfile
# 建置階段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 複製依賴檔案
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼
COPY . .

# 建置可執行檔
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-transaction cmd/cli/main.go

# 執行階段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 複製可執行檔
COPY --from=builder /go-transaction .

# 複製配置檔案
COPY config/ ./config/

# 暴露埠號（如果有 HTTP 服務）
EXPOSE 8080

CMD ["./go-transaction", "monitor", "start"]
```

**`docker-compose.yml` 範例：**

```yaml
version: '3.8'

services:
  app:
    build: .
    container_name: go-transaction
    environment:
      - APP_ENV=production
      - BINANCE_API_KEY=${BINANCE_API_KEY}
      - BINANCE_API_SECRET=${BINANCE_API_SECRET}
    volumes:
      - ./data:/root/data
      - ./config:/root/config
    restart: unless-stopped
```

### 15.3 日誌輸出標準化

**結構化日誌範例：**

```json
{
  "time": "2025-11-04T10:30:00Z",
  "level": "info",
  "msg": "fetched price",
  "symbol": "BTC",
  "price": 60000.5,
  "exchange": "binance",
  "duration_ms": 234
}
```

**日誌輪換（Log Rotation）：**

```go
import "gopkg.in/natefinch/lumberjack.v2"

func setupLogger() *logrus.Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})

    // 設定日誌輪換
    logger.SetOutput(&lumberjack.Logger{
        Filename:   "/var/log/go-transaction/app.log",
        MaxSize:    100, // MB
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true,
    })

    return logger
}
```

### 15.4 優雅關閉流程

```go
func main() {
    // ... 初始化服務

    // 監聽系統信號
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    <-sigCh
    logger.Info("received shutdown signal")

    // 停止接受新請求
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 關閉 HTTP 伺服器（如果有）
    if err := server.Shutdown(ctx); err != nil {
        logger.Error("server shutdown error", "error", err)
    }

    // 停止背景任務
    monitor.Stop()

    // 關閉資料庫連線
    sqlDB, _ := db.DB()
    sqlDB.Close()

    logger.Info("graceful shutdown completed")
}
```

### 15.5 健康檢查與監控

**健康檢查端點：**

```go
// HTTP 健康檢查
func HealthCheckHandler(c *gin.Context) {
    status := map[string]string{
        "status": "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    }
    c.JSON(200, status)
}

// 註冊路由
r.GET("/health", HealthCheckHandler)
```

**關鍵指標監控：**

```go
type Metrics struct {
    PriceCheckCount   int64
    NotificationsSent int64
    APIErrors         int64
    LastCheckTime     time.Time
}

// 定期輸出監控指標
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        logger.Info("system metrics",
            "price_checks", metrics.PriceCheckCount,
            "notifications", metrics.NotificationsSent,
            "api_errors", metrics.APIErrors,
        )
    }
}()
```

---

## 附錄 A: 程式碼審查檢查清單 (Code Review Checklist)

### 功能正確性
- [ ] 程式碼實作了所有驗收條件
- [ ] 業務邏輯與需求一致
- [ ] 邊界條件已處理（如空值、零值、極值）
- [ ] 錯誤處理完整（所有錯誤都有處理路徑）

### 程式碼品質
- [ ] 命名清晰且符合慣例
- [ ] 函式長度合理（< 50 行）
- [ ] 複雜度可接受（圈複雜度 < 15）
- [ ] 無重複程式碼（DRY 原則）
- [ ] 註解適當（解釋「為什麼」而非「是什麼」）

### 測試
- [ ] 單元測試覆蓋率 > 80%
- [ ] 測試涵蓋正常與異常情況
- [ ] 測試名稱清晰描述測試場景
- [ ] 使用 Mock/Stub 隔離外部依賴

### 效能與安全
- [ ] 無明顯效能瓶頸（如 N+1 查詢）
- [ ] 輸入驗證完整
- [ ] 無 SQL 注入風險
- [ ] 敏感資料未硬編碼

### 架構設計
- [ ] 遵循 Clean Architecture 分層
- [ ] 依賴方向正確（向內依賴）
- [ ] 介面設計合理（單一職責）
- [ ] 無循環依賴

### 資料庫
- [ ] 資料庫變更有遷移檔案
- [ ] 索引設計合理
- [ ] 事務邊界正確
- [ ] 查詢效率可接受

### 文件
- [ ] 複雜邏輯有註解說明
- [ ] API 變更有更新文件
- [ ] README 保持最新
- [ ] CHANGELOG 記錄變更

---

## 附錄 B: 常見問題與解決方案 (FAQ)

### Q1: 如何選擇 GORM 或純 SQL？

**建議：**
- 簡單 CRUD 操作使用 GORM（提高開發效率）
- 複雜查詢（如多表 JOIN、子查詢）使用純 SQL（更好的效能與可控性）

```go
// GORM: 簡單查詢
var strategies []domain.Strategy
db.Where("symbol = ? AND is_active = ?", "BTC", true).Find(&strategies)

// 純 SQL: 複雜查詢
query := `
    SELECT s.*, COUNT(t.id) as trade_count
    FROM strategies s
    LEFT JOIN trade_records t ON s.id = t.strategy_id
    WHERE s.is_active = true
    GROUP BY s.id
    HAVING trade_count > 5
`
db.Raw(query).Scan(&results)
```

### Q2: 如何處理長時間運行的 goroutine 錯誤？

**解決方案：**
- 使用 `defer recover()` 防止 panic 導致程式崩潰
- 記錄詳細的錯誤日誌（包含 stack trace）
- 實作重試機制
- 監控 goroutine 狀態

```go
func (m *Monitor) safeRun(ctx context.Context) {
    defer func() {
        if r := recover(); r != nil {
            m.logger.Error("panic recovered",
                "panic", r,
                "stack", string(debug.Stack()),
            )
            // 重啟 goroutine
            go m.safeRun(ctx)
        }
    }()

    m.run(ctx)
}
```

### Q3: 如何測試涉及時間的功能（如定時任務）？

**解決方案：使用可注入的時間介面**

```go
// 定義時間介面
type Clock interface {
    Now() time.Time
    Sleep(d time.Duration)
}

// 生產環境使用真實時間
type RealClock struct{}

func (c RealClock) Now() time.Time {
    return time.Now()
}

func (c RealClock) Sleep(d time.Duration) {
    time.Sleep(d)
}

// 測試環境使用假時間
type FakeClock struct {
    CurrentTime time.Time
}

func (c *FakeClock) Now() time.Time {
    return c.CurrentTime
}

func (c *FakeClock) Sleep(d time.Duration) {
    c.CurrentTime = c.CurrentTime.Add(d)
}

// 在服務中注入
type Service struct {
    clock Clock
}

func (s *Service) CheckExpired() {
    now := s.clock.Now()
    // ...
}
```

---

## 附錄 C: 參考資源 (References)

### 官方文件
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### 風格指南
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)

### 設計模式
- [Go Patterns](https://github.com/tmrts/go-patterns)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)

### 測試
- [Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

---

## 結語

這份 `rule.md` 是活文件（Living Document），隨著專案發展與團隊經驗積累，應持續更新與改進。

**核心精神：**
- 簡潔優於複雜
- 可測試性優先
- 顯式優於隱式
- 先讓程式碼正確執行，再優化效能

**下一步行動：**
1. 團隊成員審閱此文件，提出疑問與建議
2. 根據實際開發經驗調整規範
3. 在 Code Review 中引用此文件的章節
4. 定期（每季度）回顧與更新

**記住：好的規範是為了提升團隊協作效率，而非限制創新。當規範與實際需求衝突時，請提出討論並更新規範。**

---

**文件結束**

如有任何疑問或建議，請隨時提出！
