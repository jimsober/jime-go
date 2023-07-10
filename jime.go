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

var (
	hm_format     string
	hms_format    string
	using_list    bool
	using_percent bool
	jime          time.Time
	WarningLog    *log.Logger
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
)

type Data struct {
	Clear_screen          bool
	Military_display      bool
	Round_to_minutes      time.Duration
	Round_to_minutes_list []int
	Loop_seconds          time.Duration
	Round_up_minutes      time.Duration
	Round_up_percent      float64
}

func init() {
	file, err := os.OpenFile("jime.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		ErrorLog.Fatal(err)
	}

	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func validate_config() (bool, bool, time.Duration, time.Duration, []int, time.Duration, float64) {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		ErrorLog.Fatal("Error when opening file: ", err)
	}
	//InfoLog.Println("Configuration file opened successfully")

	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		ErrorLog.Fatal("Error during Unmarshal(): ", err)
	}
	//InfoLog.Println("Configuration data unmarshaled successfully")

	if payload.Round_to_minutes != 0 && len(payload.Round_to_minutes_list) != 0 {
		ErrorLog.Fatal("Invalid configuration: only one of round_to_minutes and round_to_minutes_list allowed", err)
	} else if payload.Round_to_minutes == 0 && payload.Round_to_minutes_list == nil {
		ErrorLog.Fatal("Invalid configuration: one of round_to_minutes or round_to_minutes_list is required", err)
	} else {
		if payload.Round_to_minutes != 0 {
			using_list = false
		} else {
			using_list = true
		}
	}
	//InfoLog.Println("payload.Round_to_minutes is", payload.Round_to_minutes)
	//InfoLog.Println("payload.Round_to_minutes_list is", payload.Round_to_minutes_list)
	//InfoLog.Println("using_list is", using_list)

	if payload.Round_up_minutes != 0 && payload.Round_up_percent != 0 {
		ErrorLog.Fatal("Invalid configuration: only one of round_up_minutes and round_up_percent allowed", err)
	} else if payload.Round_up_minutes == 0 && payload.Round_up_percent == 0 {
		ErrorLog.Fatal("Invalid configuration: one of round_up_minutes or round_up_percent is required", err)
	} else {
		if payload.Round_up_minutes != 0 {
			using_percent = false
		} else {
			using_percent = true
		}
	}
	//InfoLog.Println("payload.Round_up_minutes is", payload.Round_up_minutes)
	//InfoLog.Println("payload.Round_up_percent is", payload.Round_up_percent)
	//InfoLog.Println("using_percent is", using_percent)

	clear_screen := payload.Clear_screen
	military_display := payload.Military_display
	loop_seconds := payload.Loop_seconds
	round_to_minutes := payload.Round_to_minutes
	round_to_minutes_list := payload.Round_to_minutes_list
	round_up_minutes := payload.Round_up_minutes
	round_up_percent := payload.Round_up_percent
	return clear_screen, military_display, loop_seconds, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent
}

func run_jime(t time.Time, clear_screen bool, round_to_minutes time.Duration, round_to_minutes_list []int, round_up_minutes time.Duration, round_up_percent float64) {
	now_minute := t.Minute()
	//InfoLog.Println("now_minute is", now_minute)

	if using_list {
		var low_round_to_minute int
		var high_round_to_minute int
		for i, v := range round_to_minutes_list {
			if i < len(round_to_minutes_list)-1 && v <= now_minute && now_minute <= round_to_minutes_list[i+1] {
				low_round_to_minute = v
				high_round_to_minute = round_to_minutes_list[i+1]
				break
			} else {
				low_round_to_minute = v
				high_round_to_minute = round_to_minutes_list[0] + 60
			}
		}
		//InfoLog.Println("low_round_to_minute is", low_round_to_minute)
		//InfoLog.Println("high_round_to_minute is", high_round_to_minute)
		round_to_minutes = time.Duration(high_round_to_minute-low_round_to_minute) * time.Minute
	} else {
		round_to_minutes = time.Duration(round_to_minutes) * time.Minute
	}
	InfoLog.Println("round_to_minutes is", round_to_minutes)

	if using_percent {
		round_up_minutes = time.Duration((round_up_percent/100)*float64(round_to_minutes.Minutes())) * time.Minute
	} else {
		round_up_minutes = time.Duration(round_up_minutes) * time.Minute
	}
	InfoLog.Println("round_up_minutes is", round_up_minutes)

	round_up_minute := now_minute + int(round_up_minutes.Minutes())
	round_up_hour := 0

	if round_up_minute > 59 {
		round_up_minute = round_up_minute - 60
		round_up_hour = 1
	}
	//InfoLog.Println("round_up_hour is", round_up_hour)
	//InfoLog.Println("round_up_minute is", round_up_minute)
	round_up_minute_mod := math.Mod(float64(round_up_minute), float64(round_to_minutes.Minutes()))
	//InfoLog.Println("round_up_minute_mod is", round_up_minute_mod)

	if clear_screen {
		cmd := exec.Command("clear") //works on Darwin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	if round_up_minute_mod == 0 {
		jime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+round_up_hour, t.Round(time.Duration(round_to_minutes)).Minute(), t.Second(), t.Nanosecond(), t.Location())
	} else {
		jime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+round_up_hour, round_up_minute-int(round_up_minute_mod), t.Second(), t.Nanosecond(), t.Location())
	}
	fmt.Println("The jime is", jime.Format(hm_format))
	InfoLog.Println("jime is", jime.Format(hm_format))
}

func main() {
	clear_screen, military_display, loop_seconds, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent := validate_config()
	//InfoLog.Println("** configuration: clear_screen is", clear_screen, "; military_display is", military_display, "; loop_seconds is", loop_seconds, "; round_to_minutes is", round_to_minutes, "; round_to_minutes_list is", round_to_minutes_list, "; round_up_minutes is", round_up_minutes, "; round_up_percent is", round_up_percent)
	if military_display {
		hm_format = "15:04"
		hms_format = "15:04:05"
	} else {
		hm_format = "3:04 PM"
		hms_format = "3:04:05 PM"
	}
	//InfoLog.Println("hm_format is", hm_format)
	//InfoLog.Println("hms_format is", hms_format)

	t := time.Now()
	InfoLog.Println("* time is", t.Format(hms_format))

	run_jime(t, clear_screen, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent)

	if loop_seconds.Seconds() != 0 {
		for {
			next_loop_time := t.Round(loop_seconds * time.Second)
			sleep := next_loop_time.Sub(t)

			if sleep < 0 {
				sleep = next_loop_time.Sub(t) + loop_seconds*time.Second
			}
			//InfoLog.Println("sleep is", sleep)

			time.Sleep(sleep)

			t = time.Now()
			InfoLog.Println("* time is", t.Format(hms_format))

			run_jime(t, clear_screen, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent)
		}
	}
}
