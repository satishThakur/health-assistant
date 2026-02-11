.PHONY: help setup clean build run debug release install gen test

# Default target
help:
	@echo "Flutter App - Development Commands"
	@echo "=================================="
	@echo "make setup    - Install dependencies and generate code"
	@echo "make gen      - Generate .g.dart files"
	@echo "make clean    - Clean build artifacts"
	@echo "make build    - Build debug version"
	@echo "make run      - Build and run app"
	@echo "make debug    - Run in debug mode with hot reload"
	@echo "make release  - Build release version"
	@echo "make test     - Run tests"
	@echo "make install  - Install dependencies"

# First-time setup
setup:
	@echo "📦 Installing dependencies..."
	flutter pub get
	@echo "🔨 Generating code..."
	flutter pub run build_runner build --delete-conflicting-outputs
	@echo "✅ Setup complete!"

# Generate code
gen:
	@echo "🔨 Generating .g.dart files..."
	flutter pub run build_runner build --delete-conflicting-outputs

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	flutter clean
	rm -rf build/

# Build debug version
build:
	@echo "🔨 Building debug version..."
	export CC=gcc && export CXX=g++ && flutter build linux --debug

# Quick run (build + execute)
run: build
	@echo "🚀 Running app..."
	./build/linux/x64/debug/bundle/health_assistant

# Run with hot reload
debug:
	@echo "🔥 Running in debug mode with hot reload..."
	export CC=gcc && export CXX=g++ && flutter run -d linux

# Build release version
release:
	@echo "🚀 Building release version..."
	export CC=gcc && export CXX=g++ && flutter build linux --release
	@echo "✅ Release build complete!"
	@echo "📦 Binary: ./build/linux/x64/release/bundle/health_assistant"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	flutter pub get

# Run tests
test:
	@echo "🧪 Running tests..."
	flutter test

# Watch mode for code generation
watch:
	@echo "👀 Watching for changes..."
	flutter pub run build_runner watch --delete-conflicting-outputs
