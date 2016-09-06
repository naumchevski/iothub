package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (c Command) Handle() (counter string, value int, err error) {
	if len(c) == 0 {
		err = errors.New("counter command is not set")
		return
	}

	s := strings.Split(string(c), "=")
	if len(s) != 2 {
		err = errors.New(fmt.Sprintf("wrong counter command format (%s)", c))
		return
	}

	i, err := strconv.Atoi(s[1])
	if err != nil {
		return
	}
	value = int(i)
	counter = s[0]
	return
}
