package config

import (
	"gorm.io/gorm"
)

// DB is the global variable to access the backend DB.
// The variable is supposed to be set once when the web server starts.
var DB *gorm.DB
