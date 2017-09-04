package tokei

type ExpressionContext int

const (
	MinuteContext ExpressionContext = iota
	HourContext
	DayOfMonthContext
	MonthContext
	DayOfWeekContext
)

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
