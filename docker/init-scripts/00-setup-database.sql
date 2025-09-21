-- GameHub 數據庫設置腳本
-- 此腳本在其他初始化腳本之前執行

-- 切換到目標數據庫
\c gamehub_dev;

-- 創建必要的擴展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 設置時區
SET timezone = 'Asia/Taipei';

-- 創建更新時間戳的通用函數
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 輸出初始化開始訊息
SELECT 'GameHub 數據庫初始化開始...' as initialization_status;