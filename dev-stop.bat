@echo off
chcp 65001 >nul
title GameHub 開發環境停止

echo ====================================
echo GameHub 開發環境停止
echo ====================================
echo.

echo [停止] 停止開發環境...
docker-compose -f docker-compose.dev.yml down

if %errorLevel% equ 0 (
    echo [成功] 開發環境已停止
) else (
    echo [錯誤] 停止失敗
)

echo.
pause