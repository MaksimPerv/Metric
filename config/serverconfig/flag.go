// config/serverconfig.go
package serverconfig

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

var FlagRunAddr string
var StoreInterval time.Duration
var FileStoragePath string
var Restore bool

func ParseFlags() {
	var StoreInervalStr string
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "server address")
	flag.StringVar(&StoreInervalStr, "i", "300", "interval server for save metric in file")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "path file storage")
	flag.BoolVar(&Restore, "r", true, "download or not last data")
	flag.Parse()

	// Убедимся, что нет http:// в начале
	FlagRunAddr = strings.TrimPrefix(FlagRunAddr, "http://")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envStoreItercal := os.Getenv("STORE_INTERVAL"); envStoreItercal != "" {
		StoreInervalStr = envStoreItercal
	}
	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		FileStoragePath = envStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		Restore, _ = strconv.ParseBool(envRestore)
	}
	StoreInterval, _ = time.ParseDuration(StoreInervalStr + "s")
}
