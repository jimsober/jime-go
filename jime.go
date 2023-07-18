// jime.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	zerolog.DurationFieldUnit = time.Second
	zerolog.SetGlobalLevel(zerolog.FatalLevel) //set to 'zerolog.Disabled()' to disable logging
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	file, err := os.OpenFile("log_jime.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Logger = log.Output(file)
}

func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func validateConfig() (bool, bool, string, float64, float64, []float64, float64, float64) {
	content, err := os.ReadFile("./config_jime.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Error when opening file")
	} else {
		log.Debug().Msg("Configuration file opened successfully")
	}

	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal().Err(err).Msg("Error during Unmarshal()")
	} else {
		log.Debug().Msg("Configuration data unmarshaled successfully")
	}

	log_level_values := []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}
	if !isElementExist(log_level_values, payload.Log_level) {
		log.Fatal().Err(err).Msg("Invalid configuration: log_level must be one of 'panic', 'fatal', 'error', 'warn', 'info', 'debug', or 'trace'")
	}
	log_level := payload.Log_level
	if log_level == "panic" {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else if log_level == "fatal" {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if log_level == "error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if log_level == "warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if log_level == "info" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if log_level == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if log_level == "trace" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if payload.Round_to_minutes != 0 && len(payload.Round_to_minutes_list) != 0 {
		log.Fatal().Err(err).Msg("Invalid configuration: both of round_to_minutes and round_to_minutes_list not allowed")
	} else if payload.Round_to_minutes == 0 && payload.Round_to_minutes_list == nil {
		using_list = false
	} else {
		if payload.Round_to_minutes != 0 {
			using_list = false
		} else {
			using_list = true
		}
	}
	log.Debug().Float64("payload.Round_to_minutes", payload.Round_to_minutes).Send()
	log.Debug().Floats64("payload.Round_to_minutes_list", payload.Round_to_minutes_list).Send()
	log.Debug().Bool("using_list", using_list).Send()

	if payload.Round_up_minutes != 0 && payload.Round_up_percent != 0 {
		log.Fatal().Err(err).Msg("Invalid configuration: both of round_up_minutes and round_up_percent not allowed")
	} else if payload.Round_up_minutes == 0 && payload.Round_up_percent == 0 {
		using_percent = false
	} else {
		if payload.Round_up_minutes != 0 {
			using_percent = false
		} else {
			using_percent = true
		}
	}
	log.Debug().Float64("payload.Round_up_minutes", payload.Round_up_minutes).Send()
	log.Debug().Float64("payload.Round_up_percent", payload.Round_up_percent).Send()
	log.Debug().Bool("using_percent", using_percent).Send()

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
	log.Debug().Int("now_minute", now_minute).Send()

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
		log.Info().Float64("** low_round_to_minute", low_round_to_minute).Send()
		log.Info().Float64("** high_round_to_minute", high_round_to_minute).Send()
		round_to_duration = time.Duration((high_round_to_minute - low_round_to_minute) * 60 * float64(time.Second))
	} else {
		round_to_duration = time.Duration(round_to_minutes * 60 * float64(time.Second))
	}
	log.Info().Dur("** round_to_duration, seconds", round_to_duration).Send()

	if using_percent {
		round_up_duration = time.Duration(round_up_percent / 100 * float64(round_to_duration))
	} else {
		round_up_duration = time.Duration(round_up_minutes * 60 * float64(time.Second))
	}
	log.Info().Dur("** round_up_duration, seconds", round_up_duration).Send()

	round_down_time := t.Add(-round_up_duration)
	round_up_time := t.Add(round_up_duration)
	log.Debug().Str("round_down_time", round_down_time.Round(time.Duration(round_to_duration)).Format(hms_format)).Send()
	log.Debug().Str("round_up_time", round_up_time.Round(time.Duration(round_to_duration)).Format(hms_format)).Send()

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
	log.Info().Str("** jime", jime.Format(hm_format)).Send()
	fmt.Println("The jime is", jime.Format(hm_format))
}

func main() {
	clear_screen, military_display, log_level, loop_seconds, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent := validateConfig()
	log.Debug().Bool("clear_screen", clear_screen).Send()
	log.Debug().Bool("military_display", military_display).Send()
	log.Debug().Str("log_level", log_level).Send()
	log.Debug().Float64("round_to_minutes", round_to_minutes).Send()
	log.Debug().Floats64("round_to_minutes_list", round_to_minutes_list).Send()
	log.Debug().Float64("loop_seconds", loop_seconds).Send()
	log.Debug().Float64("round_up_minutes", round_up_minutes).Send()
	log.Debug().Float64("round_up_percent", round_up_percent).Send()

	if military_display {
		hm_format = "15:04"
		hms_format = "15:04:05"
	} else {
		hm_format = "3:04 PM"
		hms_format = "3:04:05 PM"
	}
	log.Debug().Str("hm_format", hm_format).Send()
	log.Debug().Str("hms_format", hms_format).Send()

	loop_duration = time.Duration(float64(loop_seconds)) * time.Second
	log.Debug().Dur("loop_duration, seconds", loop_duration).Send()
	t := time.Now()
	log.Info().Str("* time", t.Format(hms_format)).Send()

	calculateAndDisplayJime(t, clear_screen, log_level, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent)

	if loop_duration != 0 {
		for {
			next_loop_time := t.Round(loop_duration)
			log.Debug().Time("next_loop_time", next_loop_time).Send()
			sleep := next_loop_time.Sub(t)

			if sleep < 0 {
				sleep = next_loop_time.Sub(t) + loop_duration
			}
			log.Debug().Dur("sleep, seconds", sleep).Send()

			time.Sleep(sleep)

			t = time.Now()
			log.Info().Str("* time", t.Format(hms_format)).Send()

			calculateAndDisplayJime(t, clear_screen, log_level, round_to_minutes, round_to_minutes_list, round_up_minutes, round_up_percent)
		}
	}
}
