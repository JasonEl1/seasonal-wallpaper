# Time and season - aware wallpaper tool

Automatically change your wallpaper folder to match the season and time of day!

### Instructions:

> [!IMPORTANT]
> Only tested on MacOS. Parts of the application will currently not function correctly outside of MacOS. The tool is still usable on other platforms however there will be some delay between running the tool and the wallpaper correctly updating.

Create a master folder to store all wallpaper images and folders. This can have any name.

Within the master folder, create the following sub-folders or run the below shell script inside the master folder:

/Spring-Day  
/Spring-Night  
/Summer-Day  
/Summer-Night  
/Autumn-Day  
/Autumn-Night  
/Winter-Day  
/Winter-Night  

Important: Each sub-folder should contain a file called name.txt which contains the name of the folder. This is because the current wallpaper folder will be renamed to `Active` by the tool. 


```
folders=("Spring-Day" "Spring-Night" "Summer-Day" "Summer-Night" "Autumn-Day" "Autumn-Night" "Winter-Day" "Winter-Night")

for folder in "${folders[@]}"; do
    mkdir "$folder"
    cd "$folder"
    echo "$folder" >> name.txt
    cd ..
done
```

Once these folders and files have been created, rename one of the folders to `Active` and set it as your wallpaper folder in your system's settings.

Go back to the initial folder from github and run the `./install.sh` script. This will install the compiled binary to `/usr/local/bin` and add a cron job that run the program once an hour and update your wallpaper folder if needed.

You can also enter `wallpaper` in the terminal to update the wallpaper folder manually.

### Uninstallation

To uninstall the program, run `./uninstall.sh` from the original install folder. Or enter this in the terminal:

```
rm /usr/local/bin/wallpaper
```
