@echo off
title Reset GameHub Database

echo ====================================
echo Reset GameHub Database
echo ====================================
echo.

echo [Step 1] Stopping all containers...
docker-compose -f docker-compose.dev.yml down

echo.
echo [Step 2] Removing PostgreSQL data volume...
docker volume ls | findstr postgres_dev_data
if %errorLevel% equ 0 (
    docker volume rm stone_postgres_dev_data
    echo PostgreSQL data volume removed
) else (
    echo PostgreSQL data volume not found or already removed
)

echo.
echo [Step 3] Starting PostgreSQL container (will re-run init scripts)...
docker-compose -f docker-compose.dev.yml up -d postgres

echo.
echo [Step 4] Waiting for database initialization (about 2-3 minutes)...
echo Waiting for database startup...
timeout /t 10 /nobreak >nul

echo Running initialization scripts...
timeout /t 30 /nobreak >nul

echo Waiting for initialization complete...
timeout /t 30 /nobreak >nul

echo.
echo [Step 5] Checking initialization logs...
docker-compose -f docker-compose.dev.yml logs postgres | findstr -i "initialization"

echo.
echo [Step 6] Verifying tables and data...
echo Checking tables...
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U gamehub_dev -d gamehub_dev -c "\dt" 2>nul

echo.
echo Checking record counts...
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U gamehub_dev -d gamehub_dev -c "SELECT 'gamelist' as table_name, COUNT(*) as count FROM gamelist UNION ALL SELECT 'gameinfo' as table_name, COUNT(*) as count FROM gameinfo UNION ALL SELECT 'lobbyinfo' as table_name, COUNT(*) as count FROM lobbyinfo;" 2>nul

echo.
echo [Step 7] Starting other services...
docker-compose -f docker-compose.dev.yml up -d redis pgadmin redis-commander

echo.
echo ====================================
echo Database reset and initialization complete!
echo ====================================
echo.
echo Expected results:
echo - gamelist: about 30+ records
echo - gameinfo: about 1000+ records  
echo - lobbyinfo: about 4000+ records
echo.
echo You can now run "Launch Full Stack" in VS Code
echo or visit pgAdmin: http://localhost:5050
echo.
pause