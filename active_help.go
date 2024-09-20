package cobra

import (
	"fmt"
	"os"
)

const (
	activeHelpMarker = "_activeHelp_ "
	// The below values should not be changed: programs will be using them explicitly
	// in their user documentation, and users will be using them explicitly.
	activeHelpEnvVarSuffix  = "ACTIVE_HELP"
	activeHelpGlobalEnvVar  = configEnvVarGlobalPrefix + "_" + activeHelpEnvVarSuffix
	activeHelpGlobalDisable = "0"
)

func AppendActiveHelp(compArry []string, activeHelpStr string) []string {
	return append(compArry, fmt.Sprintf("%s%s", activeHelpMarker, activeHelpStr))
}

func GetActiveHelpConfig(cmd *Command) string {
	activeHelpCfg := os.Getenv(activeHelpGlobalEnvVar)
	if activeHelpConfig {
		activeHelpCfg = os.Getenv(activeHelpEnvVar(cmd.Root().Name()))
	}
	return activeHelpCfg
}

func activeHelpEnvVar(name string) string {
	return configEnvVar(name, activeHelpEnvVarSuffix)
}
