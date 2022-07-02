package goNixArgParser

import (
	"strings"
)

func (s *OptionSet) splitMergedToken(token *argToken) (results []*argToken, success bool) {
	flagMap := s.flagMap
	optionMap := s.flagOptionMap
	originalArg := token.text

	if token.kind != undetermArg ||
		len(originalArg) <= len(s.mergeFlagPrefix) ||
		!strings.HasPrefix(originalArg, s.mergeFlagPrefix) {
		return
	}

	if flagMap[originalArg] != nil {
		return
	}

	var prevFlag *Flag
	mergedArgs := originalArg[len(s.mergeFlagPrefix):]
	splittedTokens := make([]*argToken, 0, len(mergedArgs))
	for i, mergedArg := range mergedArgs {
		splittedArg := s.mergeFlagPrefix + string(mergedArg)
		flag := flagMap[splittedArg]

		if flag != nil {
			if !flag.canMerge {
				return
			}
			splittedTokens = append(splittedTokens, newToken(splittedArg, flagArg))
			prevFlag = flag
			continue
		}

		if prevFlag == nil {
			return
		}

		option := optionMap[prevFlag.Name]
		if option == nil || !option.AcceptValue {
			return
		}

		// re-generate standalone flag with values
		splittedTokens[len(splittedTokens)-1] = newToken(prevFlag.Name+mergedArgs[i:], undetermArg)
		break
	}

	return splittedTokens, true
}

func (s *OptionSet) splitMergedTokens(tokens []*argToken) []*argToken {
	results := make([]*argToken, 0, len(tokens))
	for _, originalToken := range tokens {
		splittedTokens, splitted := s.splitMergedToken(originalToken)
		if splitted {
			results = append(results, splittedTokens...)
		} else {
			results = append(results, originalToken)
		}
	}
	return results
}

func (s *OptionSet) splitAssignSignToken(token *argToken) (results []*argToken) {
	results = make([]*argToken, 0, 2)

	text := token.text
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
			if strings.HasPrefix(text, prefix) {
				results = append(results,
					newToken(flagName, flagArg),
					newToken(text[len(prefix):], valueArg),
				)
				return
			}

			assignIndex := strings.Index(text, assignSign)
			if assignIndex <= 0 {
				continue
			}
			prefix = text[0:assignIndex]
			if foundFlag, _ := s.findFlagByPrefix(prefix); foundFlag == flag {
				results = append(results,
					newToken(flagName, flagArg),
					newToken(text[assignIndex+len(assignSign):], valueArg),
				)
				return
			}
		}
	}

	results = append(results, token)
	return
}

func (s *OptionSet) splitAssignSignTokens(tokens []*argToken) []*argToken {
	results := make([]*argToken, 0, len(tokens))

	for _, token := range tokens {
		if token.kind == undetermArg {
			results = append(results, s.splitAssignSignToken(token)...)
		} else {
			results = append(results, token)
		}
	}

	return results
}

func (s *OptionSet) splitConcatAssignToken(token *argToken) (results []*argToken) {
	results = make([]*argToken, 0, 2)

	text := token.text
	for _, flag := range s.flagMap {
		if !flag.canConcatAssign ||
			!s.flagOptionMap[flag.Name].AcceptValue ||
			len(text) <= len(flag.Name) ||
			!strings.HasPrefix(text, flag.Name) {
			continue
		}
		flagName := flag.Name
		flagValue := text[len(flagName):]
		results = append(results,
			newToken(flagName, flagArg),
			newToken(flagValue, valueArg),
		)
		return
	}

	results = append(results, token)
	return
}

