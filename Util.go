package main

import (
	"time"
)

func timestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}