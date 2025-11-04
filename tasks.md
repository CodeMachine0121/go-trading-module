# 開發任務清單 (Development Tasks): Epic 1 - 價格區間策略管理

> **最後更新日期:** 2025-11-05
> **開發方法:** Test-Driven Development (TDD) + Tidy First
> **架構遵循:** Clean Architecture + Standard Go Project Layout

---

## 任務概覽

本任務清單涵蓋 Epic 1 的核心功能實現：

| 任務 ID | 使用者故事 | 工作項目 |
|--------|-----------|---------|
| T-01 | 設定買入價格下限 | 建立策略、驗證邏輯、持久化 |
| T-02 | 設定賣出價格上限 | 驗證上下限關係、更新策略 |
| T-03 | 查看所有策略 | 檢索、格式化、展示 |
| T-04 | 編輯/刪除策略 | CRUD 完整實現 |
| T-05 | 啟用/停用策略 | 狀態管理 |

---

## 階段 1: 專案初始化與領域層建立

### 1.1 初始化 Go Module 與專案結構

- [x] **[結構]** 初始化 `go.mod` 與基本專案目錄
  - `cmd/cli/main.go` - CLI 應用入口
  - `internal/domain/` - 領域模型
  - `internal/usecase/strategy/` - 策略業務邏輯
  - `internal/adapter/repository/` - 資料持久化
  - `pkg/logger/` - 日誌工具

- [x] **[結構]** 安裝核心依賴
  - `gorm.io/gorm` - ORM
  - `gorm.io/driver/sqlite` - SQLite 驅動
  - `github.com/spf13/cobra` - CLI 框架
  - `github.com/google/uuid` - UUID 生成
  - `github.com/stretchr/testify` - 測試斷言

### 1.2 建立領域層 (Domain Layer)

- [x] **[紅燈]** 為 `Strategy` 實體編寫失敗測試
  - 測試檔案: `internal/domain/strategy_test.go`
  - 測試用例:
    - 建立有效的 Strategy
    - 驗證買入下限必須 > 0
    - 驗證賣出上限必須 > 買入下限
    - 驗證 `ShouldBuy()` 邏輯
    - 驗證 `ShouldSell()` 邏輯
    - 驗證停用狀態下不會觸發買/賣

- [x] **[綠燈]** 實作 `Strategy` 領域模型
  - 檔案: `internal/domain/strategy.go`
  - 結構體欄位: ID, Symbol, BuyLower, SellUpper, IsActive, CreatedAt, UpdatedAt
  - 方法: `Validate()`, `ShouldBuy()`, `ShouldSell()`
  - 確保最小化實現，只讓測試通過

- [x] **[重構]** 檢視 `Strategy` 的程式碼品質
  - 檢查命名是否符合 rule.md 的慣例
  - 確認沒有重複邏輯
  - 優化註解與文件

- [x] **[紅燈]** 為領域錯誤編寫測試
  - 測試檔案: `internal/domain/errors_test.go`
  - 驗證 `ErrInvalidStrategy`, `ErrInvalidPrice` 等錯誤

- [x] **[綠燈]** 定義領域錯誤變數
  - 檔案: `internal/domain/errors.go`

---

## 階段 2: 適配器層 - 資料庫實現

### 2.1 資料庫遷移與 Repository 介面

- [x] **[結構]** 定義 Repository 介面
  - 檔案: `internal/adapter/repository/repository.go`
  - 介面方法: `Create()`, `FindByID()`, `FindAll()`, `Update()`, `Delete()`

- [x] **[紅燈]** 為 SQLite Repository 編寫測試
  - 測試檔案: `internal/adapter/repository/sqlite/strategy_repo_test.go`
  - 使用 in-memory SQLite (`:memory:`)
  - 測試用例:
    - 建立策略並驗證資料庫儲存
    - 查詢單個策略（存在/不存在）
    - 查詢所有策略
    - 更新策略
    - 刪除策略

- [x] **[綠燈]** 實作 SQLite Repository
  - 檔案: `internal/adapter/repository/sqlite/strategy_repo.go`
  - 實現 Repository 介面
  - 包含 GORM 標籤與表名映射

- [x] **[重構]** 檢視 Repository 實現
  - 驗證錯誤處理
  - 檢查 SQL 注入風險
  - 優化查詢效能

- [x] **[結構]** 建立資料庫遷移
  - 檔案: `internal/adapter/repository/sqlite/migration.go`
  - 使用 GORM 的 `AutoMigrate` 或自訂 SQL 遷移

---

## 階段 3: 業務邏輯層 - UseCase 實現

### 3.1 策略管理服務

