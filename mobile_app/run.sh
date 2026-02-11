#!/bin/bash

# Quick run script for Flutter app
# Usage: ./run.sh

set -e

echo "🚀 Health Assistant - Quick Start"
echo "=================================="

# Set compiler
export CC=gcc
export CXX=g++

# Check if dependencies are installed
if [ ! -d ".dart_tool" ]; then
    echo "📦 First run detected - installing dependencies..."
    flutter pub get
    flutter pub run build_runner build --delete-conflicting-outputs
fi

# Clean and build
echo "🔨 Building app..."
flutter clean > /dev/null 2>&1
flutter build linux --debug

# Run
echo "✅ Build complete! Starting app..."
echo ""
./build/linux/x64/debug/bundle/health_assistant
