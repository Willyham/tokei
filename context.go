package tokei

// expressionContext defines the type of expression we're parsing.
type expressionContext int

// Types of ExpressionContext
const (
	minuteContext expressionContext = iota
	hourContext
	dayOfMonthContext
	monthContext
	dayOfWeekContext
)

// Min gets the min value for the context.
func (ex expressionContext) Min() int {
	switch ex {
	case minuteContext, hourContext:
		return 0
	case dayOfWeekContext, monthContext, dayOfMonthContext:
		return 1
	default:
		panic("invalid expression context")
	}
}

// Max gets the max value for the context.
func (ex expressionContext) Max() int {
	switch ex {
	case minuteContext:
		return 59
	case hourContext:
		return 23
	case dayOfWeekContext:
		return 7
	case monthContext:
		return 12
	case dayOfMonthContext:
		return 31
	default:
		panic("invalid expression context")
	}
}
