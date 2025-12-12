#!/bin/bash

wallpaper_path=$(osascript -e 'tell app "Finder" to get posix path of (get desktop picture as alias)')
enclosing_folder=$(dirname "$wallpaper_path")
folder_count=$(find "$enclosing_folder" -maxdepth 1 -mindepth 1 -type d | wc -l | tr -d '[:space:]')

if [[ $"folder_count" -gt 7 ]]; then
    sys_type=$(uname)
    sys_type=$(echo "$sys_type" | tr '[:upper:]' '[:lower:]')

    if [[ $sys_type == "darwin" ]]; then
        sys_type="mac"
    fi

    sys_arch=$(uname -m)

    echo "Detected system as $sys_type-$sys_arch. Continue? [y/n]"

    read continue

    if [[ $continue == "y" ]]; then
        sudo cp build/wallpaper-${sys_type}-${sys_arch} /usr/local/bin/wallpaper
        crontab -l | grep -v "wallpaper" | crontab -
        (crontab -l; echo "0 * * * * wallpaper") | crontab -
        echo "Copied executable to wallpaper folder and added cron entry."
    else
        echo "Received $continue. Exiting installer."
    fi
else
    echo "Destination directory structure is incorrect or wallpaper folder is set incorrectly. See README.md."
fi
