package security

// foreground colors
// Red
var Red = "\033[0;31m"
var Green = "\033[0;32m"
var Yellow = "\u001B[33m"
var Blue = "\033[0;34m"
var Magenta = "\033[0;35m"
var Cyan = "\033[0;36m"
var White = "\033[0;37m"
var BlackBold = "\033[1;30m"
var RedBold = "\033[1;31m"
var GreenBold = "\033[1;32m"
var YellowBold = "\033[1;33m"
var BlueBold = "\033[1;34m"
var MagentaBold = "\033[1;35m"
var CyanBold = "\033[1;36m"
var WhiteBold = "\033[1;37m"

// no color
var NC = "\033[0m"

func Nocolor() {
	Red = ""
	Green = ""

	Yellow = ""
	Blue = ""
	Magenta = ""
	Cyan = ""
	White = ""
	BlackBold = ""
	RedBold = ""
	GreenBold = ""
	YellowBold = ""
	BlueBold = ""
	MagentaBold = ""
	CyanBold = ""
	WhiteBold = ""
	NC = ""
}
