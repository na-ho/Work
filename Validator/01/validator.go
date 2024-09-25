package main

import (
	"fmt"
)

// ConditionValidator is the struct that will validate conditions.
type ConditionValidator struct{}

// NewConditionValidator creates a new instance of ConditionValidator.
func NewConditionValidator() *ConditionValidator {
	return &ConditionValidator{}
}

// Validate processes both subconditions and the condition itself.
func (cv *ConditionValidator) Validate(condition Condition, cart Cart) (bool, error) {
	fmt.Println("Validate " + condition.Type)

	// Step 1: Check and validate subconditions if present
	if len(condition.Conditions) > 0 {
		subConditionValid, err := cv.validateNestedConditions(condition, cart)
		if err != nil || !subConditionValid {
			return false, fmt.Errorf("validation failed in subconditions: %v", err)
		}
	}

	// Step 2: Validate the condition itself
	return cv.validateSelfCondition(condition, cart)
}

// validateNestedConditions handles the validation of any subconditions using the aggregator (all/any).
func (cv *ConditionValidator) validateNestedConditions(condition Condition, cart Cart) (bool, error) {
	switch condition.Aggregator {
	case "all":
		for _, subCondition := range condition.Conditions {
			valid, err := cv.Validate(subCondition, cart)
			if err != nil {
				return false, fmt.Errorf("subcondition validation error: %v", err)
			}
			if !valid {
				return false, fmt.Errorf("subcondition %v failed", subCondition)
			}
		}
		return true, nil
	case "any":
		for _, subCondition := range condition.Conditions {
			valid, err := cv.Validate(subCondition, cart)
			if err != nil {
				return false, fmt.Errorf("subcondition validation error: %v", err)
			}
			if valid {
				return true, nil
			}
		}
		return false, fmt.Errorf("none of the subconditions met the 'any' aggregator condition")
	default:
		return false, fmt.Errorf("unknown aggregator: %s", condition.Aggregator)
	}
}

// validateSelfCondition validates the current condition's attributes (after any subconditions have been validated).
func (cv *ConditionValidator) validateSelfCondition(condition Condition, cart Cart) (bool, error) {
	switch condition.Type {
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Product":
		return cv.validateProduct(condition, cart)
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Product\\Subselect":
		return cv.validateSubselect(condition, cart)
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Combine":
		// Handle Combine conditions
		//return cv.validateNestedConditions(condition, cart)
		return true, nil
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Address":
		return cv.validateAddress(condition, cart)
	case "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer":
		return cv.validateCustomer(condition, cart)
	default:
		return false, fmt.Errorf("unknown condition type: %s", condition.Type)
	}
}

// validateProduct validates a product-related condition (e.g., SKU, quantity).
func (cv *ConditionValidator) validateProduct(condition Condition, cart Cart) (bool, error) {
	for _, item := range cart.Items {
		// Get the attribute from the item (e.g., SKU, quantity)
		itemValue, err := cv.getItemAttribute(item, condition.Attribute)
		if err != nil {
			return false, fmt.Errorf("failed to get attribute %s from item: %v", condition.Attribute, err)
		}

		// Compare the item's attribute value with the condition's operator and value
		valid, err := cv.compareValues(itemValue, condition.Operator, condition.Value)
		if err != nil {
			return false, fmt.Errorf("comparison failed for item attribute %s: %v", condition.Attribute, err)
		}
		if valid {
			return true, nil
		}
	}
	return false, fmt.Errorf("no item in the cart matched the condition (attribute: %s, operator: %s, value: %v)", condition.Attribute, condition.Operator, condition.Value)
}

// validateSubselect validates a subselect condition (for subsets of products in the cart).
func (cv *ConditionValidator) validateSubselect(condition Condition, cart Cart) (bool, error) {
	for _, item := range cart.Items {
		itemMatches := true

		// Loop through subconditions
		for _, subCondition := range condition.Conditions {
			// If the subcondition is a Product, validate it with cv.validateProduct
			if subCondition.Type == "Magento\\SalesRule\\Model\\Rule\\Condition\\Product" {
				valid, err := cv.validateProduct(subCondition, cart)
				if err != nil {
					return false, fmt.Errorf("product validation failed for subcondition: %v", err)
				}
				if !valid {
					itemMatches = false
					break
				}
			} else {
				// If it's another condition type (e.g., Subselect, Combine), recursively call Validate
				valid, err := cv.Validate(subCondition, cart)
				if err != nil {
					return false, fmt.Errorf("recursive validation failed for subcondition: %v", err)
				}
				if !valid {
					itemMatches = false
					break
				}
			}
		}

		// If item matches, compare its quantity directly with the condition value (do not sum quantities)
		if itemMatches {
			// Compare the item quantity with the expected value in the condition
			valid, err := cv.compareValues(float64(item.Quantity), condition.Operator, condition.Value)
			if err != nil {
				return false, fmt.Errorf("subselect comparison failed: %v", err)
			}
			if !valid {
				return false, fmt.Errorf("subselect condition failed: expected quantity %s %v but got %v for SKU %s", condition.Operator, condition.Value, item.Quantity, item.SKU)
			}
		}
	}

	return true, nil
}

// validateAddress validates an address-related condition.
func (cv *ConditionValidator) validateAddress(condition Condition, cart Cart) (bool, error) {
	addressValue, err := cv.getAddressAttribute(cart.ShippingAddress, condition.Attribute)
	if err != nil {
		return false, fmt.Errorf("failed to get address attribute %s: %v", condition.Attribute, err)
	}
	valid, err := cv.compareValues(addressValue, condition.Operator, condition.Value)
	if err != nil {
		return false, fmt.Errorf("address comparison failed: %v", err)
	}
	if !valid {
		return false, fmt.Errorf("address condition failed: expected %s %v but got %v", condition.Operator, condition.Value, addressValue)
	}
	return true, nil
}

// validateCustomer validates a customer-related condition.
func (cv *ConditionValidator) validateCustomer(condition Condition, cart Cart) (bool, error) {
	customerValue, err := cv.getCustomerAttribute(cart.Customer, condition.Attribute)
	if err != nil {
		return false, fmt.Errorf("failed to get customer attribute %s: %v", condition.Attribute, err)
	}
	valid, err := cv.compareValues(customerValue, condition.Operator, condition.Value)
	if err != nil {
		return false, fmt.Errorf("customer comparison failed: %v", err)
	}
	if !valid {
		return false, fmt.Errorf("customer condition failed: expected %s %v but got %v", condition.Operator, condition.Value, customerValue)
	}
	return true, nil
}

// compareValues compares values based on the operator (==, >=, <=, etc.)
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
