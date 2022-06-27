package goNixArgParser

import (
	"strings"
)

func (s *OptionSet) splitMergedArg(arg *argToken) (args []*argToken, success bool) {
	flagMap := s.flagMap
	optionMap := s.flagOptionMap
	argText := arg.text

	if arg.kind != undetermArg ||
		len(argText) <= len(s.mergeFlagPrefix) ||
		!strings.HasPrefix(argText, s.mergeFlagPrefix) {
		return
	}

	if flagMap[argText] != nil {
		return
	}

	var prevFlag *Flag
	mergedArgs := argText[len(s.mergeFlagPrefix):]
	splittedArgs := make([]*argToken, 0, len(mergedArgs))
	for i, mergedArg := range mergedArgs {
		splittedArg := s.mergeFlagPrefix + string(mergedArg)
		flag := flagMap[splittedArg]

		if flag != nil {
			if !flag.canMerge {
				return
			}
			splittedArgs = append(splittedArgs, newArg(splittedArg, flagArg))
			prevFlag = flag
			continue
		}

		if len(splittedArg) <= 1 {
			return
		}

		if prevFlag == nil {
			return
		}

		option := optionMap[prevFlag.Name]
		if option == nil || !option.AcceptValue {
			return
		}

		// re-generate standalone flag with values
		splittedArgs[len(splittedArgs)-1] = newArg(prevFlag.Name+mergedArgs[i:], undetermArg)
		break
	}

	return splittedArgs, true
}

func (s *OptionSet) splitMergedArgs(initArgs []*argToken) []*argToken {
	args := make([]*argToken, 0, len(initArgs))
	for _, arg := range initArgs {
		splittedArgs, splitted := s.splitMergedArg(arg)
		if splitted {
			args = append(args, splittedArgs...)
		} else {
			args = append(args, arg)
		}
	}
	return args
}

func (s *OptionSet) splitAssignSignArg(arg *argToken) (args []*argToken) {
	args = make([]*argToken, 0, 2)

	if arg.kind != undetermArg {
		args = append(args, arg)
		return
	}

	argText := arg.text
	for _, flag := range s.flagMap {
		flagName := flag.Name
		if !s.flagOptionMap[flagName].AcceptValue {
			continue
		}
		for _, assignSign := range flag.assignSigns {
			if len(assignSign) == 0 {
				continue
			}

			prefix := flagName + assignSign
			if strings.HasPrefix(argText, prefix) {
				args = append(args, newArg(flagName, flagArg))
				args = append(args, newArg(argText[len(prefix):], valueArg))
				return
			}

			assignIndex := strings.Index(argText, assignSign)
			if assignIndex <= 0 {
				continue
			}
			prefix = argText[0:assignIndex]
			if foundFlag, _ := s.findFlagByPrefix(prefix); foundFlag == flag {
				args = append(args, newArg(flagName, flagArg))
				args = append(args, newArg(argText[assignIndex+len(assignSign):], valueArg))
				return
			}
		}
	}

	args = append(args, arg)
	return
}

func (s *OptionSet) splitAssignSignArgs(initArgs []*argToken) []*argToken {
	args := make([]*argToken, 0, len(initArgs))

	for _, initArg := range initArgs {
		args = append(args, s.splitAssignSignArg(initArg)...)
	}

	return args
}

func (s *OptionSet) splitConcatAssignArg(arg *argToken) (args []*argToken) {
	args = make([]*argToken, 0, 2)

	if arg.kind != undetermArg {
		args = append(args, arg)
		return
	}

	argText := arg.text
	for _, flag := range s.flagMap {
		if !flag.canConcatAssign ||
			!s.flagOptionMap[flag.Name].AcceptValue ||
			len(argText) <= len(flag.Name) ||
			!strings.HasPrefix(argText, flag.Name) {
			continue
		}
		flagName := flag.Name
		flagValue := argText[len(flagName):]
		args = append(args, newArg(flagName, flagArg))
		args = append(args, newArg(flagValue, valueArg))
		return
	}

	args = append(args, arg)
	return
}

