#!/bin/bash

# Development mode with hot reload
# Usage: ./dev.sh

set -e

echo "🔥 Starting development mode with hot reload..."
echo "================================================"
echo ""
echo "Commands:"
echo "  r  - Hot reload (fast)"
echo "  R  - Hot restart (full restart)"
echo "  q  - Quit"
echo ""

# Set compiler
export CC=gcc
export CXX=g++

# Check if generated files exist
if [ ! -f "lib/features/checkin/domain/checkin_model.g.dart" ]; then
    echo "⚠️  Generated files missing. Running code generation..."
    flutter pub get
    flutter pub run build_runner build --delete-conflicting-outputs
fi

# Run with hot reload
flutter run -d linux
