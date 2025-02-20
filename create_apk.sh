#!/bin/bash

set -e

APP_NAME="SilverBullet"
OUTPUT_APK="$APP_NAME.apk"
BUILD_DIR="android_build"

# Ensure required tools are installed
if ! command -v gradle &> /dev/null; then
    echo "[Error] Gradle is not installed. Please install it first."
    exit 1
fi

# Create project structure
mkdir -p $BUILD_DIR/app/src/main/java/com/silverbullet
mkdir -p $BUILD_DIR/app/src/main/res
mkdir -p $BUILD_DIR/app/src/main/assets

# Generate AndroidManifest.xml
cat > $BUILD_DIR/app/src/main/AndroidManifest.xml <<EOL
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="com.silverbullet.spv">
    <application android:label="$APP_NAME">
        <activity android:name=".MainActivity">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>
</manifest>
EOL

# Create basic Gradle build file
cat > $BUILD_DIR/app/build.gradle <<EOL
apply plugin: 'com.android.application'
android {
    compileSdkVersion 30
    defaultConfig {
        applicationId "com.silverbullet.spv"
        minSdkVersion 21
        targetSdkVersion 30
        versionCode 1
        versionName "1.0"
    }
    buildTypes {
        release {
            minifyEnabled false
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
}
EOL

# Build APK
cd $BUILD_DIR && gradle assembleDebug

mv app/build/outputs/apk/debug/app-debug.apk ../$OUTPUT_APK

cd ..
echo "[Success] $OUTPUT_APK created successfully."
