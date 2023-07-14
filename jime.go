// jime.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	hm_format         string
	hms_format        string
	round_to_duration time.Duration
	loop_duration     time.Duration
	round_up_duration time.Duration
	using_list        bool
	using_percent     bool
	jime              time.Time
	DebugLog          *log.Logger
	InfoLog           *log.Logger
	WarningLog        *log.Logger
	ErrorLog          *log.Logger
)

type Data struct {
	Clear_screen          bool
	Military_display      bool
	Log_level             string
	Round_to_minutes      float64
	Round_to_minutes_list []float64
	Loop_seconds          float64
	Round_up_minutes      float64
	Round_up_percent      float64
}

func init() {
	file, err := os.OpenFile("jime.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		ErrorLog.Fatal(err)
	}

	DebugLog = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func validateConfig(default_log_level string) (bool, bool, string, float64, float64, []float64, float64, float64) {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		if default_log_level == "Critical" || default_log_level == "Warning" || default_log_level == "Info" || default_log_level == "Debug" {
			ErrorLog.Fatal("Error when opening file: ", err)
		}
	} else {
		if default_log_level == "Debug" {
			DebugLog.Println("Configuration file opened successfully")
		}
	}

	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		if default_log_level == "Critical" || default_log_level == "Warning" || default_log_level == "Info" || default_log_level == "Debug" {
			ErrorLog.Fatal("Error during Unmarshal(): ", err)
		}
	} else {
		if default_log_level == "Debug" {
			DebugLog.Println("Configuration data unmarshaled successfully")
		}
	}

	log_level_values := []string{"Debug", "Info", "Warning", "Critical"}
	if !isElementExist(log_level_values, payload.Log_level) {
		if default_log_level == "Critical" || default_log_level == "Warning" || default_log_level == "Info" || default_log_level == "Debug" {
			ErrorLog.Fatal("Invalid configuration: log_level must be one of 'Debug', 'Info', 'Warning', or 'Critical'", err)
		}
	}
	log_level := payload.Log_level

	if payload.Round_to_minutes != 0 && len(payload.Round_to_minutes_list) != 0 {
		if log_level == "Critical" || log_level == "Warning" || log_level == "Info" || log_level == "Debug" {
			ErrorLog.Fatal("Invalid configuration: both of round_to_minutes and round_to_minutes_list not allowed", err)
		}
	} else if payload.Round_to_minutes == 0 && payload.Round_to_minutes_list == nil {
		using_list = false
	} else {
		if payload.Round_to_minutes != 0 {
			using_list = false
		} else {
			using_list = true
		}
	}
	if log_level == "Debug" {
		DebugLog.Println("payload.Round_to_minutes is", payload.Round_to_minutes)
	}
	if log_level == "Debug" {
		DebugLog.Println("payload.Round_to_minutes_list is", payload.Round_to_minutes_list)
	}
	if log_level == "Debug" {
		DebugLog.Println("using_list is", using_list)
	}

	if payload.Round_up_minutes != 0 && payload.Round_up_percent != 0 {
		if log_level == "Critical" || log_level == "Warning" || log_level == "Info" || log_level == "Debug" {
			ErrorLog.Fatal("Invalid configuration: both of round_up_minutes and round_up_percent not allowed", err)
		}
	} else if payload.Round_up_minutes == 0 && payload.Round_up_percent == 0 {
		using_percent = false
	} else {
		if payload.Round_up_minutes != 0 {
			using_percent = false
		} else {
			using_percent = true
		}
	}
	if log_level == "Debug" {
		DebugLog.Println("payload.Round_up_minutes is", payload.Round_up_minutes)
	}
	if log_level == "Debug" {
		DebugLog.Println("payload.Round_up_percent is", payload.Round_up_percent)
	}
	if log_level == "Debug" {
		DebugLog.Println("using_percent is", using_percent)
	}

	clear_screen := payload.Clear_screen
	military_display := payload.Military_display
	loop_seconds := payload.Loop_seconds
	round_to_minutes := payload.Round_to_minutes
	round_to_minutes_list := payload.Round_to_minutes_list
	round_up_minutes := payload.Round_up_minutes
	round_up_percent := payload.Round_up_percent

	return clear_screen, military_display, log_level, loop_seconds, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent
}

