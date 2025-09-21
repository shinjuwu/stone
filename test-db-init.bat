@echo off
chcp 65001 >nul
title GameHub 數據庫初始化測試

echo ====================================
echo GameHub 數據庫初始化測試
echo ====================================
echo.

echo [步驟1] 停止現有容器並清理數據...
docker-compose -f docker-compose.dev.yml down -v

echo.
echo [步驟2] 啟動 PostgreSQL 容器 (會自動執行初始化腳本)...
docker-compose -f docker-compose.dev.yml up -d postgres

echo.
echo [步驟3] 等待數據庫啟動完成...
timeout /t 30 /nobreak

echo.
echo [步驟4] 檢查初始化日誌...
docker-compose -f docker-compose.dev.yml logs postgres | findstr "initialization_status"

echo.
echo [步驟5] 驗證數據庫表和數據...
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U gamehub_dev -d gamehub_dev -c "\dt"

echo.
echo [步驟6] 檢查遊戲數據統計...
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U gamehub_dev -d gamehub_dev -c "SELECT * FROM v_database_summary;"

echo.
echo ====================================
echo 測試完成！
echo ====================================
echo.
echo 如果看到上方有數據輸出，表示初始化成功
echo 現在可以啟動完整的開發環境：
echo   docker-compose -f docker-compose.dev.yml up -d
echo.
pause