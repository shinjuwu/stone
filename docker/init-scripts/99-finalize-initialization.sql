-- GameHub 數據庫初始化完成腳本
-- 此腳本在所有其他初始化腳本之後執行

-- 切換到目標數據庫
\c gamehub_dev;

-- 更新數據庫統計資訊
ANALYZE;

-- 創建數據庫完整性檢查視圖
CREATE OR REPLACE VIEW v_database_summary AS
SELECT 
    'gamelist' as table_name,
    COUNT(*) as record_count,
    COUNT(DISTINCT game_code) as unique_games
FROM gamelist
WHERE status = 1
UNION ALL
SELECT 
    'gameinfo' as table_name,
    COUNT(*) as record_count,
    COUNT(DISTINCT agent_id) as unique_agents
FROM gameinfo
WHERE status = 1
UNION ALL
SELECT 
    'lobbyinfo' as table_name,
    COUNT(*) as record_count,
    COUNT(DISTINCT CONCAT(agent_id, '-', game_id)) as unique_combinations
FROM lobbyinfo
WHERE status = 1;

-- 輸出初始化完成統計
SELECT '=== GameHub 數據庫初始化完成 ===' as status;
SELECT * FROM v_database_summary;

-- 輸出表計數統計
SELECT 
    'Tables created' as summary,
    COUNT(*) as count
FROM information_schema.tables 
WHERE table_schema = 'public' 
    AND table_type = 'BASE TABLE';

-- 輸出索引統計
SELECT 
    'Indexes created' as summary,
    COUNT(*) as count
FROM pg_indexes 
WHERE schemaname = 'public';

-- 設置數據庫完成標記
INSERT INTO game_settings (game_id, game_code, game_name, status, config) 
VALUES (0, 'SYSTEM', 'Database Initialization', 1, 
    json_build_object('initialized_at', NOW(), 'version', '1.0'))
ON CONFLICT DO NOTHING;

SELECT '✅ GameHub 數據庫初始化全部完成！' as final_status;