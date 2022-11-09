package prog

import "github.com/fatih/color"

var MnemonicColor *color.Color
var NumberColor *color.Color
var LabelColor *color.Color
var AddressColor *color.Color
var GrayColor *color.Color

func init() {
	MnemonicColor = color.New(color.FgCyan, color.Bold)
	NumberColor = color.New(color.FgHiRed)
	LabelColor = color.New(color.FgHiGreen, color.Bold)
	AddressColor = color.New(color.FgYellow)
	GrayColor = color.New(color.FgHiBlack)
}
