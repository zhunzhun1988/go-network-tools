package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
}

func NextIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ""
	}

	for i := 3; i >= 0; i-- {
		tmp, err := strconv.Atoi(parts[i])
		if err != nil {
			return ""
		}
		if tmp < 254 {
			parts[i] = strconv.Itoa(tmp + 1)
			return strings.Join(parts, ".")
		}
		if i != 3 {
			parts[i] = "0"
		} else {
			parts[i] = "1"
		}
	}
	return ""
}
func IPCmp(ip1, ip2 string) int {
	parts1 := strings.Split(ip1, ".")
	parts2 := strings.Split(ip2, ".")
	if len(parts1) != len(parts2) || len(parts1) != 4 {
		return -2
	}
	for i := 0; i <= 3; i++ {
		tmp1, err1 := strconv.Atoi(parts1[i])
		tmp2, err2 := strconv.Atoi(parts2[i])
		if err1 != nil || err2 != nil {
			return -1
		}
		if tmp1 != tmp2 {
			if tmp1 < tmp2 {
				return -1
			} else {
				return 1
			}
		}
	}
	return 0
}
