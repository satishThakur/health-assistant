# Development Quick Start

## Prerequisites

- Flutter SDK 3.0+
- GCC compiler (already installed)
- Backend API running on `localhost:8083`

## Quick Commands

### First Time Setup
```bash
make setup
# or
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
```

### Run the App

**Option 1: Quick Run (fastest)**
```bash
./run.sh
```

**Option 2: Development Mode (with hot reload)**
```bash
./dev.sh
# or
make debug
```

**Option 3: Using Makefile**
```bash
make run        # Build and run
make debug      # Run with hot reload
make release    # Build release version
```

## Makefile Commands

```bash
make help      # Show all commands
make setup     # First-time setup
make gen       # Generate .g.dart files
make clean     # Clean build
make build     # Build debug
make run       # Build and run
make debug     # Hot reload mode
make release   # Release build
make test      # Run tests
```

## Common Workflows

### Daily Development
```bash
# Start development with hot reload
./dev.sh

# Press 'r' for hot reload
# Press 'R' for hot restart
# Press 'q' to quit
```

### After Changing Models
```bash
# Regenerate .g.dart files
make gen
# or
flutter pub run build_runner build --delete-conflicting-outputs
```

### Clean Build
```bash
make clean
make run
```

### Release Build
```bash
make release
./build/linux/x64/release/bundle/health_assistant
```

## Troubleshooting

### Build Fails
```bash
# Clean everything and rebuild
make clean
make setup
make run
```

### Missing Generated Files
```bash
make gen
```

### Permission Errors
```bash
# Use the run.sh script (no system install)
./run.sh
```

## File Watching (Auto-regenerate)

Watch for model changes and auto-generate:
```bash
make watch
# or
flutter pub run build_runner watch --delete-conflicting-outputs
```

Keep this running in a separate terminal while developing.

## Backend Connection

By default, the app connects to `http://localhost:8083`.

To use a different API:
```bash
flutter run -d linux --dart-define=API_BASE_URL=http://192.168.1.100:8083
```

## IDE Setup

### VS Code
1. Install "Flutter" extension
2. Install "Dart" extension
3. Open Command Palette (Ctrl+Shift+P)
4. Select "Flutter: Select Device" → Linux (desktop)
5. Press F5 to run

### IntelliJ/Android Studio
1. Install Flutter plugin
2. Open project
3. Select "Linux" as device
4. Click Run button

## Tips

- **Hot reload (r)**: Fast, preserves state
- **Hot restart (R)**: Full restart, clears state
- **First run**: Takes 2-3 minutes (compiling)
- **Subsequent runs**: ~30 seconds
- **Release builds**: Much faster than debug

## Performance

- **Debug build**: Slower, full debugging support
- **Release build**: Optimized, production-ready
- **Profile build**: Performance profiling enabled

## Next Steps

1. Start backend: `cd ../infra && docker compose up -d`
2. Run app: `./run.sh`
3. Submit a check-in
4. View dashboard

Enjoy! 🚀
