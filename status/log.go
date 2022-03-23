package status

import (
	"fmt"
	"log"
)

func Log(format string, v ...interface{}) {
	log.Printf("[LeeCahce] %s", fmt.Sprintf(format, v...))
}