func (s *OptionSet) splitConcatAssignArgs(initArgs []*argToken) []*argToken {
	args := make([]*argToken, 0, len(initArgs))

	for _, initArg := range initArgs {
		args = append(args, s.splitConcatAssignArg(initArg)...)
	}

	return args
}

func (s *OptionSet) markAmbiguPrefixArgsValues(args []*argToken) {
	foundAmbiguFlag := false
	for _, arg := range args {
		if arg.kind != undetermArg {
			foundAmbiguFlag = false
			continue
		}
		actualFlag, ambiguous := s.findFlagByPrefix(arg.text)
		if ambiguous {
			arg.kind = ambiguousFlagArg
			foundAmbiguFlag = true
		} else if actualFlag != nil {
			arg.kind = flagArg
			arg.text = actualFlag.Name
			foundAmbiguFlag = false
		} else if foundAmbiguFlag {
			arg.kind = ambiguousFlagValueArg
		}
	}
}

func (s *OptionSet) markUndefArgsValues(args []*argToken) {
	foundUndefFlag := false
	for _, arg := range args {
		if arg.kind != undetermArg {
			foundUndefFlag = false
			continue
		}
		if s.isUdefFlag(arg.text) {
			arg.kind = undefFlagArg
			foundUndefFlag = true
		} else if foundUndefFlag {
			arg.kind = undefFlagValueArg
		}
	}
}

func isValueArg(flag *Flag, arg *argToken) bool {
	switch arg.kind {
	case valueArg:
		return true
	case undetermArg:
		return flag.canFollowAssign
	default:
		return false
	}
}

func (s *OptionSet) parseArgsInGroup(argObjs []*argToken) (args map[string][]string, rests, ambigus, undefs []string) {
	args = map[string][]string{}
	rests = []string{}
	ambigus = []string{}
	undefs = []string{}

	flagOptionMap := s.flagOptionMap
	flagMap := s.flagMap

	if s.hasCanMerge {
		argObjs = s.splitMergedArgs(argObjs)
	}
	if s.hasAssignSigns {
		argObjs = s.splitAssignSignArgs(argObjs)
	}
	if s.hasCanConcatAssign {
		argObjs = s.splitConcatAssignArgs(argObjs)
	}

	s.markAmbiguPrefixArgsValues(argObjs)
	s.markUndefArgsValues(argObjs)

	// walk
	for i, argCount, peeked := 0, len(argObjs), 0; i < argCount; i, peeked = i+1+peeked, 0 {
		arg := argObjs[i]

		// rests
		if arg.kind == restSignArg {
			continue
		}

		if arg.kind == undetermArg {
			arg.kind = restArg
		}
		if arg.kind == restArg {
			rests = append(rests, arg.text)
			continue
		}

		// ambigus
		if arg.kind == ambiguousFlagValueArg {
			continue
		}

		if arg.kind == ambiguousFlagArg {
			ambigus = append(ambigus, arg.text)
			continue
		}

		// undefs
		if arg.kind == undefFlagValueArg {
			continue
		}

		if arg.kind == undefFlagArg {
			undefs = append(undefs, arg.text)
			continue
		}

		// normal
		opt := flagOptionMap[arg.text]
		flag := flagMap[arg.text]

		if !opt.AcceptValue { // option has no value
			args[opt.Key] = []string{}
			continue
		}

		if !opt.MultiValues { // option has 1 value
			if i == argCount-1 || !isValueArg(flag, argObjs[i+1]) { // no more value
				if opt.OverridePrev || args[opt.Key] == nil {
					args[opt.Key] = []string{}
				}
			} else {
				if opt.OverridePrev || args[opt.Key] == nil {
					nextArg := argObjs[i+1]
					nextArg.kind = valueArg
					args[opt.Key] = []string{nextArg.text}
				}
				peeked++
			}
			continue
		}

		//option have multi values
		values := []string{}
		for {
			if i+peeked == argCount-1 { // last arg reached
				break
			}

			if !isValueArg(flag, argObjs[i+peeked+1]) { // no more value
				break
			}

			peeked++
			peekedArg := argObjs[i+peeked]
			peekedArg.kind = valueArg
			value := peekedArg.text
			var appending []string
			if len(opt.Delimiters) == 0 {
				appending = []string{value}
			} else {
				appending = strings.FieldsFunc(value, opt.isDelimiter)
			}

			if opt.UniqueValues {
				values = appendUnique(values, appending...)
			} else {
				values = append(values, appending...)
			}
		}

		if opt.OverridePrev || args[opt.Key] == nil {
			args[opt.Key] = values
		} else {
			args[opt.Key] = append(args[opt.Key], values...)
		}
	}

	return args, rests, ambigus, undefs
}

