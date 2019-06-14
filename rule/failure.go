package rule

import "fmt"

func failureHeader(rule Rule, message string) string {
	return fmt.Sprintf("[%s] %s: %s\n", rule.ID(), rule.Name(), message)
}

type Failure interface {
	fmt.Stringer
}

type failure struct {
	msg string
}

func (f *failure) String() string {
	return f.msg
}

//

type BasicFailure struct{}

func NewBasicFailure(rule Rule, message string) Failure {
	return &failure{
		msg: failureHeader(rule, message),
	}
}

//

type ActualExpectedFailure struct{}

func NewActualExpectedFailure(rule Rule, message string, actual interface{}, expected interface{}) Failure {
	msg := fmt.Sprintf("%s\tactual: %s\n\texpected: %s", failureHeader(rule, message), actual, expected)
	return &failure{
		msg: msg,
	}
}

//

type RegexpFailure struct{}

func NewRegexpFailure(rule Rule, message string, pattern interface{}, test string, expectedMatch bool) Failure {
	var msg string

	h := failureHeader(rule, message)

	if expectedMatch {
		msg = fmt.Sprintf("%s\texpected a pattern: %s\n\tto match: %s",
			h,
			pattern,
			test,
		)
	} else {
		msg = fmt.Sprintf("%s\texpected pattern: %s\n\tto not match: %s",
			h,
			pattern,
			test,
		)
	}

	return &failure{
		msg: msg,
	}
}
