package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	// Complex condition example
	jsonStr := `
	{
		"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Combine",
		"aggregator": "all",
		"conditions": [
			{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product\\Subselect",
				"attribute": "qty",
				"operator": ">=",
				"value": "2",
				"conditions": [
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
						"attribute": "category_ids",
						"operator": "()",
						"value": ["1", "2", "3"]
					}
				]
			},
			{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Address",
				"attribute": "base_subtotal",
				"operator": ">=",
				"value": "100"
			},
			{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
				"attribute": "group_id",
				"operator": "==",
				"value": "2"
			}
		]
	}`

	var condition Condition
	err := json.Unmarshal([]byte(jsonStr), &condition)
	if err != nil {
		log.Fatalf("Error parsing condition: %v", err)
	}

	// Create a sample cart
	cart := Cart{
		Items: []Item{
			{SKU: "SKU001", Name: "Product 1", Quantity: 2, Price: 50.0, CategoryIDs: []int{1, 4}},
			{SKU: "SKU002", Name: "Product 2", Quantity: 1, Price: 75.0, CategoryIDs: []int{2, 3}},
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
	isValid, err := validator.Validate(condition, cart)
	if err != nil {
		log.Fatalf("Error validating condition: %v", err)
	}

	fmt.Printf("Condition is valid: %v\n", isValid)

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
