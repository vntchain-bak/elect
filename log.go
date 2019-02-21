package elect

import (
	golog "log"
	"os"
)

var log *golog.Logger

func InitLog() {
	log = golog.New(os.Stdout, "", golog.LstdFlags|golog.Lshortfile)
}
