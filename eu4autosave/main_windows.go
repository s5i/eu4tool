package main

import (
	"fmt"
	"os"
)

var defaultBackupPath = fmt.Sprintf("%s\\Documents\\Paradox Interactive\\Europa Universalis IV\\save games", os.Getenv("userprofile"))
