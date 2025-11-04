# 架構文件

## 概述

本專案採用 **清潔架構（Clean Architecture）** 設計模式，將應用程序分為多個獨立的層次，確保高內聚、低耦合，便於測試和維護。

## 系統架構

```
┌─────────────────────────────────────────────────────────────────┐
│                    應用程序 (cmd/cli/main.go)                    │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│               介面層 (Interface Layer)                           │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  CLI 命令層                                                │ │
│  │  - CreateStrategy, ListStrategies, GetStrategy            │ │
│  │  - UpdateStrategy, DeleteStrategy, ToggleStrategy         │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│              業務邏輯層 (Use Case / Application Layer)           │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  StrategyService                                           │ │
│  │  - 處理策略的業務規則和流程                              │ │
│  │  - 依賴注入：Repository, Logger                          │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│                領域層 (Domain Layer)                             │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Strategy 實體                                             │ │
│  │  - ID, Symbol, BuyLower, SellUpper, IsActive             │ │
│  │  - Validate(), ShouldBuy(), ShouldSell()                 │ │
│  │                                                            │ │
│  │  領域錯誤                                                  │ │
│  │  - ErrInvalidStrategy, ErrInvalidPrice                   │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│            適配器層 (Adapter / Infrastructure Layer)            │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Repository 層                                             │ │
│  │  - IStrategyRepository 介面                               │ │
│  │  - SQLite Repository 實現                                 │ │
│  │  - CRUD 操作：Create, Read, Update, Delete               │ │
│  │                                                            │ │
│  │  數據庫層                                                  │ │
│  │  - GORM ORM 框架                                          │ │
│  │  - SQLite 驅動                                            │ │
│  │  - 自動遷移（AutoMigrate）                               │ │
│  │                                                            │ │
│  │  日誌層                                                    │ │
│  │  - SimpleLogger 實現                                      │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

## 目錄結構

```
project/
├── cmd/
│   └── cli/
│       └── main.go              # 應用程序入口點
├── internal/
│   ├── domain/                  # 領域層
│   │   ├── strategy.go          # Strategy 實體
│   │   ├── errors.go            # 領域錯誤
│   │   └── strategy_test.go     # 領域層測試
│   ├── usecase/
│   │   └── strategy/            # 業務邏輯層
│   │       ├── service.go       # StrategyService 實現
│   │       ├── service_test.go  # 業務層測試
│   │       └── dto.go           # 數據傳輸對象
│   ├── adapter/
│   │   └── repository/
│   │       ├── repository.go    # Repository 介面
│   │       └── sqlite/
│   │           ├── strategy_repo.go        # SQLite 實現
│   │           ├── strategy_repo_test.go   # Repository 測試
│   │           └── migration.go            # 數據庫遷移
│   └── interface/
│       └── cli/
│           ├── strategy_cmd.go  # CLI 命令實現
│           └── root.go          # CLI 根命令
├── pkg/
│   └── logger/
│       ├── logger.go            # Logger 實現
│       └── logger_test.go       # Logger 測試
├── test/
│   └── integration/
│       └── cli_strategy_test.go # 集成測試
├── docs/
│   ├── architecture.md          # 本文件
│   └── api.md                   # API 文檔
├── go.mod                       # Go 模塊定義
├── go.sum                       # 依賴項校驗
└── tasks.md                     # 開發任務清單
```

## 各層職責

### 領域層 (Domain Layer)

- **目的**：定義業務領域的核心概念
- **包含內容**：
  - 實體（Entity）：Strategy
  - 值對象：價格範圍
  - 領域錯誤
- **特點**：
  - 完全獨立，不依賴任何其他層
  - 包含業務規則的驗證邏輯
  - 純粹的 Go 代碼，無框架依賴

### 業務邏輯層 (Use Case Layer)

- **目的**：實現應用程序的業務流程
- **包含內容**：
  - StrategyService：協調領域邏輯和數據訪問
  - DTO（Data Transfer Objects）：定義請求和響應結構
- **特點**：
  - 依賴反轉：依賴抽象（Repository 介面、Logger 介面）
  - 包含應用程序特定的業務規則
  - 獨立於表示層和持久化層

### 適配器層 (Adapter / Infrastructure Layer)

- **目的**：實現與外部系統的交互
- **包含內容**：
  - Repository：數據訪問層
  - Logger：日誌記錄
  - 數據庫：SQLite
- **特點**：
  - 實現 Repository 和 Logger 介面
  - 處理數據庫操作、日誌輸出等細節
  - 框架和技術棧的具體實現

### 介面層 (Interface Layer)

- **目的**：提供用戶交互的入口
- **包含內容**：
  - CLI 命令：使用 Cobra 框架
  - 命令行參數解析
  - 用戶輸出格式化
- **特點**：
  - 最外層，直接與用戶交互
  - 依賴業務邏輯層
  - 處理用戶輸入的驗證和轉換

## 數據流

### 建立策略的流程

```
用戶輸入 CLI 命令
  ↓
