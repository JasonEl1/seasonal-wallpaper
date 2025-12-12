#!/bin/bash

# Target combinations: OS/ARCH/OutputSuffix
targets=(
    "darwin/arm64/mac-arm64"
    "darwin/amd64/mac-x86"
    "linux/amd64/linux-x86_64"
    "windows/amd64/windows-x64.exe"
)

for target in "${targets[@]}"; do
    IFS='/' read -r GOOS GOARCH OUTPUT_SUFFIX <<< "$target"
    OUTPUT_FILE="wallpaper-${OUTPUT_SUFFIX}"

    echo "Building $GOOS/$GOARCH ($OUTPUT_FILE)"

    CGO_ENABLED=0
    GOOS=$GOOS
    GOARCH=$GOARCH
    go build -o "build/$OUTPUT_FILE" .

    if [ $? -eq 0 ]; then
        echo "     Success: build/$OUTPUT_FILE"
    else
        echo "     FAILURE for $GOOS/$GOARCH"
    fi
done

echo "Done building"
