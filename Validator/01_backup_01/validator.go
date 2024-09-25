package main

import (
	"fmt"
)

type ConditionValidator struct{}

func NewConditionValidator() *ConditionValidator {
	return &ConditionValidator{}
}

func (cv *ConditionValidator) Validate(condition Condition, cart Cart) (bool, error) {
	fmt.Println("Validate " + condition.Type)
	switch condition.Type {
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Combine":
		return cv.validateCombine(condition, cart)
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Product\\Subselect":
		return cv.validateSubselect(condition, cart)
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Product":
		return cv.validateProduct(condition, cart)
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Address":
		return cv.validateAddress(condition, cart)
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer":
		return cv.validateCustomer(condition, cart)
	default:
		return false, fmt.Errorf("unknown condition type: %s", condition.Type)
	}
}

func (cv *ConditionValidator) validateCombine(condition Condition, cart Cart) (bool, error) {
	if len(condition.Conditions) == 0 {
		return true, nil // Empty conditions always validate to true
	}

	switch condition.Aggregator {
	case "all":
		for _, subCondition := range condition.Conditions {
			valid, err := cv.Validate(subCondition, cart)
			if err != nil {
				return false, err
			}
			if !valid {
				return false, nil
			}
		}
		return true, nil

	case "any":
		for _, subCondition := range condition.Conditions {
			valid, err := cv.Validate(subCondition, cart)
			if err != nil {
				return false, err
			}
			if valid {
				return true, nil
			}
		}
		return false, nil

	default:
		return false, fmt.Errorf("unknown aggregator: %s", condition.Aggregator)
	}
}

func (cv *ConditionValidator) validateSubselect(condition Condition, cart Cart) (bool, error) {
	matchingQuantity := 0
	for _, item := range cart.Items {
		itemMatches := true
		for _, subCondition := range condition.Conditions {
			valid, err := cv.validateProduct(subCondition, Cart{Items: []Item{item}})
			if err != nil {
				return false, err
			}
			if !valid {
				itemMatches = false
				break
			}
		}
		if itemMatches {
			matchingQuantity += item.Quantity
		}
	}
	return cv.compareValues(float64(matchingQuantity), condition.Operator, condition.Value)
}

func (cv *ConditionValidator) validateProduct(condition Condition, cart Cart) (bool, error) {
	for _, item := range cart.Items {
		itemValue, err := cv.getItemAttribute(item, condition.Attribute)
		if err != nil {
			return false, err
		}
		valid, err := cv.compareValues(itemValue, condition.Operator, condition.Value)
		if err != nil {
			return false, err
		}
		if valid {
			return true, nil
		}
	}
	return false, nil
}

func (cv *ConditionValidator) validateAddress(condition Condition, cart Cart) (bool, error) {
	addressValue, err := cv.getAddressAttribute(cart.ShippingAddress, condition.Attribute)
	if err != nil {
		return false, err
	}
	return cv.compareValues(addressValue, condition.Operator, condition.Value)
}

func (cv *ConditionValidator) validateCustomer(condition Condition, cart Cart) (bool, error) {
	customerValue, err := cv.getCustomerAttribute(cart.Customer, condition.Attribute)
	if err != nil {
		return false, err
	}
	return cv.compareValues(customerValue, condition.Operator, condition.Value)
}

func (cv *ConditionValidator) compareValues(a interface{}, operator string, b interface{}) (bool, error) {
	switch operator {
	case "==", "!=", ">", ">=", "<", "<=":
		return cv.compareNumericOrString(a, operator, b)
	case "{}", "!{}":
		return cv.compareContains(a, operator, b)
	case "()", "!()":
		return cv.compareInSet(a, operator, b)
	case "null", "notnull":
		return cv.compareNull(a, operator)
	case "like", "nlike":
		return cv.compareLike(a, operator, b)
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

// Implement the following methods in a separate comparisons.go file:
// compareNumericOrString, compareContains, compareInSet, compareNull, compareLike
