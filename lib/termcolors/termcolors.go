package termcolors

import (
	"github.com/fatih/color"
)

var WarningColor = color.New(color.FgYellow, color.Bold)

var ErrorColor = color.New(color.FgRed, color.Bold)

func hello() string {
	return "hello world"
}

const helloConst = "hello world"
