package valitor

import "github.com/fatih/color"

var Red func(...interface{}) string
var Green func(...interface{}) string
var Yellow func(...interface{}) string
var Blue func(...interface{}) string

func InitColors() {
	Red = color.New(color.FgRed).SprintFunc()
	Green = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Blue = color.New(color.FgBlue).SprintFunc()
}