func (s *OptionSet) splitConcatAssignTokens(tokens []*argToken) []*argToken {
	results := make([]*argToken, 0, len(tokens))

	for _, token := range tokens {
		if token.kind == undetermArg {
			results = append(results, s.splitConcatAssignToken(token)...)
		} else {
			results = append(results, token)
		}
	}

	return results
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

func (s *OptionSet) parseArgsInGroup(tokens []*argToken) (options map[string][]string, rests, ambigus, undefs []string) {
	options = map[string][]string{}
	rests = []string{}
	ambigus = []string{}
	undefs = []string{}

	flagOptionMap := s.flagOptionMap
	flagMap := s.flagMap

	if s.hasCanMerge {
		tokens = s.splitMergedTokens(tokens)
	}
	if s.hasAssignSigns {
		tokens = s.splitAssignSignTokens(tokens)
	}
	if s.hasCanConcatAssign {
		tokens = s.splitConcatAssignTokens(tokens)
	}

	s.markAmbiguPrefixArgsValues(tokens)
	s.markUndefArgsValues(tokens)

	// walk
	for i, argCount, peeked := 0, len(tokens), 0; i < argCount; i, peeked = i+1+peeked, 0 {
		arg := tokens[i]

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
			options[opt.Key] = []string{}
			continue
		}

		if !opt.MultiValues { // option has 1 value
			if i == argCount-1 || !isValueArg(flag, tokens[i+1]) { // no more value
				if opt.OverridePrev || options[opt.Key] == nil {
					options[opt.Key] = []string{}
				}
			} else {
				if opt.OverridePrev || options[opt.Key] == nil {
					nextArg := tokens[i+1]
					nextArg.kind = valueArg
					options[opt.Key] = []string{nextArg.text}
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

			if !isValueArg(flag, tokens[i+peeked+1]) { // no more value
				break
			}

			peeked++
			peekedArg := tokens[i+peeked]
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

		if opt.OverridePrev || options[opt.Key] == nil {
			options[opt.Key] = values
		} else {
			options[opt.Key] = append(options[opt.Key], values...)
		}
	}

	return options, rests, ambigus, undefs
}

func (s *OptionSet) parseInGroup(specifiedTokens, configTokens []*argToken) *ParseResult {
	keyOptionMap := s.keyOptionMap

	args, argRests, argAmbigus, argUndefs := s.parseArgsInGroup(specifiedTokens)
	envs := s.keyEnvMap
	configs, configRests, configAmbigus, configUndefs := s.parseArgsInGroup(configTokens)
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

func (s *OptionSet) argsToTokensGroups(args []string) (tokensGroups [][]*argToken) {
	tokensGroups = make([][]*argToken, 1)
	groupIndex := 0

	foundRestSign := false
	for _, arg := range args {
		switch {
		case s.isGroupSep(arg):
			tokensGroups = append(tokensGroups, make([]*argToken, 0, 4))
			groupIndex++
			foundRestSign = false
		case foundRestSign:
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, restArg))
		case s.isRestSign(arg):
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, restSignArg))
			foundRestSign = true
		case s.flagMap[arg] != nil:
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, flagArg))
		default:
			tokensGroups[groupIndex] = append(tokensGroups[groupIndex], newToken(arg, undetermArg))
		}
	}

	return
}

func (s *OptionSet) getAlignedTokensGroups(specifiedArgs, configArgs []string) ([][]*argToken, [][]*argToken) {
	specifiedTokensGroups := s.argsToTokensGroups(specifiedArgs)
	specifiedTokensGroupsCount := len(specifiedTokensGroups)

	configTokensGroups := s.argsToTokensGroups(configArgs)
	configTokensGroupsCount := len(configTokensGroups)

	maxCount := specifiedTokensGroupsCount
	if configTokensGroupsCount > maxCount {
		maxCount = configTokensGroupsCount
	}

	for i := specifiedTokensGroupsCount; i < maxCount; i++ {
		specifiedTokensGroups = append(specifiedTokensGroups, []*argToken{})
	}

	for i := configTokensGroupsCount; i < maxCount; i++ {
		configTokensGroups = append(configTokensGroups, []*argToken{})
	}

	return specifiedTokensGroups, configTokensGroups
}

func (s *OptionSet) ParseGroups(specifiedArgs, configArgs []string) []*ParseResult {
	specifiedTokensGroups, configTokensGroups := s.getAlignedTokensGroups(specifiedArgs, configArgs)

	length := len(specifiedTokensGroups)
	results := make([]*ParseResult, length)
	for i := 0; i < length; i++ {
		results[i] = s.parseInGroup(specifiedTokensGroups[i], configTokensGroups[i])
	}

	return results
}

func (s *OptionSet) Parse(specifiedArgs, configArgs []string) *ParseResult {
	specifiedTokensGroups, configTokensGroups := s.getAlignedTokensGroups(specifiedArgs, configArgs)

	var specifiedTokens []*argToken
	if len(specifiedTokensGroups) > 0 {
		specifiedTokens = specifiedTokensGroups[0]
	} else {
		specifiedTokens = []*argToken{}
	}

	var configTokens []*argToken
	if len(configTokensGroups) > 0 {
		configTokens = configTokensGroups[0]
	} else {
		configTokens = []*argToken{}
	}

	result := s.parseInGroup(specifiedTokens, configTokens)

	return result
}
