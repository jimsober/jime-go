# Jime
### Description
Jime is a program that rounds the time based on configurations you choose.  Configuration options are defined in the file `config.json` in the `jime/` directory.

Jime can be run once or can be configured to loop until an interrupt signal is generated. 

Log entries are written to the file `jime.log` in the `jime/` directory.

### Options
* Clear Screen
  * choose from true or false
* Military Display
  * choose from true (24-hour) or false (AM/PM)
* Round To options
  * minutes
  * minutes list, for example [0, 15, 30, 45]
* Loop Seconds
  * 0 for no loop
* Round Up options
  * minutes
  * percentage of minutes in current interval

### Configuration Details in config.json
You may use only one non-null Round To option (`round_to_min` or `round_to_min_list`). The values in round_to_min_list should be sorted from low to high and the highest value is succeeded by the lowest + 60 (representing the lowest value of the next hour). 

You may use only one non-null Round Up option (`round_up_min` or `round_up_per`).

#### Examples:
```
# Configuration with no loop
{
  "clear_screen": true,
  "military_display": true,
  "round_to_minutes": 5,
  "round_to_minutes_list": null,
  "loop_seconds": 0,
  "round_up_minutes": 2,
  "round_up_percent": null
}
```
```
# Configuration with 60 second loop, round to minutes list combining 10 and 15 minute intervals, and 40% round up percentage
{
  "clear_screen": true,
  "military_display": false,
  "round_to_minutes": null,
  "round_to_minutes_list": [0,10,15,20,30,40,45,50],
  "loop_seconds": 60,
  "round_up_minutes": null,
  "round_up_percent": 40
}
```
