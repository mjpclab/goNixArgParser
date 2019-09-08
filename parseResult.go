package goNixArgParser

func (r *ParseResult) HasFlagKey(key string) bool {
	_, found := r.params[key]
	return found
}

func (r *ParseResult) HasFlagValue(key string) bool {
	return len(r.params[key]) > 0
}

func (r *ParseResult) HasEnvKey(key string) bool {
	_, found := r.envs[key]
	return found
}

func (r *ParseResult) HasEnvValue(key string) bool {
	return len(r.envs[key]) > 0
}

func (r *ParseResult) HasDefaultKey(key string) bool {
	_, found := r.defaults[key]
	return found
}

func (r *ParseResult) HasDefaultValue(key string) bool {
	return len(r.defaults[key]) > 0
}

func (r *ParseResult) HasKey(key string) bool {
	return r.HasFlagKey(key) || r.HasEnvKey(key) || r.HasDefaultKey(key)
}

func (r *ParseResult) HasValue(key string) bool {
	return r.HasFlagValue(key) || r.HasEnvValue(key) || r.HasDefaultValue(key)
}

func _getValue(source map[string][]string, key string) (value string, found bool) {
	var values []string
	values, found = source[key]

	if found && len(values) > 0 {
		value = values[0]
	}

	return
}

func (r *ParseResult) GetValue(key string) (value string, found bool) {
	value, found = _getValue(r.params, key)
	if found {
		return
	}

	value, found = _getValue(r.envs, key)
	if found {
		return
	}

	value, found = _getValue(r.defaults, key)
	if found {
		return
	}

	return
}

func _getValues(source map[string][]string, key string) (values []string, found bool) {
	values, found = source[key]
	if found {
		values = copys(values)
		return values, true
	}
	return
}

func (r *ParseResult) GetValues(key string) (values []string, found bool) {
	values, found = _getValues(r.params, key)
	if found {
		return
	}

	values, found = _getValues(r.envs, key)
	if found {
		return
	}

	values, found = _getValues(r.defaults, key)
	if found {
		return
	}

	return
}

func (r *ParseResult) GetRests() []string {
	rests := make([]string, len(r.rests))
	copy(rests, r.rests)
	return rests
}

func copys(input []string) []string {
	if input == nil {
		return nil
	}

	output := make([]string, len(input))
	copy(output, input)
	return output
}
