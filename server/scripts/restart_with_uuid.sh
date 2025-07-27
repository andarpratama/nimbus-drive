#!/bin/bash

echo "🔄 Restarting application with UUID support..."

# Stop the current containers
echo "📦 Stopping current containers..."
docker compose down

# Reset the database (this will drop and recreate it)
echo "🗄️  Resetting database..."
docker compose exec -T mysql mysql -u root -ppassword -e "DROP DATABASE IF EXISTS nimbus_drive;"
docker compose exec -T mysql mysql -u root -ppassword -e "CREATE DATABASE nimbus_drive CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Start the application
echo "🚀 Starting application with UUID support..."
docker compose up -d

echo "✅ Application restarted with UUID support!"
echo "📝 Note: All existing data has been cleared. New users and files will use UUIDs." 