// v0.0.2

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/nathan-osman/go-sunrise"
	"github.com/reujab/wallpaper"
)

func get_season(month int, day int) string {
	if month > 11 || month < 3 {
		return "Winter"
	} else if month > 2 && month < 6 {
		return "Spring"
	} else if month > 5 && month < 9 {
		return "Summer"
	} else if month == 9 {
		if day < 22 {
			return "Summer"
		} else {
			return "Autumn"
		}
	} else if month > 9 && month < 12 {
		return "Autumn"
	} else {
		return "Winter"
	}
}

func get_day_night(time time.Time) string {

	destination, err := os.Readlink("/etc/localtime")
	if err != nil {
		fmt.Printf("/etc/localtime does not exist\n")
		os.Exit(1)
	}

	time_zone_full := strings.Split(destination, "/")
	time_zone_full = time_zone_full[len(time_zone_full)-2:]

	time_zone := strings.Join(time_zone_full, "/")

	db, err := os.Open("/usr/share/zoneinfo/zone.tab")
	if err != nil {
		fmt.Printf("/usr/share/zoneinfo/zone.tab does not exist")
		os.Exit(1)
	}

	var lat float64
	var long float64
	var coordinates [2]float64

	scanner := bufio.NewScanner(db)

	for scanner.Scan() {
		entry := scanner.Text()
		if !strings.HasPrefix(entry, "#") {
			info := strings.Split(entry, " ")[0]
			if strings.Contains(info, time_zone) {
				info = info[2:]
				info = strings.TrimLeftFunc(info, unicode.IsSpace)

				time_zone_check := info[len(info)-len(time_zone):]
				if time_zone_check != time_zone {
					for time_zone_check != time_zone {
						info = info[:len(info)-1]
						time_zone_check = info[len(info)-len(time_zone):]
					}
				}

				info = info[:len(info)-len(time_zone)]
				info = strings.TrimRightFunc(info, unicode.IsSpace)

				plus := strings.Index(info, "+")
				minus := strings.Index(info, "-")

				if minus == -1 { // coordinates are (+_,+_)
					split := strings.Split(info[1:], "+")
					lat, _ = strconv.ParseFloat(split[0], 64)
					long, _ = strconv.ParseFloat(split[1], 64)
				} else if plus == -1 { // coordinates are (-_, -_)
					split := strings.Split(info[1:], "-")
					lat, _ = strconv.ParseFloat(split[0], 64)
					lat = lat * -1
					long, _ = strconv.ParseFloat(split[1], 64)
					long = long * -1
				} else if plus == 1 { // coordinates are (-_, +_)
					split := strings.Split(info[1:], "+")
					lat, _ = strconv.ParseFloat(split[0], 64)
					lat = lat * -1
					long, _ = strconv.ParseFloat(split[1], 64)
				} else { // coordinates are (+_, -_)
					split := strings.Split(info[1:], "-")
					lat, _ = strconv.ParseFloat(split[0], 64)
					long, _ = strconv.ParseFloat(split[1], 64)
					long = long * -1
				}
				coordinates = [2]float64{lat / 100, long / 100}
				break
			}
		}
	}
	db.Close()

	sunrise, sunset := sunrise.SunriseSunset(
		coordinates[0], coordinates[1],
		time.Year(), time.Month(), time.Day(),
	)

	time = time.UTC()

	if time.Before(sunset) && time.After(sunrise) {
		return "Day"
	} else {
		return "Night"
	}
}

func main() {
	current_time := time.Now()
	tod := get_day_night(current_time)

	wallpaper_path, err := wallpaper.Get()

	current_file, err := os.ReadFile(wallpaper_path + "name.txt")
	if err != nil {
		fmt.Println("Active/name.txt not found")
		os.Exit(1)
	}

	current_folder := string(current_file)
	current_attributes := strings.Split(string(current_folder), "-")
	current_season := current_attributes[0]

	current_tod := current_attributes[1]

	temp_path := strings.Split(wallpaper_path, "/")
	wallpapers_folder := strings.Join(temp_path[:len(temp_path)-2], "/")

	if current_season != get_season(int(current_time.Month()), current_time.Day()) || current_tod != tod { //if active wallpaper does not match current season or tod
		entries, err := os.ReadDir(wallpapers_folder)
		if err != nil {
			fmt.Println("Couldn't load wallpaper folders")
		}

		for _, entry := range entries {
			if entry.IsDir() {
				file, _ := os.ReadFile(wallpapers_folder + "/" + entry.Name() + "/name.txt")
				folder := string(file)
				attributes := strings.Split(folder, "-")
				if attributes[0] == get_season(int(current_time.Month()), current_time.Day()) && attributes[1] == tod {
					_ = os.Rename(wallpapers_folder+"/Active", wallpapers_folder+"/"+current_season+"-"+current_tod)
					_ = os.Rename(wallpapers_folder+"/"+entry.Name(), wallpapers_folder+"/Active")
					break
				}
			}
		}
		ascript := `tell application "System Events" to tell current desktop to set change interval of current desktop to (change interval of current desktop)`

		cmd := exec.Command("osascript", "-e", ascript)
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to run Applescript command: %w", err)
			os.Exit(1)
		}
	}
}
