#!/bin/bash

# Build script for Ghost Host Manager
# This script creates a temporary manifest file for Windows builds and cleans up afterward

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
MANIFEST_DEV="$PROJECT_DIR/build/windows/wails.exe.dev.manifest"
MANIFEST_TEMP="$PROJECT_DIR/build/windows/wails.exe.manifest"

echo "========================================="
echo "Ghost Host Manager Build Script"
echo "========================================="
echo ""

# Check if the development manifest exists
if [ ! -f "$MANIFEST_DEV" ]; then
    echo "Error: $MANIFEST_DEV not found!"
    exit 1
fi

echo "Step 1: Copying manifest file..."
cp "$MANIFEST_DEV" "$MANIFEST_TEMP"
if [ $? -ne 0 ]; then
    echo "Error: Failed to copy manifest file!"
    exit 1
fi
echo "✓ Copied $MANIFEST_DEV to $MANIFEST_TEMP"
echo ""

# Build the application
echo "Step 2: Building application with wails build..."
cd "$PROJECT_DIR"
wails build
BUILD_STATUS=$?
echo ""

# Cleanup: Remove the temporary manifest file regardless of build status
echo "Step 3: Cleaning up temporary manifest file..."
rm -f "$MANIFEST_TEMP"
echo "✓ Removed $MANIFEST_TEMP"
echo ""

# Report build status
echo "========================================="
if [ $BUILD_STATUS -eq 0 ]; then
    echo "✓ Build completed successfully!"
else
    echo "✗ Build failed with exit code: $BUILD_STATUS"
fi
echo "========================================="

exit $BUILD_STATUS
