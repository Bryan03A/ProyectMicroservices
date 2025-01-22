// config/config.go
package config

import (
    "log"
    "os"
    "strconv"
)

var (
    AppName = "search-and-filters-service"
    Debug   = true
)

func init() {
    if debug, err := strconv.ParseBool(os.Getenv("DEBUG")); err == nil {
        Debug = debug
    }
    log.Printf("App Name: %s, Debug: %v", AppName, Debug)
}
