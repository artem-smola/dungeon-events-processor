package parser

import (
	"errors"
	"strconv"
	"strings"
)

func ParseTime(raw string) (int, error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 3 {
		return 0, errors.New("time must be HH:MM:SS")
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	s, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	if h < 0 || h > 23 || m < 0 || m > 59 || s < 0 || s > 59 {
		return 0, errors.New("time is out of range")
	}

	return h*3600 + m*60 + s, nil
}
