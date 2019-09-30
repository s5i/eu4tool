package main

import (
	"fmt"
	"os"
)

var defaultBackupPath = fmt.Sprintf("%s/Library/Application Support/Steam/SteamApps/common/Europa Universalis IV/save games", os.Getenv("HOME"))
