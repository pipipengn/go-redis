package database

import "strings"

var cmdTable = make(map[string]*command)

type command struct {
	executor ExecFunc
	argNum   int
}

func RegisterCommand(name string, executor ExecFunc, argNum int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		argNum:   argNum,
	}
}
