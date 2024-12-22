package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	runtime := fmt.Sprintf("%d mins", r)

	runtime = strconv.Quote(runtime)

	return []byte(runtime), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unquoteStr, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	splited := strings.Split(unquoteStr, " ")

	if len(splited) != 2 || splited[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(splited[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
