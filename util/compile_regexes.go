package util

import (
	"regexp"
)

func CompileRegexes(patterns []string) ([]*regexp.Regexp, error) {
	var err error

	result := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		result[i], err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
