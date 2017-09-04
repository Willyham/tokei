package tokei

// ExpressionContext defines the type of expression we're parsing.
type ExpressionContext int

// Types of ExpressionContext
const (
	MinuteContext ExpressionContext = iota
	HourContext
	DayOfMonthContext
	MonthContext
	DayOfWeekContext
)

// Min gets the min value for the context.
func (ex ExpressionContext) Min() int {
	switch ex {
	case MinuteContext, HourContext:
		return 0
	case DayOfWeekContext, MonthContext, DayOfMonthContext:
		return 1
	default:
		panic("invalid expression context")
	}
}

// Max gets the max value for the context.
func (ex ExpressionContext) Max() int {
	switch ex {
	case MinuteContext:
		return 59
	case HourContext:
		return 23
	case DayOfWeekContext:
		return 7
	case MonthContext:
		return 12
	case DayOfMonthContext:
		return 31
	default:
		panic("invalid expression context")
	}
}
