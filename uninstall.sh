#!/bin/bash

wallpaper_path=$(osascript -e 'tell app "Finder" to get posix path of (get desktop picture as alias)')
rm $wallpaper_path/../seasons.json

crontab -l | grep -v "wallpaper" | crontab -

rm /usr/local/bin/wallpaper