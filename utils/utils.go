package utils

import (
	"strconv"
)

func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}
