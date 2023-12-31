# Jime
### Description
Jime is a program that rounds the time based on configurations you choose.  Configuration options are defined in the file `config_jime.json` in the `jime/` directory.

Jime can be run once or can be configured to loop until an interrupt signal is generated. 

Log entries are written to the file `log_jime.txt` in the `jime/` directory.

### Configuration Options
* Clear Screen
  * choose from true or false
* Military Display
  * choose from true (24-hour) or false (AM/PM)
* Time Zone
  * choose from true or false
* Log Level
  * choose from panic, fatal, error, warn, info, debug, or trace
* Round To options
  * minutes
  * minutes list, for example [0, 15, 30, 45]
* Loop Seconds
  * 0 for no loop
* Round Up options
  * minutes
  * percentage of minutes in current interval

### Configuration Details in config.json
You may use only one non-null Round To option (`round_to_minutes` or `round_to_minutes_list`). The values in round_to_minutes_list must be sorted from low to high, for example [0, 10, ...]. The highest value in the list will be succeeded by a value representing the lowest value of the next hour, for example [0, 30] calculates the jime using 0 round_to_minutes, 30 round_to_minutes, and 0 round_to_minutes the succeeding hour (and repeating).

You may use only one non-null Round Up option (`round_up_minutes` or `round_up_percent`).

#### Examples:
##### Configuration with no loop
```
{
  "clear_screen": false,
  "military_display": true,
  "time_zone": true,
  "log_level": "trace",
  "round_to_minutes": 5,
  "round_to_minutes_list": null,
  "loop_seconds": 0,
  "round_up_minutes": 2,
  "round_up_percent": null
}
```
##### Configuration with 60 second loop, round to minutes list combining 10 and 15 minute intervals, and 40% round up percent
```
{
  "clear_screen": true,
  "military_display": false,
  "time_zone": false,
  "log_level": "fatal",
  "round_to_minutes": null,
  "round_to_minutes_list": [0,10,15,20,30,40,45,50],
  "loop_seconds": 60,
  "round_up_minutes": null,
  "round_up_percent": 40
}
```
##### Configuration with 6 second loop, round to minutes list combining 10 and 15 minute intervals plus 1 minute intervals at 25-29 and 55-59 minutes, and 40% round up percent
```
{
  "clear_screen": true,
  "military_display": false,
  "time_zone": false,
  "log_level": "fatal",
  "round_to_minutes": null,
  "round_to_minutes_list": [0,10,15,20,25,26,27,28,29,30,40,45,50,55,56,57,58,59],
  "loop_seconds": 6,
  "round_up_minutes": null,
  "round_up_percent": 40
}
```
##### Configuration of a standard minute clock with 0 round to minutes and 0 round up minutes, updated once per minute
```
{
  "clear_screen": true,
  "military_display": false,
  "time_zone": false,
  "log_level": "fatal",
  "round_to_minutes": 0,
  "round_to_minutes_list": null,
  "loop_seconds": 60,
  "round_up_minutes": 0,
  "round_up_percent": null
}
```
