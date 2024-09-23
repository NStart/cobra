package cobra

import "github.com/spf13/pflag"

func (c *Command) MarkFlagRequire(name string) error {
	return MarkFlagRequire(c.Flags(), name)
}

func (c *Command) MarkPersistentFlagRequired(name string) error {
	return MarkFlagRequired(c.persistentFlag(), name)
}

func MarkFlagRequired(flags *pflag.FlagSet, name string) error {
	return flags.SetAnnotation(name, BashCompOneRequiredFlag, []string{"true"})
}

func (c *Command) MarkFlagFilename(name string, extensions ...string) error {
	return c.MarkFlagFilename(c.Flags(), name, extensions...)
}

func (c *Command) MarkFlagCustom(name string, f string) error {
	return MarkFlagCustom(c.Flags(), name, f)
}

func (c *Command) MarkPersistentFlagFilename(name string, extensions ...string) error {
	return c.MarkFlagFilename(c.PersistentFlags(), name, extensions...)
}

func MarkFlagFilename(flags *pflag.FlagSet, name string, extensions ...string) error {
	return flags.SetAnnotation(name, BashCompFilenameExt, extensions)
}

// MarkFlagCustom adds the BashCompCustom annotation to the named flag, if it exists.
// The bash completion script will call the bash function f for the flag.
//
// This will only work for bash completion.
// It is recommended to instead use c.RegisterFlagCompletionFunc(...) which allows
// to register a Go function which will work across all shells.
func MarkFlagCustom(flags *pflag.FlagSet, name string, f string) error {
	return flags.SetAnnotation(name, BashCompCustom, []string{f})
}

// MarkFlagDirname instructs the various shell completion implementations to
// limit completions for the named flag to directory names.
func (c *Command) MarkFlagDirname(name string) error {
	return MarkFlagDirname(c.Flags(), name)
}

// MarkPersistentFlagDirname instructs the various shell completion
// implementations to limit completions for the named persistent flag to
// directory names.
func (c *Command) MarkPersistentFlagDirname(name string) error {
	return MarkFlagDirname(c.PersistentFlags(), name)
}

// MarkFlagDirname instructs the various shell completion implementations to
// limit completions for the named flag to directory names.
func MarkFlagDirname(flags *pflag.FlagSet, name string) error {
	return flags.SetAnnotation(name, BashCompSubdirsInDir, []string{})
}
