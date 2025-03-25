// config/serverconfig.go
package serverconfig

import (
	"flag"
	"os"
	"strings"
)

var FlagRunAddr string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "server address")
	flag.Parse()

	// Убедимся, что нет http:// в начале
	FlagRunAddr = strings.TrimPrefix(FlagRunAddr, "http://")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

}