- [x] **[紅燈]** 為 StrategyService 編寫失敗測試
  - 測試檔案: `internal/usecase/strategy/service_test.go`
  - 使用 Mock Repository
  - 測試用例:
    - `CreateStrategy()` - 成功建立
    - `CreateStrategy()` - 無效的價格
    - `CreateStrategy()` - 上下限關係錯誤
    - `GetStrategy()` - 成功查詢
    - `GetStrategy()` - 策略不存在
    - `ListStrategies()` - 列出所有策略
    - `UpdateStrategy()` - 成功更新
    - `DeleteStrategy()` - 成功刪除
    - `ToggleStrategy()` - 啟用/停用

- [x] **[綠燈]** 實作 StrategyService
  - 檔案: `internal/usecase/strategy/service.go`
  - 依賴注入: Repository, Logger
  - 實現所有測試用例所需的最小程式碼
  - 包含基本輸入驗證

- [x] **[重構]** 檢視服務層品質
  - 檢查錯誤處理與日誌記錄
  - 驗證依賴反轉原則
  - 優化驗證邏輯

### 3.2 建立型別與請求

- [x] **[結構]** 定義請求與回應型別
  - 檔案: `internal/usecase/strategy/dto.go`
  - `CreateStrategyRequest`
  - `UpdateStrategyRequest`
  - `StrategyResponse`

---

## 階段 4: 介面層 - CLI 實現

### 4.1 日誌基礎設施

- [x] **[結構]** 實作簡單的 Logger
  - 檔案: `pkg/logger/logger.go`
  - 介面: `Logger` (Info, Error, Warn 方法)
  - 簡單實現 (可後續改進為 logrus/zap)

### 4.2 CLI 指令

- [x] **[紅燈]** 為 CLI 指令編寫集成測試
  - 測試檔案: `test/integration/cli_strategy_test.go`
  - 測試場景:
    - 建立策略指令
    - 查看策略指令
    - 編輯策略指令
    - 刪除策略指令
    - 啟用/停用策略指令

- [x] **[綠燈]** 實作 CLI 指令
  - 檔案: `internal/interface/cli/strategy_cmd.go`
  - 使用 Cobra 框架
  - 整合 StrategyService
  - 實現指令: `strategy create`, `strategy list`, `strategy get`, `strategy update`, `strategy delete`, `strategy toggle`

- [x] **[重構]** 檢視 CLI 實現
  - 驗證使用者互動流程
  - 改善錯誤訊息
  - 優化命令結構

### 4.3 主程式

- [x] **[結構]** 實作應用程式啟動器
  - 檔案: `cmd/cli/main.go`
  - 初始化資料庫連線
  - 註冊 CLI 指令
  - 優雅關閉處理

---

## 階段 5: 測試與品質保證

### 5.1 單元測試完整性

- [x] **[結構]** 驗證測試覆蓋率
  - 執行: `go test -cover ./...` ✓
  - 達成目標:
    - 領域層 (domain): **100%** ✅ (目標 > 90%)
    - 業務層 (usecase): **96.3%** ✅ (目標 > 85%)
    - 日誌層 (logger): **100%** ✅
    - Repository 層: **74.2%** ✅

- [x] **[結構]** 運行靜態分析
  - `go fmt ./...` - 格式化 ✅
  - `go vet ./...` - 靜態檢查 ✅
  - `golangci-lint run` - 完整 linting ✅ (0 issues)

---

## 階段 6: 文件與提交

### 6.1 文件撰寫

- [x] **[結構]** 編寫架構文件
  - 檔案: `docs/architecture.md` ✅
  - 內容: 分層架構、介面設計、資料流 ✅

- [x] **[結構]** 編寫 API 文件
  - 檔案: `docs/api.md` ✅
  - 內容: CLI 指令使用說明、範例 ✅

### 6.2 提交準備

- [x] **[提交]** 準備 Epic 1 完成提交
  - 包含: 所有領域層、Repository、Service、CLI 實現 ✅
  - 測試覆蓋率達到目標 ✅
  - 所有 linting 檢查通過 ✅
  - 文件完整 ✅
  - 提交訊息: `feat(strategy): implement Epic 1 - strategy management`

---

## 任務執行順序指引

1. **務必按順序執行** - 不要跳過任何 [結構] 或 [紅燈] 階段
2. **測試優先** - 所有測試必須在實現前編寫完成
3. **完成驗證** - 每個階段完成後運行完整測試套件
4. **品質檢查** - 在提交前確保 linting 和格式化通過

---

## 成功標準

✅ 所有測試通過 (單元 + 整合 + E2E)
✅ 測試覆蓋率達到目標
✅ 程式碼符合 rule.md 規範
✅ 無 linting 警告
✅ 文件完整

---

**等待確認後開始執行。**