CLI 層解析命令和參數
  ↓
调用 StrategyService.CreateStrategy()
  ↓
StrategyService 驗證輸入，創建 Strategy 實體
  ↓
調用領域層 Strategy.Validate()
  ↓
調用 Repository.Create()
  ↓
Repository 使用 GORM 保存到 SQLite
  ↓
返回結果給 CLI 層
  ↓
CLI 層格式化輸出給用戶
```

### 查詢策略的流程

```
用戶輸入 CLI 命令
  ↓
CLI 層解析命令
  ↓
調用 StrategyService.GetStrategy()
  ↓
調用 Repository.FindByID()
  ↓
Repository 從 SQLite 查詢
  ↓
返回 Strategy 實體
  ↓
StrategyService 轉換為 StrategyResponse
  ↓
返回給 CLI 層
  ↓
CLI 層格式化輸出給用戶
```

## 設計原則

### 1. 依賴反轉原則 (Dependency Inversion Principle)

- 高層模塊不依賴低層模塊，都依賴抽象
- StrategyService 依賴 IStrategyRepository 介面，而非具體實現
- 便於測試時使用 Mock

### 2. 單一職責原則 (Single Responsibility Principle)

- 每個模塊只有一個改變的理由
- Strategy：業務規則
- StrategyService：業務流程
- Repository：數據訪問

### 3. 開閉原則 (Open/Closed Principle)

- 對擴展開放，對修改關閉
- 可以輕鬆添加新的 Repository 實現（如 PostgreSQL）
- 而不需要修改 StrategyService

### 4. 介面分離原則 (Interface Segregation Principle)

- 定義精確的介面
- IStrategyRepository 只包含必要的方法
- Logger 介面簡潔清晰

## 測試策略

### 層級測試

1. **單元測試**（Unit Tests）
   - 領域層：測試 Strategy 實體和驗證邏輯
   - 業務層：使用 Mock Repository 測試 StrategyService
   - Repository 層：使用 in-memory SQLite 測試數據訪問

2. **集成測試**（Integration Tests）
   - 測試真實的 SQLite 數據庫操作
   - 測試完整的業務流程
   - 測試 CLI 命令與服務層的集成

### 測試覆蓋率目標

- 領域層：> 90% (當前: 100%)
- 業務層：> 85% (當前: 96.3%)
- Repository 層：> 70% (當前: 74.2%)
- 整體：> 80% (當前: 達到目標)

## 依賴項

### 核心依賴

- **gorm.io/gorm**：ORM 框架
- **gorm.io/driver/sqlite**：SQLite 驅動
- **github.com/google/uuid**：UUID 生成
- **github.com/spf13/cobra**：CLI 框架

### 測試依賴

- **github.com/stretchr/testify**：測試斷言和 Mock
- **gorm.io/gorm**：in-memory SQLite 測試數據庫

## 未來的改進方向

1. **增強 Logger**
   - 支援 logrus 或 zap 進行更強大的日誌管理
   - 支援不同的日誌級別（DEBUG, INFO, WARN, ERROR）

2. **添加新的 Repository 實現**
   - PostgreSQL
   - MongoDB
   - 基於內存的實現（用於緩存）

3. **擴展 CLI 功能**
   - 添加配置文件支援
   - 添加交互式 CLI 模式
   - 添加 JSON 輸出格式

4. **API 層**
   - 添加 REST API（使用 Gin 或 Echo）
   - 添加 gRPC 服務

5. **持久化層增強**
   - 添加事務支援
   - 添加查詢優化
   - 添加數據庫連接池配置

## 架構決策記錄

### 為什麼選擇清潔架構？

1. **可測試性**：各層獨立，易於單元測試
2. **可維護性**：清晰的職責分離
3. **可擴展性**：易於添加新功能或替換實現
4. **獨立於框架**：業務邏輯不依賴特定技術

### 為什麼選擇 SQLite？

1. **簡單**：無需配置，開箱即用
2. **輕量級**：適合小型應用
3. **便於測試**：支援 in-memory 數據庫
4. **易於部署**：單個文件即可

### 為什麼使用 GORM？

1. **功能完整**：支援複雜查詢和遷移
2. **易用**：簡潔的 API
3. **多數據庫支援**：易於遷移到其他數據庫
4. **活躍社區**：良好的文檔和支援

