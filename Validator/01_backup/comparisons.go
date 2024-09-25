package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (cv *ConditionValidator) compareNumericOrString(a interface{}, operator string, b interface{}) (bool, error) {
	aFloat, aErr := cv.toFloat64(a)
	bFloat, bErr := cv.toFloat64(b)

	if aErr == nil && bErr == nil {
		switch operator {
		case "==":
			return aFloat == bFloat, nil
		case "!=":
			return aFloat != bFloat, nil
		case ">":
			return aFloat > bFloat, nil
		case ">=":
			return aFloat >= bFloat, nil
		case "<":
			return aFloat < bFloat, nil
		case "<=":
			return aFloat <= bFloat, nil
		}
	}

	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if aOk && bOk {
		switch operator {
		case "==":
			return aStr == bStr, nil
		case "!=":
			return aStr != bStr, nil
		case ">":
			return aStr > bStr, nil
		case ">=":
			return aStr >= bStr, nil
		case "<":
			return aStr < bStr, nil
		case "<=":
			return aStr <= bStr, nil
		}
	}

	return false, fmt.Errorf("unable to compare values: %v %s %v", a, operator, b)
}

func (cv *ConditionValidator) compareContains(a interface{}, operator string, b interface{}) (bool, error) {
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if !aOk || !bOk {
		return false, fmt.Errorf("contains operator requires string values")
	}
	contains := strings.Contains(strings.ToLower(aStr), strings.ToLower(bStr))
	if operator == "{}" {
		return contains, nil
	}
	return !contains, nil
}

func (cv *ConditionValidator) compareInSet(a interface{}, operator string, b interface{}) (bool, error) {
	bSlice, ok := b.([]interface{})
	if !ok {
		return false, fmt.Errorf("in/nin operator requires a slice value")
	}
	for _, v := range bSlice {
		if reflect.DeepEqual(a, v) {
			return operator == "()" || operator == "in", nil
		}
	}
	return operator == "!()" || operator == "nin", nil
}

func (cv *ConditionValidator) compareNull(a interface{}, operator string) (bool, error) {
	isNull := a == nil || reflect.ValueOf(a).IsZero()
	return (operator == "null" && isNull) || (operator == "notnull" && !isNull), nil
}

func (cv *ConditionValidator) compareLike(a interface{}, operator string, b interface{}) (bool, error) {
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if !aOk || !bOk {
		return false, fmt.Errorf("like operator requires string values")
	}
	pattern := strings.ReplaceAll(bStr, "%", ".*")
	matched, err := regexp.MatchString("(?i)"+pattern, aStr)
	if err != nil {
		return false, err
	}
	return (operator == "like" && matched) || (operator == "nlike" && !matched), nil
}

func (cv *ConditionValidator) toFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case time.Time:
		return float64(v.Unix()), nil
	default:
		return 0, fmt.Errorf("unable to convert to float64: %v", v)
	}
}

func (cv *ConditionValidator) compareDate(a interface{}, operator string, b interface{}) (bool, error) {
	aTime, aErr := cv.toTime(a)
	bTime, bErr := cv.toTime(b)

	if aErr != nil || bErr != nil {
		return false, fmt.Errorf("invalid date format: %v or %v", a, b)
	}

	switch operator {
	case "==":
		return aTime.Equal(bTime), nil
	case "!=":
		return !aTime.Equal(bTime), nil
	case ">":
		return aTime.After(bTime), nil
	case ">=":
		return aTime.After(bTime) || aTime.Equal(bTime), nil
	case "<":
		return aTime.Before(bTime), nil
	case "<=":
		return aTime.Before(bTime) || aTime.Equal(bTime), nil
	default:
		return false, fmt.Errorf("unsupported date comparison operator: %s", operator)
	}
}

func (cv *ConditionValidator) toTime(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case time.Time:
		return v, nil
	case string:
		return time.Parse(time.RFC3339, v)
	case int64:
		return time.Unix(v, 0), nil
	case float64:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("unable to convert to time.Time: %v", v)
	}
}
