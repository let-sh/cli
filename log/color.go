package log

import "github.com/fatih/color"

var Green = color.New(color.FgGreen).SprintFunc()

var Red = color.New(color.FgRed).SprintFunc()

var Yellow = color.New(color.FgHiYellow).SprintFunc()

var Cyan = color.New(color.FgCyan).SprintFunc()

var CyanUnderline = color.New(color.FgCyan).Add(color.Underline).SprintFunc()

var CyanBold = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
