package transcoding

import (
	"strconv"
	"strings"
)

func ParseFrameRates(f string) int8 {
	slice := strings.Split(f, "/")
	i, err := strconv.ParseInt(slice[0], 10, 16)
	if err != nil {
		return 0
	}

	return int8(i)
}
