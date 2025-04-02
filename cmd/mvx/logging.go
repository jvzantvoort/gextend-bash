package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

func IndentStr(num int) string {
	return strings.Repeat(indentchar, num)
}

func StartFunc(functioname string) string {
	indent += 1
	var retv string
	retv = IndentStr(indent - 1)
	retv += functioname
	retv += " START"
	log.Debug(retv)
	return retv
}

func EndFunc(functioname string) string {
	var retv string
	retv = IndentStr(indent - 1)
	retv += functioname
	retv += " END"

	indent -= 1
	log.Debug(retv)
	return retv
}

func StatMsg(msg, status string) string {
	if len(msg) > 80 {
		msg = msg[:70] + " ..."
	}
	statstr := CenterText(status, 12)
	var c func(a ...interface{}) string

	switch strings.ToLower(status) {
	case "success":
		c = color.New(color.FgGreen).SprintFunc()
	case "failure":
		c = color.New(color.FgRed).SprintFunc()
	default:
		c = color.New(color.FgWhite).SprintFunc()
	}

	return fmt.Sprintf("%-80s [%12s]", msg, c(statstr))
}

func LogSuccess(format string, a ...any) {
	retv := IndentStr(indent) + format
	msg := fmt.Sprintf(retv, a...)
	msg = StatMsg(msg, "SUCCESS")
	fmt.Println("        " + msg)
}

func LogFailure(format string, a ...any) {
	retv := IndentStr(indent) + format
	msg := fmt.Sprintf(retv, a...)
	msg = StatMsg(msg, "FAILURE")
	fmt.Println("        " + msg)
}

func LogDebugf(format string, a ...any) {
	retv := IndentStr(indent) + format
	log.Debugf(retv, a...)
}

func LogInfof(format string, a ...any) {
	retv := IndentStr(indent) + format
	log.Infof(retv, a...)
}

func LogFatalf(format string, a ...any) {
	retv := IndentStr(indent) + format
	log.Fatalf(retv, a...)
}

func LogErrorf(format string, a ...any) {
	retv := IndentStr(indent) + format
	log.Errorf(retv, a...)
}

func CurFunc() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	retv := f.Name()
	return strings.Replace(retv, "main.", "", -1)
}
