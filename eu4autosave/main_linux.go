package main

import (
	"fmt"
	"os"
)

var defaultBackupPath = fmt.Sprintf("%s/.local/share/Paradox Interactive/Europa Universalis IV/save games", os.Getenv("HOME"))
