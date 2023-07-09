// jime.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"time"
)

var using_list bool
var using_per bool
var jime time.Time

type Data struct {
	Clear_screen      bool
	Round_to_min      time.Duration
	Round_to_min_list []int
	Loop_sec          time.Duration
	Round_up_min      time.Duration
	Round_up_per      float64
}

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func init() {
	file, err := os.OpenFile("jime.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		ErrorLog.Fatal(err)
	}

	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func validate_config() (bool, time.Duration, time.Duration, []int, time.Duration, float64) {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		ErrorLog.Fatal("Error when opening file: ", err)
	}

	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		ErrorLog.Fatal("Error during Unmarshal(): ", err)
	}

	if payload.Round_to_min != 0 && len(payload.Round_to_min_list) != 0 {
		ErrorLog.Fatal("Invalid configuration: only one of round_to_min and round_to_min_list allowed", err)
	} else if payload.Round_to_min == 0 && payload.Round_to_min_list == nil {
		ErrorLog.Fatal("Invalid configuration: one of round_to_min or round_to_min list is required", err)
	} else {
		if payload.Round_to_min != 0 {
			using_list = false
		} else {
			using_list = true
		}
	}
	if payload.Round_up_min != 0 && payload.Round_up_per != 0 {
		ErrorLog.Fatal("Invalid configuration: only one of round_up_min and round_up_per allowed", err)
	} else if payload.Round_up_min == 0 && payload.Round_up_per == 0 {
		ErrorLog.Fatal("Invalid configuration: one of round_up_min or round_up_per is required", err)
	} else {
		if payload.Round_up_min != 0 {
			using_per = false
		} else {
			using_per = true
		}
	}
	clear_screen := payload.Clear_screen
	loop_sec := payload.Loop_sec
	round_to_min := payload.Round_to_min
	round_to_min_list := payload.Round_to_min_list
	round_up_min := payload.Round_up_min
	round_up_per := payload.Round_up_per
	return clear_screen, loop_sec, round_to_min, round_to_min_list, round_up_min, round_up_per
}

func main() {
	clear_screen, loop_sec, round_to_min, round_to_min_list, round_up_min, round_up_per := validate_config()
	t := time.Now()
	now_minute := t.Minute()

	if using_list {
		var low_rtm int
		var high_rtm int
		for i, v := range round_to_min_list {
			if i < len(round_to_min_list)-1 && v <= now_minute && now_minute <= round_to_min_list[i+1] {
				low_rtm = v
				high_rtm = round_to_min_list[i+1]
				break
			} else {
				low_rtm = v
				high_rtm = round_to_min_list[0] + 60
			}
		}
		round_to_min = time.Duration(high_rtm-low_rtm) * time.Minute
	}

	if using_per {
		round_up_min = time.Duration((round_up_per/100)*float64(round_to_min.Minutes())) * time.Minute
	} else {
		round_up_min = time.Duration(round_up_min.Minutes())
	}

	for {
		t = time.Now()
		InfoLog.Println("t is", t.Format("3:04:05 PM"))
		now_minute = t.Minute()
		InfoLog.Println("now_minute is", now_minute)
		minute_round_up := now_minute + int(round_up_min.Minutes())
		InfoLog.Println("minute_round_up is", minute_round_up)
		hour_round_up := 0
		minute_round_up_mod := math.Mod(float64(minute_round_up), float64(round_to_min.Minutes()))

		if minute_round_up > 59 {
			minute_round_up = minute_round_up - 60
			hour_round_up = 1
		}

		if clear_screen {
			cmd := exec.Command("clear") //works on Darwin
			cmd.Stdout = os.Stdout
			cmd.Run()
		}

		if minute_round_up_mod == 0 {
			jime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+hour_round_up, t.Round(time.Duration(round_to_min)).Minute(), t.Second(), t.Nanosecond(), t.Location())
		} else {
			jime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+hour_round_up, minute_round_up-int(minute_round_up_mod), t.Second(), t.Nanosecond(), t.Location())
		}
		fmt.Println("The jime is", jime.Format("3:04 PM"))
		InfoLog.Println("jime is", jime.Format("3:04 PM"))

		var round time.Duration = loop_sec * time.Second
		next_loop_time := t.Round(round)
		sleep := next_loop_time.Sub(t)
		if sleep < 0 {
			sleep = next_loop_time.Sub(t) + loop_sec*time.Second
		}
		time.Sleep(sleep)
	}
}
