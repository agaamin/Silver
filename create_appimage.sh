#!/bin/bash

set -e

APP_NAME="SilverBullet"
APP_DIR="appdir"
OUTPUT_FILE="$APP_NAME.AppImage"

# Ensure dependencies are installed
if ! command -v appimagetool &> /dev/null; then
    echo "[Error] appimagetool not found. Please install it first."
    exit 1
fi

# Create AppDir structure
mkdir -p $APP_DIR/usr/bin
mkdir -p $APP_DIR/usr/share/applications
mkdir -p $APP_DIR/usr/share/icons/hicolor/256x256/apps

# Copy SPV resolver binary
cp spvproc $APP_DIR/usr/bin/

# Create .desktop file
cat > $APP_DIR/usr/share/applications/$APP_NAME.desktop <<EOL
[Desktop Entry]
Name=$APP_NAME
Exec=/usr/bin/spvproc
Icon=$APP_NAME
Type=Application
Categories=Network;
EOL

# Copy icon
cp assets/icon.png $APP_DIR/usr/share/icons/hicolor/256x256/apps/$APP_NAME.png

# Generate AppImage
appimagetool $APP_DIR $OUTPUT_FILE

echo "[Success] $OUTPUT_FILE created successfully."
