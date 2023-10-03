package logging

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	Debug    *log.Logger
	Info     *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	Critical *log.Logger
)

func init() {
	Debug = log.New(os.Stdout, color.BlackString("DEBUG   : "), log.Ldate|log.Ltime)
	Info = log.New(os.Stdout, color.BlueString("INFO    : "), log.Ldate|log.Ltime)
	Warning = log.New(os.Stdout, color.YellowString("WARNING : "), log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, color.MagentaString("ERROR   : "), log.Ldate|log.Ltime)
	Critical = log.New(os.Stderr, color.RedString("CRITICAL: "), log.Ldate|log.Ltime)
}
