package gombokey

import (
	"context"
	"flag"
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/livelace/logrus"
	"os"
	"os/exec"
	"strings"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: false,
		ForceColors:            true,
		ForceQuote:             true,
		FullTimestamp:          true,
		SortingFunc:            SortLogFields,
		TimestampFormat:        DEFAULT_LOG_TIME_FORMAT,
		QuoteEmptyFields:       true,
	})
}

func RunApp() {
	// Greetings.
	log.Info(fmt.Sprintf("%s %s", APP_NAME, APP_VERSION))

	showConfig := flag.Bool("c", false, "Show config sample")
	showDevices := flag.Bool("d", false, "Show input devices")
	flag.Parse()

	if *showConfig {
		fmt.Println(CONFIG_SAMPLE)
		os.Exit(0)
	}

	if *showDevices {
		listDevices()
		os.Exit(0)
	}

	// Get configuration and rules.
	config := getConfig()

	// Set log level.
	ll, err := log.ParseLevel(config.LogLevel)
	log.SetLevel(ll)

	// Open input device.
	dev, err := evdev.Open(config.InputDevice)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(LOG_DEVICE_OPEN_ERROR)
		os.Exit(1)
	}

	// Read device events.
	holdKeys := make(map[string]int64, 0)
	holdKeysSnapshot := make([]string, 0)

	pressedKeys := make([]string, 0)
	resetPressedKeys := func() {
		pressedKeys = make([]string, 0)
	}

	for {
		events, err := dev.Read()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(LOG_DEVICE_OPEN_ERROR)
			os.Exit(1)
		}

		for _, event := range events {
			code := int(event.Code)
			key := evdev.KEY[code]
			state := event.Value
			timestamp := event.Time.Sec

			// save info about hold keys.
			if state == 2 {
				holdKeys[key] = timestamp
				holdKeysSnapshot = ToSort(MapKeysToStringSlice(&holdKeys))
				pressedKeys = DeleteValueFromSlice(pressedKeys, key)
			} else {
				delete(holdKeys, key)
			}

			// save info about pressed keys (only if some hold keys are existed).
			if state == 1 && len(holdKeys) > 0 && key != "KEY_RESERVED" {
				pressedKeys = append(pressedKeys, key)
			}

			// 1. try to detect rule signature and execute command, right after hold keys were released.
			// 2. flush all pressed keys data, if input timeout is reached.
			if len(holdKeys) == 0 {
				// generate current rule signature.
				signature := holdKeysSnapshot

				// append pressed keys to signature.
				for _, key := range pressedKeys {
					signature = append(signature, key)
				}

				// check if generated signature exists inside rules.
				s := strings.Join(signature, "")

				if rule, ok := config.Rules[s]; ok {
					log.WithFields(log.Fields{
						"rule":  rule.Name,
						"hold":  rule.Hold,
						"press": rule.Press,
						"exec":  rule.Exec,
					}).Debug(LOG_RULE_MATCHED)

					go runRule(rule)
				}

				resetPressedKeys()

			} else {
				currentTimestamp := time.Now().Unix()
				earliestTimestamp := currentTimestamp

				for _, timestamp := range holdKeys {
					if timestamp < earliestTimestamp {
						earliestTimestamp = timestamp
					}
				}

				if currentTimestamp-earliestTimestamp > int64(config.InputTimeout) {
					log.WithFields(log.Fields{
						"info": LOG_INPUT_TIMEOUT,
					}).Warn(LOG_MONITOR_EVENTS)

					resetPressedKeys()
				}

				log.WithFields(log.Fields{
					"hold":  MapKeysToStringSlice(&holdKeys),
					"press": pressedKeys,
				}).Debug(LOG_MONITOR_EVENTS)
			}
		}
	}
}

func runRule(rule *Rule) {
	if rule.AddTask() {
		defer rule.DelTask()

		log.WithFields(log.Fields{
			"rule":     rule.Name,
			"parallel": rule.Parallel,
			"task":     rule.ParallelCounter,
		}).Info(LOG_RULE_EXEC)

		// Execute command with timeout.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rule.ExecTimeout)*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, rule.Exec[0], rule.Exec[1:]...)
		output, err := command.Output()

		// timeout.
		if ctx.Err() == context.DeadlineExceeded {
			log.WithFields(log.Fields{
				"rule":  rule.Name,
				"exec":  command,
				"error": ctx.Err().Error(),
			}).Error(LOG_RULE_TIMEOUT)
			return
		}

		// return code is not 0.
		if err != nil {
			log.WithFields(log.Fields{
				"rule":   rule.Name,
				"exec":   command,
				"stdout": string(output),
				"stderr": err,
			}).Error(LOG_RULE_FAILED)
			return

		} else {
			log.WithFields(log.Fields{
				"rule":   rule.Name,
				"exec":   command,
				"stdout": string(output),
				"stderr": err,
			}).Debug(LOG_RULE_EXEC)
		}

		// success.
		log.WithFields(log.Fields{
			"rule":     rule.Name,
			"parallel": rule.Parallel,
			"task":     rule.ParallelCounter,
		}).Infof(LOG_RULE_SUCCESS)

	} else {
		log.WithFields(log.Fields{
			"rule":     rule.Name,
			"parallel": rule.Parallel,
			"task":     rule.ParallelCounter,
		}).Warn(LOG_RULE_LIMIT)
	}
}
