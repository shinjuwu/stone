@echo off
:: Set the console to UTF-8 to handle any special characters in paths or outputs
chcp 65001 >nul

echo ====================================================
echo  GameHub Windows Development Environment Setup Script
echo ====================================================
echo.

:: Check for Administrator privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] Please run this script as an Administrator.
    goto :error_exit
)

:: Check for Docker Desktop
echo [CHECK] Checking for Docker Desktop...
docker --version >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] Docker Desktop is not installed or not running.
    echo Please install Docker Desktop first: https://www.docker.com/products/docker-desktop
    goto :error_exit
) else (
    echo [SUCCESS] Docker Desktop is installed.
)

:: Check for Docker Compose
echo [CHECK] Checking for Docker Compose...
docker-compose --version >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] Docker Compose not found.
    goto :error_exit
) else (
    echo [SUCCESS] Docker Compose is installed.
)

:: Check for VS Code
echo [CHECK] Checking for VS Code...
call code --version >nul 2>&1
if %errorLevel% neq 0 (
    echo [WARNING] VS Code not found in PATH.
    echo Please ensure VS Code is installed and added to the system PATH.
) else (
    echo [SUCCESS] VS Code is installed.
)

:: Check for WSL2
echo [CHECK] Checking for WSL2...
wsl --list --verbose >nul 2>&1
if %errorLevel% neq 0 (
    echo [WARNING] WSL2 is not installed or enabled.
    echo Installing WSL2 is recommended for better performance.
) else (
    echo [SUCCESS] WSL2 is installed.
)

:: Create necessary directories
echo [SETUP] Creating required directories...
if not exist "docker\nginx\ssl" mkdir "docker\nginx\ssl"
if not exist "data\postgres" mkdir "data\postgres"
if not exist "data\redis" mkdir "data\redis"
if not exist "logs\nginx" mkdir "logs\nginx"
if not exist "logs\gamehub" mkdir "logs\gamehub"

:: Copy environment variable file
echo [SETUP] Setting up environment variable file...
if not exist ".env" (
    if exist ".env.example" (
        copy ".env.example" ".env" >nul
        echo [SUCCESS] .env file created successfully from .env.example.
    ) else (
        echo [ERROR] .env.example not found, cannot create .env file.
        goto :error_exit
    )
) else (
    echo [SKIP] .env file already exists.
)

:: Configure Git settings
echo [SETUP] Configuring Git...
git config --global core.autocrlf input
git config --global core.eol lf
echo [SUCCESS] Git configuration complete.

:: Install VS Code extensions
echo [SETUP] Installing VS Code extensions...
set vscode_extensions_failed=0
for %%e in (golang.go ms-vscode-remote.remote-containers ms-azuretools.vscode-docker ms-vscode.vscode-json redhat.vscode-yaml) do (
    call code --install-extension %%e >nul 2>&1
    if %errorLevel% neq 0 (
        echo [WARNING] Failed to install VS Code extension %%e. Please check if VS Code is in your PATH.
        set vscode_extensions_failed=1
    )
)
if %vscode_extensions_failed% equ 0 (
    echo [SUCCESS] VS Code extensions installed successfully.
)

echo.
echo ====================================
echo  Setup Complete!
echo ====================================
echo.
echo Next steps:
echo 1. Start Docker Desktop if it's not already running.
echo 2. Open this project folder in VS Code.
echo 3. Press Ctrl+Shift+P and run "Remote-Containers: Reopen in Container".
echo 4. Alternatively, you can run the dev-start.bat script.
echo.
pause
exit /b 0

:: Error handling block
:error_exit
echo.
echo ====================================
echo  Script terminated due to an error.
echo ====================================
echo.
pause
exit /b 1