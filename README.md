# Jime
### Description
Jime is program that rounds the time the way you configure it to. It can be configured to loop until there is an interrupt signal. Log entries are writtem to the file `jime.log` in the `jime/` directory. Configuration options are defined in the file `config.json` in the `jime/` directory.

### Options
* Clear Screen
  * choose from true or false
* Round To options
  * minutes
  * minutes list, for example [0, 15, 30, 45]
* Loop Seconds
  * 0 for no loop
* Round Up options
  * minutes
  * percentage of minutes in current interval

### Configuration Details in config.json
You may use only one non-null Round To option (`round_to_min` or `round_to_min_list`). The values in round_to_min_list should be sorted from low to high and the highest value is followed by the lowest plus 60 (representing the lowest value of the next hour). 

You may use only one non-null Round Up option (`round_up_min` or `round_up_per`).

#### Examples:
```
# Configuration with no loop
{
  "clear_screen": true,
  "log_level": "CRITICAL",
  "round_to_min": 5,
  "round_to_min_list": null,
  "loop_sec": 0,
  "round_up_min": 2,
  "round_up_per": null
}
```
```
# Configuration with 60 second loop, round to minutes list combining 10 and 15 minute intervals, and 40% round up percentage
{
  "clear_screen": true,
  "log_level": "CRITICAL",
  "round_to_min": null,
  "round_to_min_list": [0,10,15,20,30,40,45,50],
  "loop_sec": 60,
  "round_up_min": null,
  "round_up_per": 40
}
```