func (s *OptionSet) parseInGroup(argObjs, configObjs []*argToken) *ParseResult {
	keyOptionMap := s.keyOptionMap

	args, argRests, argAmbigus, argUndefs := s.parseArgsInGroup(argObjs)
	envs := s.keyEnvMap
	configs, configRests, configAmbigus, configUndefs := s.parseArgsInGroup(configObjs)
	defaults := s.keyDefaultMap

	return &ParseResult{
		keyOptionMap: keyOptionMap,

		args:     args,
		envs:     envs,
		configs:  configs,
		defaults: defaults,

		argRests:    argRests,
		configRests: configRests,

		argAmbigus:    argAmbigus,
		configAmbigus: configAmbigus,

		argUndefs:    argUndefs,
		configUndefs: configUndefs,
	}
}

func (s *OptionSet) getNormalizedArgs(initArgs []string) []*argToken {
	args := make([]*argToken, 0, len(initArgs)+1)

	foundRestSign := false
	for _, arg := range initArgs {
		switch {
		case s.isGroupSep(arg):
			foundRestSign = false
			args = append(args, newArg(arg, groupSepArg))
		case foundRestSign:
			args = append(args, newArg(arg, restArg))
		case s.isRestSign(arg):
			foundRestSign = true
			args = append(args, newArg(arg, restSignArg))
		case s.flagMap[arg] != nil:
			args = append(args, newArg(arg, flagArg))
		default:
			args = append(args, newArg(arg, undetermArg))
		}
	}

	return args
}

func splitArgsIntoGroups(argObjs []*argToken) [][]*argToken {
	argObjs = append(argObjs, newArg("", groupSepArg))

	groups := [][]*argToken{}
	items := []*argToken{}
	for _, argObj := range argObjs {
		if argObj.kind != groupSepArg {
			items = append(items, argObj)
			continue
		}

		groups = append(groups, items)
		items = []*argToken{}
	}

	return groups
}

func (s *OptionSet) getArgsConfigsGroups(initArgs, initConfigs []string) ([][]*argToken, [][]*argToken) {
	args := s.getNormalizedArgs(initArgs)
	argsGroups := splitArgsIntoGroups(args)
	argsGroupsCount := len(argsGroups)

	configs := s.getNormalizedArgs(initConfigs)
	configsGroups := splitArgsIntoGroups(configs)
	configsGroupsCount := len(configsGroups)

	length := argsGroupsCount
	if configsGroupsCount > length {
		length = configsGroupsCount
	}

	for i := 0; i < length-argsGroupsCount; i++ {
		argsGroups = append(argsGroups, []*argToken{})
	}

	for i := 0; i < length-configsGroupsCount; i++ {
		configsGroups = append(configsGroups, []*argToken{})
	}

	return argsGroups, configsGroups
}

func (s *OptionSet) ParseGroups(initArgs, initConfigs []string) []*ParseResult {
	argsGroups, configsGroups := s.getArgsConfigsGroups(initArgs, initConfigs)

	results := []*ParseResult{}
	for i, length := 0, len(argsGroups); i < length; i++ {
		result := s.parseInGroup(argsGroups[i], configsGroups[i])
		results = append(results, result)
	}

	return results
}

func (s *OptionSet) Parse(initArgs, initConfigs []string) *ParseResult {
	argsGroups, configsGroups := s.getArgsConfigsGroups(initArgs, initConfigs)

	var args []*argToken
	if len(argsGroups) > 0 {
		args = argsGroups[0]
	} else {
		args = []*argToken{}
	}

	var configs []*argToken
	if len(configsGroups) > 0 {
		configs = configsGroups[0]
	} else {
		configs = []*argToken{}
	}

	result := s.parseInGroup(args, configs)

	return result
}
