package cobra

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

const (
	requiredAsGroupAnnotation   = "cobra_annotation_required_if_others_set"
	oneRequiredAnnotation       = "cobra_annotation_one_required"
	mutuallyExclusiveAnnotation = "cobra_annotation_mutually_exclusive"
)

func (c *Command) MarkFlagsRequiredTogether(flagNames ...string) {
	c.mergePersistentFlags()
	for _, v := range flagNames {
		f := c.Flags().Lookup(v)
		if f == nil {
			panic(fmt.Sprintf("Failed to find flag %q and mark it as begin required in a flag group", v))
		}
		if err := c.Flags().SetAnnotation(v, requiredAsGroupAnnotation, append(f.Annotations[requiredAsGroupAnnotation], strings.Join(flagNames, " "))); err != nil {
			panic(err)
		}
	}
}

func (c *Command) MarkFlagsOneRequired(flagNames ...string) {
	c.mergePersistentFlags()
	for _, v := range flagNames {
		f := c.Flags().Lookup(v)
		if f == nil {
			panic(fmt.Sprintf("Failed to find flag %q and mark it as begin in a one-required flag group", v))
		}
		if err := c.Flags().SetAnnotation(v, oneRequiredAnnotation, append(f.Annotations[oneRequiredAnnotation], strings.Join(flagNames, " "))); err != nil {
			panic(err)
		}
	}
}

func (c *Command) MarkFlagsMutuallyExclusive(flagNames ...string) {
	c.mergePersistentFlags()
	for _, v := range flagNames {
		f := c.Flags().Lookup(v)
		if f == nil {
			panic(fmt.Sprintf("Failed to find flag %q and mark it as being in a mutually exclusive flag group", v))
		}
		// Each time this is called is a single new entry; this allows it to be a member of multiple groups if needed.
		if err := c.Flags().SetAnnotation(v, mutuallyExclusiveAnnotation, append(f.Annotations[mutuallyExclusiveAnnotation], strings.Join(flagNames, " "))); err != nil {
			panic(err)
		}
	}
}

func (c *Command) ValidateFlagGroups() error {
	if c.DisableFlagParsing {
		return nil
	}

	flags := c.Flags()

	groupStatus := map[string]map[string]bool{}
	oneRequiredGroupStatus := map[string]map[string]bool{}
	mutuallyExclusiveGroupStatus := map[string]map[string]bool{}
	flag.VisitAll(func(pflag *flag.Flag) {
		processFlagGroupAnnotation(flags, pflag, requiredAsGroupAnnotation, groupStatus)
		processFlagGroupAnnotation(flags, pflag, oneRequiredAnnotation, oneRequiredGroupStatus)
		processFlagGroupAnnotation(flags, pflag, mutuallyExclusiveAnnotation, mutuallyExclusiveGroupStatus)
	})

	if err := validateRequireFlagGroups(groupStatus); err != nil {
		return err
	}
	if err := validateOneRequiredFlagGroups(oneRequiredGroupStatus); err != nil {
		return err
	}
	if err := validateExclusiveFlagGroups(mutuallyExclusiveGroupStatus); err != nil {
		return err
	}
	return nil
}

func hasAllFlags(fs *flag.FlagSet, flagnames ...string) bool {
	for _, fname := range flagnames {
		f := fs.Lookup(fname)
		if f == nil {
			return false
		}
	}
	return true
}

func hasAllFlags(fs *flag.FlagSet, flagnames ...string) bool {
	for _, fname := range flagnames {
		f := fs.Lookup(fname)
		if f == nil {
			return false
		}
	}
	return true
}

func processFlagGroupAnnotation(flags *flag.FlagSet, pflag *flag.Flag, annotation string, groupStatus map[string]map[string]bool) {
	groupInfo, found := pflag.Annotations[annotation]
	if found {
		for _, group := range groupInfo {
			if groupStatus[group] == nil {
				flagnames := strings.Split(group, " ")

				if !hasAllFlags(flags, flagnames...) {
					continue
				}

				groupStatus[group] = make(map[string]bool, len(flagnames))
				for _, name := range flagnames {
					groupStatus[group][name] = false
				}
			}

			groupStatus[group][pflag.Name] = pflag.Changed
		}
	}
}

