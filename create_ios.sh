#!/bin/bash

set -e

APP_NAME="SilverBullet"
OUTPUT_IPA="$APP_NAME.ipa"
BUILD_DIR="ios_build"

# Ensure required tools are installed
if ! command -v xcodebuild &> /dev/null; then
    echo "[Error] Xcode and command-line tools are not installed. Please install them first."
    exit 1
fi

# Create project structure
mkdir -p $BUILD_DIR/$APP_NAME
cd $BUILD_DIR
xcodebuild -create-xcframework -framework $APP_NAME.framework -output $APP_NAME.xcframework

# Build IPA
xcodebuild -scheme $APP_NAME -archivePath $APP_NAME.xcarchive archive
xcodebuild -exportArchive -archivePath $APP_NAME.xcarchive -exportOptionsPlist exportOptions.plist -exportPath .

mv $APP_NAME.ipa ../$OUTPUT_IPA

cd ..
echo "[Success] $OUTPUT_IPA created successfully."
