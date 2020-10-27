package gombokey

const (
	APP_NAME    = "gombokey"
	APP_VERSION = "v1.0.0"

	DEFAULT_DEVICE_PATH     = "/dev/input/event*"
	DEFAULT_EXEC_TIMEOUT    = 60
	DEFAULT_INPUT_TIMEOUT   = 5
	DEFAULT_LOG_LEVEL       = "INFO"
	DEFAULT_LOG_TIME_FORMAT = "02.01.2006 15:04:05"
	DEFAULT_PARALLEL        = 1

	LOG_CONFIG_ERROR         = "config error"
	LOG_DEVICE_NOT_SET       = "input device not set. exit."
	LOG_DEVICE_OPEN_ERROR    = "cannot open input device"
	LOG_INPUT_TIMEOUT        = "input timeout"
	LOG_MONITOR_EVENTS       = "monitor events"
	LOG_NO_DEVICES           = "Cannot find any input devices. Exit!"
	LOG_NO_VALID_RULES       = "no valid rules. exit."
	LOG_RULE_DEBUG           = "rule debug"
	LOG_RULE_ERROR           = "rule error"
	LOG_RULE_EXEC            = "rule exec"
	LOG_RULE_EXEC_NOT_SET    = "exec not set. skip."
	LOG_RULE_FAILED          = "rule failed"
	LOG_RULE_HOLD_NOT_SET    = "hold not set. skip."
	LOG_RULE_LIMIT           = "rule limit"
	LOG_RULE_MATCHED         = "rule matched"
	LOG_RULE_PRESS_NOT_SET   = "press not set. skip."
	LOG_RULE_SIGNATURE_ERROR = "signature error"
	LOG_RULE_SUCCESS         = "rule success"
	LOG_RULE_TIMEOUT         = "rule timeout"
	LOG_VALID_RULES_FOUND    = "valid rules"

	CONFIG_SAMPLE = `  
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
#parallel       = 1                             # global parallel limit.

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
`
)
