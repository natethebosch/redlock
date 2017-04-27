package redlock

import (
	"github.com/snksoft/crc"
	"strings"
)

func ComputeHashSlot(key string) int {

	startIndex := strings.IndexRune(key, '{')

	if startIndex != -1 {

		endIndex := strings.IndexRune(key, '}')

		if endIndex != -1 && endIndex > startIndex+1 {
			key = key[startIndex+1 : endIndex]
		}
	}

	return int(crc.CalculateCRC(crc.XMODEM, []byte(key)) % uint64(16384))

}