func calculateAndDisplayJime(t time.Time, clear_screen bool, log_level string, round_to_minutes float64, round_to_minutes_list []float64, round_up_minutes float64, round_up_percent float64) {
	now_minute := t.Minute()
	if log_level == "Debug" {
		DebugLog.Println("now_minute is", now_minute)
	}

	if using_list {
		var low_round_to_minute float64
		var high_round_to_minute float64
		for i, v := range round_to_minutes_list {
			if i < len(round_to_minutes_list)-1 && v <= float64(now_minute) && float64(now_minute) <= round_to_minutes_list[i+1] {
				low_round_to_minute = v
				high_round_to_minute = round_to_minutes_list[i+1]
				break
			} else {
				low_round_to_minute = v
				high_round_to_minute = round_to_minutes_list[0] + 60
			}
		}
		if log_level == "Debug" {
			DebugLog.Println("low_round_to_minute is", low_round_to_minute)
		}
		if log_level == "Debug" {
			DebugLog.Println("high_round_to_minute is", high_round_to_minute)
		}
		round_to_duration = time.Duration((high_round_to_minute - low_round_to_minute) * 60 * float64(time.Second))
	} else {
		round_to_duration = time.Duration(round_to_minutes * 60 * float64(time.Second))
	}
	if log_level == "Info" || log_level == "Debug" {
		InfoLog.Println("round_to_duration is", round_to_duration)
	}

	if using_percent {
		round_up_duration = time.Duration(round_up_percent / 100 * float64(round_to_duration))
	} else {
		round_up_duration = time.Duration(round_up_minutes * 60 * float64(time.Second))
	}
	if log_level == "Info" || log_level == "Debug" {
		InfoLog.Println("round_up_duration is", round_up_duration)
	}

	round_up_time := t.Add(round_up_duration)
	round_down_time := t.Add(-round_up_duration)
	if log_level == "Debug" {
		DebugLog.Println("round_up_time is", round_up_time.Round(time.Duration(round_to_duration)).Format(hms_format))
	}
	if log_level == "Debug" {
		DebugLog.Println("round_down_time is", round_down_time.Round(time.Duration(round_to_duration)).Format(hms_format))
	}

	if clear_screen {
		cmd := exec.Command("clear") //works on Darwin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	if round_up_time.Round(time.Duration(round_to_duration)).Sub(t.Add(round_up_duration)) < 0 {
		jime = round_up_time.Round(time.Duration(round_to_duration))
	} else {
		jime = round_down_time.Round(time.Duration(round_to_duration))
	}
	fmt.Println("The jime is", jime.Format(hm_format))
	if log_level == "Info" || log_level == "Debug" {
		InfoLog.Println("jime is", jime.Format(hm_format))
	}
}

func main() {
	default_log_level := "Critical"
	clear_screen, military_display, log_level, loop_seconds, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent := validateConfig(default_log_level)
	if log_level == "Debug" {
		DebugLog.Println("clear_screen is", clear_screen)
	}
	if log_level == "Debug" {
		DebugLog.Println("military_display is", military_display)
	}
	if log_level == "Debug" {
		DebugLog.Println("log_level is", log_level)
	}
	if log_level == "Debug" {
		DebugLog.Println("round_to_minutes is", round_to_minutes)
	}
	if log_level == "Debug" {
		DebugLog.Println("round_to_minutes_list is", round_to_minutes_list)
	}
	if log_level == "Debug" {
		DebugLog.Println("loop_seconds is", loop_seconds)
	}
	if log_level == "Debug" {
		DebugLog.Println("round_up_minutes is", round_up_minutes)
	}
	if log_level == "Debug" {
		DebugLog.Println("round_up_percent is", round_up_percent)
	}

	if military_display {
		hm_format = "15:04"
		hms_format = "15:04:05"
	} else {
		hm_format = "3:04 PM"
		hms_format = "3:04:05 PM"
	}
	if log_level == "Debug" {
		DebugLog.Println("hm_format is", hm_format)
	}
	if log_level == "Debug" {
		DebugLog.Println("hms_format is", hms_format)
	}

	loop_duration = time.Duration(float64(loop_seconds)) * time.Second
	if log_level == "Debug" {
		DebugLog.Println("loop_duration is", loop_duration)
	}
	t := time.Now()
	if log_level == "Info" || log_level == "Debug" {
		InfoLog.Println("* time is", t.Format(hms_format))
	}

	calculateAndDisplayJime(t, clear_screen, log_level, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent)

	if loop_duration != 0 {
		for {
			next_loop_time := t.Round(loop_duration)
			if log_level == "Debug" {
				DebugLog.Println("next_loop_time is", next_loop_time)
			}
			sleep := next_loop_time.Sub(t)

			if sleep < 0 {
				sleep = next_loop_time.Sub(t) + loop_duration
			}
			if log_level == "Debug" {
				DebugLog.Println("sleep is", sleep)
			}

			time.Sleep(sleep)

			t = time.Now()
			if log_level == "Info" || log_level == "Debug" {
				InfoLog.Println("* time is", t.Format(hms_format))
			}

			calculateAndDisplayJime(t, clear_screen, log_level, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent)
		}
	}
}
