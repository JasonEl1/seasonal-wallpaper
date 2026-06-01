// v0.1.4

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/nathan-osman/go-sunrise"
	"github.com/reujab/wallpaper"
)

var wallpapers_folder string

type Season struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func get_season(month int, day int) string {
	file, err := os.Open(wallpapers_folder + "/seasons.json")
	if err != nil {
		fmt.Println("Could not load seasons.json")
		os.Exit(1)
		return "-1"
	}
	defer file.Close()

	var season_data map[string]Season
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&season_data); err != nil {
		return ""
	}

	current_time, _ := time.Parse("1/2", strconv.Itoa(month)+"/"+strconv.Itoa(day))

	earliest_end, _ := time.Parse("1/2", "12/31")
	var earliest_season string

	for season, dates := range season_data {
		season_end, _ := time.Parse("1/2", dates.End)
		season_start, _ := time.Parse("1/2", dates.Start)

		if season_end.Before(earliest_end) {
			earliest_end = season_end
			earliest_season = season
		}

		if (current_time.Before(season_end) || current_time.Equal(season_end)) && (current_time.After(season_start) || current_time.Equal(season_start)) {
			return season
		}
	}
	return earliest_season
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
	OS_TYPE := runtime.GOOS

	current_time := time.Now()
	tod := get_day_night(current_time)

	wallpaper_path, _ := wallpaper.Get()

	var err error
	var current_file []byte

	for err != nil || current_file == nil {

		current_file, err = os.ReadFile(wallpaper_path + "name.txt")
		if err != nil {
			fmt.Println(wallpaper_path)
			fmt.Println("Active/name.txt not found. Recovering...")
			season := get_season(int(current_time.Month()), current_time.Day())
			os.Rename(wallpaper_path+"/"+season+"-"+tod, wallpaper_path+"/Active")
		}
	}

	current_folder := string(current_file)
	current_attributes := strings.Split(string(current_folder), "-")
	current_season := current_attributes[0]

	current_tod := current_attributes[1]

	temp_path := strings.Split(wallpaper_path, "/")
	wallpapers_folder = strings.Join(temp_path[:len(temp_path)-2], "/")

	var changed bool = false

	if (current_season != get_season(int(current_time.Month()), current_time.Day())) || (current_tod != tod) { //if active wallpaper does not match current season or tod
		fmt.Println("Found mismatch: current season is " + get_season(int(current_time.Month()), current_time.Day()) + " and current tod is " + tod + " but Active folder is " + current_season + "-" + current_tod)

		entries, err := os.ReadDir(wallpapers_folder)
		if err != nil {
			fmt.Println("Couldn't load wallpaper folders")
			os.Exit(1)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				file, err := os.ReadFile(wallpapers_folder + "/" + entry.Name() + "/name.txt")
				if err != nil {
					continue
				}
				folder := string(file)
				attributes := strings.Split(folder, "-")
				if attributes[0] == get_season(int(current_time.Month()), current_time.Day()) && attributes[1] == tod {
					err = os.Rename(wallpapers_folder+"/Active", wallpapers_folder+"/"+current_season+"-"+current_tod)
					if err != nil {
						fmt.Println("Could not rename Active folder")
						os.Exit(1)
					}
					err = os.Rename(wallpapers_folder+"/"+entry.Name(), wallpapers_folder+"/Active")
					if err != nil {
						fmt.Println("Could not set Active folder")
						os.Exit(1)
					}
					changed = true
					fmt.Println("Changed active folder" + " from " + current_season + "-" + current_tod + " to " + entry.Name())
					break
				}
			}
		}
	}

	if !changed {
		fmt.Println("did not change active folder")
	}

	if changed && OS_TYPE == "darwin" {
		cmd := exec.Command("killall", "WallpaperAgent")

		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to refresh wallpaper")
			os.Exit(1)
		}
	}

}
