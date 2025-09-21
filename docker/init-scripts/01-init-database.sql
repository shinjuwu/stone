-- GameHub 核心表結構初始化腳本
-- 此腳本創建系統核心表，在遊戲數據表之前執行

-- 切換到目標數據庫
\c gamehub_dev;

-- 創建基礎表結構示例（根據你的實際需求調整）

-- 用戶表
CREATE TABLE IF NOT EXISTS members (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL,
    agent_id INTEGER NOT NULL DEFAULT 0,
    token VARCHAR(255),
    nick_name VARCHAR(100),
    trans_name VARCHAR(100),
    gold DECIMAL(15,4) DEFAULT 0,
    lock_gold DECIMAL(15,4) DEFAULT 0,
    icon_id INTEGER DEFAULT 1,
    status INTEGER DEFAULT 1,
    level_code VARCHAR(50),
    icon_list TEXT,
    is_re_login BOOLEAN DEFAULT FALSE,
    wallet_type INTEGER DEFAULT 0,
    not_kill_dive_cal INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 遊戲記錄表
CREATE TABLE IF NOT EXISTS game_records (
    id BIGSERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL,
    table_id INTEGER NOT NULL,
    user_id BIGINT NOT NULL,
    bet_amount DECIMAL(15,4) DEFAULT 0,
    win_amount DECIMAL(15,4) DEFAULT 0,
    game_result TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 代理表
CREATE TABLE IF NOT EXISTS agents (
    id SERIAL PRIMARY KEY,
    agent_name VARCHAR(100) NOT NULL,
    status INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 遊戲設置表
CREATE TABLE IF NOT EXISTS game_settings (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL,
    game_code VARCHAR(50) NOT NULL,
    game_name VARCHAR(100),
    status INTEGER DEFAULT 1,
    config JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 創建索引

-- 插入初始數據
INSERT INTO agents (id, agent_name) VALUES (1, 'Default Agent') ON CONFLICT DO NOTHING;

INSERT INTO game_settings (game_id, game_code, game_name, status) VALUES
(1001, 'baccarat', '百家樂', 1),
(1002, 'blackjack', '21點', 1),
(1003, 'roulette', '輪盤', 1),
(2001, 'texas', '德州撲克', 1),
(3001, 'rcfishing', '捕魚遊戲', 1),
(4001, 'fruitslot', '水果老虎機', 1)
ON CONFLICT DO NOTHING;

-- 創建更新時間戳的函數
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 創建觸發器
CREATE TRIGGER update_members_updated_at 
    BEFORE UPDATE ON members 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_settings_updated_at 
    BEFORE UPDATE ON game_settings 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 授權給數據庫用戶
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO gamehub_dev;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO gamehub_dev;

-- 輸出核心表創建完成訊息
SELECT 'GameHub 核心表結構創建完成' as initialization_status;