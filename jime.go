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
	Round_up_per      int
}

func validate_config() (bool, time.Duration, time.Duration, []int, time.Duration, int) {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	if payload.Round_to_min != 0 && len(payload.Round_to_min_list) != 0 {
		log.Fatal("Invalid configuration: only one of round_to_min and round_to_min_list allowed", err)
	} else if payload.Round_to_min == 0 && payload.Round_to_min_list == nil {
		log.Fatal("Invalid configuration: one of round_to_min or round_to_min list is required", err)
	} else {
		if payload.Round_to_min != 0 {
			using_list = false
		} else {
			using_list = true
		}
	}
	if payload.Round_up_min != 0 && payload.Round_up_per != 0 {
		log.Fatal("Invalid configuration: only one of round_up_min and round_up_per allowed", err)
	} else if payload.Round_up_min == 0 && payload.Round_up_per == 0 {
		log.Fatal("Invalid configuration: one of round_up_min or round_up_per is required", err)
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

/*
	 func jime(round_to_min int, round_up_min int) {
		dt := time.Now()
		now_hour := dt.Hour()
		now_min := dt.Minute()
		round_to := 60 * round_to_min
		round_up := 60 * round_up_min
		//t = dt + datetime.timedelta(0,rounding-seconds,-dt.microsecond)
		//	return str(t.hour).zfill(2) + ":" + str(t.minute).zfill(2)
		fmt.Println("The jime is", now_hour, ":", now_min)
	}
*/
func main() {
	clear_screen, loop_sec, round_to_min, round_to_min_list, round_up_min, round_up_per := validate_config()
	log.Println("round_to_min_list is", round_to_min_list)
	log.Println("round_up_per is", round_up_per)
	for {
		t := time.Now()
		log.Println("t is", t.Format("3:04:05 PM"))
		now_minute := t.Minute()
		minute_round_up := now_minute + int(round_up_min)
		hour_round_up := 0
		minute_round_up_mod := math.Mod(float64(minute_round_up), float64(round_to_min))

		if minute_round_up > 59 {
			minute_round_up = minute_round_up - 60
			hour_round_up = 1
		}

		log.Println("using_list is", using_list)
		/*  		if using_list {
			for _, i := range round_to_min_list {
				if round_to_min_list[i] <= now_minute && now_minute <= round_to_min_list[0]+60 {
					low_rtm := round_to_min_list[i]
					high_rtm := round_to_min_list[0] + 60
					round_to_min = high_rtm - low_rtm
				} else if round_to_min_list[i] <= now_minute && now_minute <= round_to_min_list[i+1] {
					low_rtm := round_to_min_list[i]
					high_rtm := round_to_min_list[i+1]
					round_to_min = high_rtm - low_rtm
				}
			}
		}

		*/
		log.Println("using_per is", using_per)
		/*			if using_per {
						round_up_min = round_up_per / 100 * round_to_min
					}
		*/
		if clear_screen {
			cmd := exec.Command("clear") //works on Darwin
			cmd.Stdout = os.Stdout
			cmd.Run()
		}

		//jime(round_to_min, round_up_min)

		if minute_round_up_mod == 0 {
			jime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+hour_round_up, t.Round(time.Duration(round_to_min*time.Minute)).Minute(), t.Second(), t.Nanosecond(), t.Location())
		} else {
			jime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+hour_round_up, minute_round_up-int(minute_round_up_mod), t.Second(), t.Nanosecond(), t.Location())
		}
		fmt.Println("The jime is", jime.Format("3:04 PM"))

		var round time.Duration = loop_sec * time.Second
		next_loop_time := t.Round(round)
		sleep := next_loop_time.Sub(t)
		if sleep < 0 {
			sleep = next_loop_time.Sub(t) + loop_sec*time.Second
		}
		time.Sleep(sleep)
	}
}