func validateRequireFlagGroups(data map[string]map[string]bool) error {
	keys := sortedKeys(data)
	for _, flagList := range keys {
		flagnameAndStatus := data[flagList]

		unset := []string{}
		for flagname, isSet := range flagnameAndStatus {
			if !isSet {
				unset = append(unset, flagname)
			}
		}
		if len(unset) == len(flagnameAndStatus) || len(unset) == 0 {
			continue
		}

		sort.Strings(unset)
		return fmt.Errorf("if any flags in the group [%v] are set they must all be set; missing %v", flagList, unset)
	}

	return nil
}

func validateOneRequiredFlagGroups(data map[string]map[string]bool) error {
	keys := sortedKeys(data)
	for _, flagList := range keys {
		flagnameAndStatus := data[flagList]
		var set []string
		for flagname, isSet := range flagnameAndStatus {
			if isSet {
				set = append(set, flagname)
			}
		}
		if len(set) >= 1 {
			continue
		}

		sort.Strings(set)
		return fmt.Errorf("at least one of the flags in the group [%v] is required", flagList)
	}
	return nil
}

func validateExclusiveFlagGroups(data map[string]map[string]bool) error {
	keys := sortedKeys(data)
	for _, flagList := range keys {
		flagnameAndStatus := data[flagList]
		var set []string
		for flagname, isSet := range flagnameAndStatus {
			if isSet {
				set = append(set, flagname)
			}
		}
		if len(set) == 0 || len(set) == 1 {
			continue
		}

		// Sort values, so they can be tested/scripted against consistently.
		sort.Strings(set)
		return fmt.Errorf("if any flags in the group [%v] are set none of the others can be; %v were all set", flagList, set)
	}
	return nil
}

func sortedKeys(m map[string]map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (c *Command) enforceFlagGroupsForCompletion() {
	if c.DisableFlagParsing {
		return
	}

	flags := c.Flags()
	groupStatus := map[string]map[string]bool{}
	oneRequiredGroupStatus := map[string]map[string]bool{}
	mutuallyExclusiveGroupStatus := map[string]map[string]bool{}
	c.Flags().VisitAll(func(pflag *flag.Flag) {
		processFlagForGroupAnnotation(flags, pflag, requiredAsGroupAnnotation, groupStatus)
		processFlagForGroupAnnotation(flags, pflag, oneRequiredAnnotation, oneRequiredGroupStatus)
		processFlagForGroupAnnotation(flags, pflag, mutuallyExclusiveAnnotation, mutuallyExclusiveGroupStatus)
	})

	// If a flag that is part of a group is present, we make all the other flags
	// of that group required so that the shell completion suggests them automatically
	for flagList, flagnameAndStatus := range groupStatus {
		for _, isSet := range flagnameAndStatus {
			if isSet {
				// One of the flags of the group is set, mark the other ones as required
				for _, fName := range strings.Split(flagList, " ") {
					_ = c.MarkFlagRequired(fName)
				}
			}
		}
	}

	// If none of the flags of a one-required group are present, we make all the flags
	// of that group required so that the shell completion suggests them automatically
	for flagList, flagnameAndStatus := range oneRequiredGroupStatus {
		isSet := false

		for _, isSet = range flagnameAndStatus {
			if isSet {
				break
			}
		}

		// None of the flags of the group are set, mark all flags in the group
		// as required
		if !isSet {
			for _, fName := range strings.Split(flagList, " ") {
				_ = c.MarkFlagRequired(fName)
			}
		}
	}

	// If a flag that is mutually exclusive to others is present, we hide the other
	// flags of that group so the shell completion does not suggest them
	for flagList, flagnameAndStatus := range mutuallyExclusiveGroupStatus {
		for flagName, isSet := range flagnameAndStatus {
			if isSet {
				// One of the flags of the mutually exclusive group is set, mark the other ones as hidden
				// Don't mark the flag that is already set as hidden because it may be an
				// array or slice flag and therefore must continue being suggested
				for _, fName := range strings.Split(flagList, " ") {
					if fName != flagName {
						flag := c.Flags().Lookup(fName)
						flag.Hidden = true
					}
				}
			}
		}
	}
}
