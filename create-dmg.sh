#!/bin/bash

set -e

APP_NAME="SilverBullet"
DMG_NAME="$APP_NAME.dmg"
VOLUME_NAME="$APP_NAME Installer"
APP_BUNDLE="$APP_NAME.app"
OUTPUT_DIR="dmg_output"

# Ensure required tools are installed
if ! command -v create-dmg &> /dev/null; then
    echo "[Error] create-dmg not found. Please install it first."
    exit 1
fi

# Create the application bundle structure
mkdir -p $OUTPUT_DIR/$APP_BUNDLE/Contents/MacOS
mkdir -p $OUTPUT_DIR/$APP_BUNDLE/Contents/Resources

# Copy SPV resolver binary
cp spvproc $OUTPUT_DIR/$APP_BUNDLE/Contents/MacOS/

# Create Info.plist
cat > $OUTPUT_DIR/$APP_BUNDLE/Contents/Info.plist <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>spvproc</string>
    <key>CFBundleIdentifier</key>
    <string>com.silverbullet.spv</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
</dict>
</plist>
EOL

# Generate the DMG
create-dmg   --volname "$VOLUME_NAME"   --window-pos 200 120   --window-size 800 400   --icon-size 100   --app-drop-link 600 185   $OUTPUT_DIR/$DMG_NAME   $OUTPUT_DIR/

echo "[Success] $DMG_NAME created successfully."
