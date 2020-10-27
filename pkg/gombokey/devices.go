//  Original source code: https://github.com/gvalkov/golang-evdev/blob/master/bin/evtest.go
package gombokey

import (
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/livelace/logrus"
	"strings"
)

func listDevices() {
	devices, _ := evdev.ListInputDevices(DEFAULT_DEVICE_PATH)

	lines := make([]string, 0)
	max := 0
	if len(devices) > 0 {
		for i := range devices {
			dev := devices[i]
			str := fmt.Sprintf("%-3d %-20s %-35s %s", i, dev.Fn, dev.Name, dev.Phys)
			if len(str) > max {
				max = len(str)
			}
			lines = append(lines, str)
		}
		fmt.Println()
		fmt.Printf("%-3s %-20s %-35s %s\n", "ID", "Device", "Name", "Phys")
		fmt.Printf(strings.Repeat("-", max) + "\n")
		fmt.Printf(strings.Join(lines, "\n") + "\n")

	} else {
		log.Warn(LOG_NO_DEVICES)
	}
}
