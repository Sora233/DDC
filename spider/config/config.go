package config

import (
	"github.com/Sora233/DDC/db"
)

type Config struct {
	DB                  db.Config `toml:"db"`
	ClipProcessorNum    int       `toml:"clip_processor_num"`
	CommentProcessorNum int       `toml:"comment_processor_num"`
	RequestLimit        int       `toml:"request_limit"`
	LogLevel            string    `toml:"log_level"`
	UserAgent           string    `toml:"user_agent"`
}

var Global Config
