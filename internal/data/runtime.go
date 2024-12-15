package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	runtime := fmt.Sprintf("%d mins", r)

	runtime = strconv.Quote(runtime)

	return []byte(runtime), nil
}
