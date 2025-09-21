@echo off
chcp 65001 >nul
title GameHub 開發環境啟動

echo ====================================
echo GameHub 開發環境啟動
echo ====================================
echo.

:: 檢查 Docker 是否運行
echo [檢查] 檢查 Docker 狀態...
docker info >nul 2>&1
if %errorLevel% neq 0 (
    echo [錯誤] Docker 未運行，請先啟動 Docker Desktop
    pause
    exit /b 1
)

:: 檢查環境變量文件
if not exist ".env" (
    echo [錯誤] .env 文件不存在
    echo 請先運行: scripts\windows-setup.bat
    pause
    exit /b 1
)

:: 停止可能正在運行的容器
echo [清理] 停止現有容器...
docker-compose -f docker-compose.dev.yml down >nul 2>&1

:: 啟動開發環境
echo [啟動] 啟動開發環境...
docker-compose -f docker-compose.dev.yml up -d

if %errorLevel% equ 0 (
    echo [成功] 開發環境已啟動
    echo.
    echo 服務訪問地址：
    echo - 遊戲服務: http://localhost
    echo - pgAdmin: http://localhost:5050
    echo - Redis Commander: http://localhost:8081
    echo - API 文檔: http://localhost:8083
    echo.
    echo 在 VS Code 中：
    echo 1. 按 Ctrl+Shift+P
    echo 2. 輸入並選擇 "Remote-Containers: Attach to Running Container"
    echo 3. 選擇 "gamehub-server-dev"
    echo.
    
    :: 詢問是否打開 VS Code
    set /p open_vscode="是否打開 VS Code？(y/n): "
    if /i "%open_vscode%"=="y" (
        echo [啟動] 打開 VS Code...
        code .
    )
) else (
    echo [錯誤] 啟動失敗，請檢查 Docker 日誌
    docker-compose -f docker-compose.dev.yml logs
)

pause