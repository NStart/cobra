package cobra

import (
	"sync"

	"github.com/spf13/pflag"
)

const (
	ShellCompRequestCmd       = "__complete"
	ShellCompNoDescRequestCmd = "__completeNoDesc"
)

var flagCompletetionFunctions = map[*pflag.Flag]func(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective){}

var flagCompleteMutex = &sync.RWMutex{}

type ShellCompDirective int

type flagCompError struct {
	subCommand string
	flagName string
}

