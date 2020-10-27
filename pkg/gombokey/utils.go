package gombokey

import (
	"sort"
	"strings"
)

// append "KEY_" prefix to keys names (A -> KEY_A).
func AppendPrefix(s []string) []string {
	var temp []string

	for _, v := range s {
		temp = append(temp, "KEY_"+v)
	}

	return temp
}

func ToUpperAndSort(s []string) []string {
	temp := ToUpper(s)

	return ToSort(temp)
}

// sort hold keys (order independent, Ctrl+Alt || Alt+Ctrl).
func ToSort(s []string) []string {
	sort.Strings(s)

	return s
}

// bring rule keys to uppercase (CtRl -> CTRL).
func ToUpper(s []string) []string {
	var temp []string

	for _, v := range s {
		temp = append(temp, strings.ToUpper(v))
	}

	return temp
}

func DeleteValueFromSlice(s []string, v string) []string {
	temp := make([]string, 0)

	for _, i := range s {
		if i != v {
			temp = append(temp, i)
		}
	}

	return temp
}

func MapKeysToStringSlice(m *map[string]int64) []string {
	temp := make([]string, 0)

	for k := range *m {
		temp = append(temp, k)
	}

	return temp
}

func SortLogFields(s []string) {
	// Ordered fields list.
	order := []string{
		"rule", "parallel", "task",
		"exec", "stdout", "stderr",
		"hold", "press", "signature",
		"error",
	}

	// Mark found fields.
	found := make(map[string]bool, 0)

	for _, v := range s {
		found[v] = true
	}

	// Counter for ordering.
	c := 0

	// Set values according order.
	for _, v := range order {
		if _, ok := found[v]; ok {
			s[c] = v
			c++
		}
	}
}
