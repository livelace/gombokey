# gombokey


***gombokey*** ("go" + "combo" + "key) is a tool for binding specific actions to keys combinations, but with some additional features.

### Features:

* "Hold and press" approach instead of classic shortcuts. Example: **"LEFTALT"**, **"A + A"**. Where combination consists of: **"LEFTALT"** - hold key,  **"A + A"** - pressed keys. Action will be performed right after "hold keys" **are released**.
* "Pressed keys" are ordered. Example: **"A + B + A"** and **"A + A + B"** = are different. 
* "Hold keys" can be hold arbitrary. Example: **"LEFTALT + LEFTCTRL"** and **"LEFTCTRL + LEFTALT"** = the same.
* Parallel executions limits. Example: you can restrict task with only 3 parallel executions.
* Execution timeout. Example: task will be killed, if it exceeds the specific amount of time.

### Quick start:

```shell script
go get github.com/livelace/gombokey/cmd/gombokey
```

### Config sample:

```toml
# Configuration file may be placed:
# 1. /etc/gombokey/config.toml
# 2. $HOME/.gombokey/config.toml
# 3. $(pwd)/config.toml

[default]

input_device    = "/dev/input/event9"           # path to device.
input_timeout   = 5                             # seconds.

log_level       = "INFO"

hold            = ["LEFTALT"]                   # global hold keys.
#exec_timeout   = 60                            # global execution timeout.
#parallel       = 1                             # global unrestricted parallel execution.

[rule1]
hold            = ["LEFTALT", "LEFTCTRL"]       # independent order, case insensitive.
press           = ["A", "Q"]                    # dependent order, case insensitive.

exec            = ["bash", "-c", "echo 123"]    # command and arguments for execution.

exec_timeout    = 10
parallel        = 3


[rule2]
# LEFTALT + B + B + B
press           = ["b", "B", "b"]
exec            = ["beep"]
```