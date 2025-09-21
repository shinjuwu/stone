@echo off
title Check GameHub Database Status

echo ====================================
echo Check GameHub Database Status
echo ====================================
echo.

echo [Check 1] Container running status...
docker ps --filter "name=gamehub" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo.
echo [Check 2] PostgreSQL connection test...
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U gamehub_dev -d gamehub_dev -c "SELECT version();" 2>nul
if %errorLevel% equ 0 (
    echo Database connection: SUCCESS
) else (
    echo Database connection: FAILED
)

echo.
echo [Check 3] Table existence check...
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U gamehub_dev -d gamehub_dev -c "\dt" 2>nul

echo.
echo [Check 4] Key table record counts...
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U gamehub_dev -d gamehub_dev -c "SELECT 'gamelist' as table_name, COUNT(*) as count FROM gamelist WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gamelist') UNION ALL SELECT 'gameinfo' as table_name, COUNT(*) as count FROM gameinfo WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gameinfo') UNION ALL SELECT 'lobbyinfo' as table_name, COUNT(*) as count FROM lobbyinfo WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'lobbyinfo');" 2>nul

echo.
echo [Check 5] Initialization status check...
docker-compose -f docker-compose.dev.yml logs postgres | findstr -i "initialization_status"

echo.
echo ====================================
echo Status check complete
echo ====================================
echo.
echo If you see tables and record counts above, initialization was successful
echo If not, please run: reset-and-init-db.bat
echo.
pause