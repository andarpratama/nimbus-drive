#!/bin/bash

echo "Removing comments from project files..."

# Remove single-line comments from Vue files
find client/src -name "*.vue" -type f -exec sed -i 's|//.*$||g' {} \;

# Remove single-line comments from TypeScript files
find client/src -name "*.ts" -type f -exec sed -i 's|//.*$||g' {} \;

# Remove single-line comments from JavaScript files
find client/src -name "*.js" -type f -exec sed -i 's|//.*$||g' {} \;

# Remove single-line comments from Go files
find server -name "*.go" -type f -exec sed -i 's|//.*$||g' {} \;

# Remove empty lines that were left after comment removal
find client/src -name "*.vue" -type f -exec sed -i '/^[[:space:]]*$/d' {} \;
find client/src -name "*.ts" -type f -exec sed -i '/^[[:space:]]*$/d' {} \;
find client/src -name "*.js" -type f -exec sed -i '/^[[:space:]]*$/d' {} \;
find server -name "*.go" -type f -exec sed -i '/^[[:space:]]*$/d' {} \;

echo "Comments removed successfully!" 