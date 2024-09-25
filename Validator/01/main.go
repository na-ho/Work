package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func main() {
	// Start timing the entire process
	startTime := time.Now()

	// Complex condition example
	jsonStr := `
    {
    "type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Combine",
    "attribute": null,
    "operator": null,
    "value": "1",
    "is_value_processed": null,
    "aggregator": "all",
    "conditions": [
        {
            "type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product\\Subselect",
            "attribute": "qty",
            "operator": "==",
            "value": "1",
            "is_value_processed": null,
            "aggregator": "all",
            "conditions": [
                {
                    "type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
                    "attribute": "sku",
                    "operator": "==",
                    "value": "1012096",
                    "is_value_processed": false,
                    "attribute_scope": null
                }
            ]
        },
        {
            "type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product\\Subselect",
            "attribute": "qty",
            "operator": "==",
            "value": "1",
            "is_value_processed": null,
            "aggregator": "all",
            "conditions": [
                {
                    "type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
                    "attribute": "sku",
                    "operator": "==",
                    "value": "1132342",
                    "is_value_processed": false,
                    "attribute_scope": null
                }
            ]
        }
    ]
}`

	// Measure JSON unmarshaling
	unmarshalStart := time.Now()
	var condition Condition
	err := json.Unmarshal([]byte(jsonStr), &condition)
	unmarshalDuration := time.Since(unmarshalStart)

	if err != nil {
		log.Fatalf("Error parsing condition: %v", err)
	}

	// Create a sample cart
	cart := Cart{
		Items: []Item{
			{SKU: "1012096", Name: "Product 1", Quantity: 1, Price: 50.0, CategoryIDs: []int{1, 2, 3}},
			{SKU: "1132342", Name: "Product 2", Quantity: 1, Price: 75.0, CategoryIDs: []int{1, 2, 3}},
		},
		Subtotal: 175.0,
		Customer: Customer{
			ID:      1,
			GroupID: 2,
			Email:   "customer@example.com",
		},
	}

	// Initialize and use the validator
	validator := NewConditionValidator()

	// Measure validation time
	validateStart := time.Now()
	isValid, err := validator.Validate(condition, cart)
	validateDuration := time.Since(validateStart)

	if err != nil {
		log.Fatalf("Error validating condition: %v", err)
	}

	// Calculate total execution time
	totalDuration := time.Since(startTime)

	fmt.Printf("Condition is valid: %v\n", isValid)
	fmt.Printf("\nPerformance Metrics:\n")
	fmt.Printf("JSON Unmarshal time: %v\n", unmarshalDuration)
	fmt.Printf("Validation time: %v\n", validateDuration)
	fmt.Printf("Total execution time: %v\n", totalDuration)

	// Explain the condition
	fmt.Println("\nCondition explanation:")
	fmt.Println("1. There must be at least 2 items in the cart from categories 1, 2, or 3")
	fmt.Println("2. The cart subtotal must be at least $100")
	fmt.Println("3. The customer must belong to group ID 2")

	// Explain the result
	fmt.Println("\nResult explanation:")
	if isValid {
		fmt.Println("The condition is met because:")
		fmt.Println("- There are 2 items in the cart from the specified categories")
		fmt.Println("- The cart subtotal ($175.0) is greater than $100")
		fmt.Println("- The customer belongs to group ID 2")
	} else {
		fmt.Println("The condition is not met. One or more of the following is false:")
		fmt.Println("- There are fewer than 2 items in the cart from the specified categories")
		fmt.Println("- The cart subtotal is less than $100")
		fmt.Println("- The customer does not belong to group ID 2")
	}
}
