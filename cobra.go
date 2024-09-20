package cobra

import (
	"strings"
	"text/template"
)

var templateFunc = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"appendIfNotPresent":      appendIfNotPresent,
	"rpad":                    rpad,
	"gt":                      Gt,
	"eq":                      Eq,
}

var initializers []func()
var finalizers []func()

const (
	defaultPrefixMatching   = false
	defaultCommandSorting   = true
	defaultCaseInsensitive  = false
	defaultTraverseRunHooks = false
)


