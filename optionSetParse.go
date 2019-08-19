package goNixArgParser

import (
	"strings"
)

func (s *OptionSet) Parse(initArgs []string) *ParseResult {
	flagOptionMap := s.flagOptionMap
	params := map[string][]string{}
	rests := []string{}

	var args []string
	if s.canMergeOption {
		args = make([]string, 0, len(initArgs))
	EACH_ARG:
		for _, arg := range initArgs {
			if len(arg) <= len(s.mergeOptionPrefix) ||
				!strings.HasPrefix(arg, s.mergeOptionPrefix) ||
				flagOptionMap[arg] != nil {
				args = append(args, arg)
				continue
			}

			mergedArgs := arg[len(s.mergeOptionPrefix):]
			splitedArgs := make([]string, 0, len(mergedArgs))
			for _, mergedArg := range mergedArgs {
				splitedArg := s.mergeOptionPrefix + string(mergedArg)
				if flagOptionMap[splitedArg] == nil {
					args = append(args, arg)
					continue EACH_ARG
				}
				splitedArgs = append(splitedArgs, splitedArg)
			}

			args = append(args, splitedArgs...)
		}
	} else {
		args = initArgs
	}

	for i, argCount, peeked := 0, len(args), 0;
		i < argCount;
	i, peeked = i+1+peeked, 0 {
		arg := args[i]
		opt := flagOptionMap[arg]

		if opt == nil {
			rests = append(rests, arg)
			continue;
		}

		if !opt.AcceptValue { // option has no value
			params[opt.Key] = []string{}
		} else if !opt.MultiValues { // option has 1 value
			if i == argCount-1 || flagOptionMap[args[i+1]] != nil { // no more value or next flag found
				if params[opt.Key] == nil {
					params[opt.Key] = []string{}
				}
			} else {
				params[opt.Key] = []string{args[i+1]}
				peeked++
			}
		} else { //option have multi values
			values := []string{}
			for {
				if i+peeked == argCount-1 { // last arg reached
					break
				}

				if flagOptionMap[args[i+peeked+1]] != nil { // next flag found
					break
				}

				peeked++
				value := args[i+peeked]
				if len(opt.Delimiter) == 0 {
					values = append(values, value)
				} else {
					values = append(values, strings.Split(value, opt.Delimiter)...)
				}
			}

			if params[opt.Key] == nil {
				params[opt.Key] = values
			} else {
				params[opt.Key] = append(params[opt.Key], values...)
			}
		}
	}

	defaults := s.keyDefaultMap

	return &ParseResult{
		params:   params,
		defaults: defaults,
		rests:    rests,
	}
}
