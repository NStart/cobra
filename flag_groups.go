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
