project/
│
├── config.yaml                // 全域設定檔，存放預設參數與工作時間等配置
│
├── cmd/
│   ├── main/                    
│   │   ├── main.go               // UI 主程式入口
│   │   ├── ui_handler.go         // UI 元件初始化與事件綁定
│   │   │   ├── initUI()          // 建立主視窗、表單、按鈕等
│   │   │   └── bindUIEvents()    // 綁定按鈕點擊、設定更新等事件
│   │   └── event_listener.go     // 監聽 UI 與系統/後台的事件通知
│   │
│   └── daemon/                  
│       ├── main.go               // 常駐程式入口
│       └── daemon_controller.go  // 控制常駐模組啟動/停止/重啟
│           ├── StartDaemon()     // 啟動防閒置流程與健康檢查
│           ├── StopDaemon()      // 停止防閒置模組
│           ├── RestartDaemon()   // 配置更新後重啟防閒置流程
│           └── health_check()    // 定期檢查防閒置功能的健康狀態
│
├── internal/
│   ├── config/                  
│   │   ├── config.go             // 定義 Config 結構與全域設定管理
│   │   │   ├── LoadConfig()      // 讀取並反序列化設定檔 (包含 config.yaml)
│   │   │   ├── SaveConfig()      // 保存設定檔（原子寫入）
│   │   │   └── ValidateConfig()  // 驗證各欄位格式與範圍
│   │   ├── parser.go             // 支援 JSON、YAML 等格式解析
│   │   └── config_test.go        // Config 模組單元測試
│   │
│   ├── preventidle/             
│   │   ├── idle_controller.go    // 整合防止閒置邏輯
│   │   │   ├── StartIdlePrevention()   // 啟動防閒置作業
│   │   │   ├── StopIdlePrevention()    // 停止模擬操作
│   │   │   └── RestartIdlePrevention() // 設定更新後重啟流程
│   │   ├── input_simulator.go    // 模擬輸入操作
│   │   │   ├── SimulateKeyPress()  // 模擬鍵盤按鍵
│   │   │   ├── SimulateMouseMove() // 模擬滑鼠移動
│   │   │   └── SimulateActivity()  // 整合輸入模擬操作
│   │   ├── windows_api.go        // 封裝 Windows API 呼叫
│   │   │   ├── CallSendInput()   // 執行 SendInput 呼叫
│   │   │   └── GetIdleTime()     // 查詢系統閒置時間
│   │   ├── linux_api.go          // 封裝 Linux API 呼叫
│   │   │   ├── CallSendInput()   // 執行 SendInput 呼叫
│   │   │   └── GetIdleTime()     // 查詢系統閒置時間
│   │   └── error_handler.go      // 統一錯誤處理與重試機制
│   │
│   └── schedule/                
│       ├── scheduler.go          // 時間排程管理與任務觸發
│       │   ├── InitScheduler()   // 初始化排程器
│       │   ├── CheckWorkTime()   // 判斷是否處於工作時間
│       │   ├── ScheduleTask()    // 排程觸發防閒置操作
│       │   └── StopScheduler()   // 終止排程器
│       ├── time_manager.go       // 時間處理輔助函式
│       │   ├── ParseTimeString() // 將字串轉成標準時間格式
│       │   ├── IsTimeInRange()   // 檢查是否在指定時間區間
│       │   └── CalculateNextInterval() // 計算下一次間隔
│       └── scheduler_test.go     // 排程邏輯單元測試
│
├── pkg/                        
│   ├── logger/                
│   │   ├── logger.go           // 日誌系統封裝與初始化
│   │   │   ├── InitLogger()    // 初始化 logger
│   │   │   ├── LogInfo()       // 記錄 Info 級別日誌
│   │   │   ├── LogDebug()      // 記錄 Debug 級別日誌
│   │   │   └── LogError()      // 記錄 Error 級別日誌
│   │   └── logger_test.go      // 日誌模組測試
│   │
│   └── utils/                 
│       ├── timeutils.go        // 時間格式化與解析工具
│       │   ├── FormatTime()    // 格式化 time.Time 為字串
│       │   └── ParseDuration() // 將字串轉換為 time.Duration
│       ├── fileutils.go        // 檔案讀寫工具
│       │   ├── ReadFile()      // 讀取檔案內容
│       │   └── WriteFile()     // 寫入資料到檔案
│       └── utils_test.go       // 通用工具函式單元測試
│
├── docs/                      
│   ├── architecture.md         // 系統架構總覽
│   ├── design.md               // 詳細設計與介面規範
│   └── usage.md                // 安裝與使用指南
│
├── Makefile                    // 一鍵編譯、測試與部署
├── go.mod                      // Go module 管理檔案
├── go.sum                      // 依賴鎖定檔案
└── README.md                   // 專案簡介與快速上手指南
