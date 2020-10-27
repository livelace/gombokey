package gombokey

import (
	"fmt"
	log "github.com/livelace/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
)

type Config struct {
	ExecTimeout  int
	Hold         []string
	InputDevice  string
	InputTimeout int
	LogLevel     string
	Parallel     int
	Rules        map[string]*Rule
}

type Rule struct {
	m               sync.Mutex
	Name            string
	Exec            []string
	ExecTimeout     int
	Hold            []string
	Parallel        int
	ParallelCounter int
	Press           []string
}

func (r *Rule) AddTask() bool {
	r.m.Lock()
	defer r.m.Unlock()

	if r.Parallel == 0 || r.ParallelCounter < r.Parallel {
		r.ParallelCounter += 1
		return true
	} else {
		return false
	}
}

func (r *Rule) DelTask() bool {
	r.m.Lock()
	defer r.m.Unlock()

	if r.ParallelCounter <= 0 {
		return false
	} else {
		r.ParallelCounter -= 1
		return true
	}
}

func getConfig() *Config {
	// Get base configuration.
	v := viper.New()
	v.SetConfigName("config.toml")
	v.SetConfigType("toml")
	v.AddConfigPath("/etc/gombokey")
	v.AddConfigPath("$HOME/.gombokey")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(LOG_CONFIG_ERROR)
		os.Exit(1)
	}

	// Set defaults.
	v.SetDefault("default.exec_timeout", DEFAULT_EXEC_TIMEOUT)
	v.SetDefault("default.input_timeout", DEFAULT_INPUT_TIMEOUT)
	v.SetDefault("default.hold", []string{"LEFTALT", "LEFTCTRL"})
	v.SetDefault("default.log_level", DEFAULT_LOG_LEVEL)
	v.SetDefault("default.parallel", DEFAULT_PARALLEL)

	// Set log level.
	ll, err := log.ParseLevel(v.GetString("default.log_level"))
	log.SetLevel(ll)

	inputDevice := v.GetString("default.input_device")
	if inputDevice == "" {
		log.WithFields(log.Fields{
			"error": LOG_DEVICE_NOT_SET,
		}).Error(LOG_CONFIG_ERROR)
		os.Exit(1)
	}

	config := &Config{
		ExecTimeout:  v.GetInt("default.exec_timeout"),
		Hold:         ToUpperAndSort(AppendPrefix(v.GetStringSlice("default.hold"))),
		InputDevice:  inputDevice,
		InputTimeout: v.GetInt("default.input_timeout"),
		LogLevel:     v.GetString("default.log_level"),
		Parallel:     v.GetInt("default.parallel"),
		Rules:        make(map[string]*Rule, 0),
	}

	// Get rules configurations.
	rules := make([]string, 0)
	validRules := make([]string, 0)

	// separate global sections from rules sections.
	for k := range v.AllSettings() {
		if k != "default" {
			rules = append(rules, k)
		}
	}

	for _, rule := range rules {
		exec := v.GetStringSlice(fmt.Sprintf("%s.exec", rule))
		execTimeout := v.GetInt(fmt.Sprintf("%s.exec_timeout", rule))
		hold := ToUpperAndSort(AppendPrefix(v.GetStringSlice(fmt.Sprintf("%s.hold", rule))))
		press := ToUpper(AppendPrefix(v.GetStringSlice(fmt.Sprintf("%s.press", rule))))
		parallel := v.GetInt(fmt.Sprintf("%s.parallel", rule))

		// "hold" might be set globally or individually.
		if len(hold) == 0 {
			hold = config.Hold
		}

		// "exec_timeout" might be set globally or individually.
		if execTimeout == 0 {
			execTimeout = config.ExecTimeout
		}

		// "parallel" might be set globally or individually.
		if parallel == 0 {
			parallel = config.Parallel
		}

		// generate rule signature.
		signature := strings.Join(append(hold, press...), "")

		// Check parameters.
		if len(config.Hold) == 0 && len(hold) == 0 {
			log.WithFields(log.Fields{
				"rule":  rule,
				"error": LOG_RULE_HOLD_NOT_SET,
			}).Warn(LOG_RULE_ERROR)
			continue
		}

		if len(exec) == 0 {
			log.WithFields(log.Fields{
				"rule":  rule,
				"error": LOG_RULE_EXEC_NOT_SET,
			}).Warn(LOG_RULE_ERROR)
			continue
		}

		if len(press) == 0 {
			log.WithFields(log.Fields{
				"rule":  rule,
				"error": LOG_RULE_PRESS_NOT_SET,
			}).Warn(LOG_RULE_ERROR)
			continue
		}

		if _, ok := config.Rules[signature]; ok {
			log.WithFields(log.Fields{
				"rule":      rule,
				"signature": signature,
			}).Warn(LOG_RULE_SIGNATURE_ERROR)
			continue
		}

		validRules = append(validRules, rule)

		config.Rules[signature] = &Rule{
			Exec:        exec,
			ExecTimeout: execTimeout,
			Hold:        hold,
			Name:        rule,
			Parallel:    parallel,
			Press:       press,
		}

		log.WithFields(log.Fields{
			"rule":      rule,
			"hold":      hold,
			"press":     press,
			"signature": signature,
		}).Debug(LOG_RULE_DEBUG)
	}

	if len(config.Rules) == 0 {
		log.Error(LOG_NO_VALID_RULES)
		os.Exit(1)
	} else {
		log.WithFields(log.Fields{
			"rules": validRules,
		}).Info(LOG_VALID_RULES_FOUND)
	}

	return config
}
