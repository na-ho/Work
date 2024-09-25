package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestConditionValidator(t *testing.T) {
	validator := NewConditionValidator()

	// Helper function to create a basic cart for testing
	createTestCart := func() Cart {
		return Cart{
			Items: []Item{
				{SKU: "SKU001", Name: "Product 1", Quantity: 2, Price: 10.0, FinalPrice: 9.0, CategoryIDs: []int{1, 2}},
				{SKU: "SKU002", Name: "Product 2", Quantity: 1, Price: 20.0, FinalPrice: 18.0, CategoryIDs: []int{2, 3}},
			},
			Subtotal:        38.0,
			GrandTotal:      40.0,
			ShippingAddress: Address{Country: "US", Region: "CA", City: "Los Angeles", PostalCode: "90001"},
			Customer: Customer{
				ID: 1, GroupID: 2, Email: "test@example.com",
				FirstName: "John", LastName: "Doe",
				CreatedAt: time.Now().AddDate(-1, 0, 0), // Customer created 1 year ago
				Orders:    5, TotalSpent: 500.0,
			},
			CouponCode: "TESTCODE",
			CreatedAt:  time.Now(),
		}
	}

	tests := []struct {
		name      string
		condition string
		want      bool
	}{
		{
			name: "Product SKU equals",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
				"attribute": "sku",
				"operator": "==",
				"value": "SKU001"
			}`,
			want: true,
		},
		{
			name: "Product quantity greater than",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
				"attribute": "quantity",
				"operator": ">",
				"value": "1"
			}`,
			want: true,
		},
		{
			name: "Product category in set",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
				"attribute": "category_ids",
				"operator": "()",
				"value": ["2", "3", "4"]
			}`,
			want: true,
		},
		{
			name: "Cart subtotal greater than or equal",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Address",
				"attribute": "base_subtotal",
				"operator": ">=",
				"value": "35"
			}`,
			want: true,
		},
		{
			name: "Customer group equals",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
				"attribute": "group_id",
				"operator": "==",
				"value": "2"
			}`,
			want: true,
		},
		{
			name: "Combined condition - all true",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Combine",
				"aggregator": "all",
				"conditions": [
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Address",
						"attribute": "country_id",
						"operator": "==",
						"value": "US"
					},
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
						"attribute": "email",
						"operator": "{}",
						"value": "example.com"
					}
				]
			}`,
			want: true,
		},
		{
			name: "Combined condition - any true",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Combine",
				"aggregator": "any",
				"conditions": [
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Address",
						"attribute": "postcode",
						"operator": "==",
						"value": "90000"
					},
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
						"attribute": "email",
						"operator": "{}",
						"value": "example.com"
					}
				]
			}`,
			want: true,
		},
		{
			name: "Product subselect condition",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product\\Subselect",
				"attribute": "qty",
				"operator": ">=",
				"value": "3",
				"conditions": [
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
						"attribute": "category_ids",
						"operator": "()",
						"value": ["1", "2"]
					}
				]
			}`,
			want: true,
		},
		{
			name: "Customer orders count greater than",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
				"attribute": "orders_count",
				"operator": ">",
				"value": "3"
			}`,
			want: true,
		},
		{
			name: "Customer registration date before",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
				"attribute": "created_at",
				"operator": "<",
				"value": "2023-01-01"
			}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var condition Condition
			err := json.Unmarshal([]byte(tt.condition), &condition)
			if err != nil {
				t.Fatalf("Failed to unmarshal condition: %v", err)
			}

			cart := createTestCart()
			got, err := validator.Validate(condition, cart)
			if err != nil {
				t.Fatalf("Validator.Validate() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Validator.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionValidatorNegativeCases(t *testing.T) {
	validator := NewConditionValidator()

	createTestCart := func() Cart {
		return Cart{
			Items: []Item{
				{SKU: "SKU003", Name: "Product 3", Quantity: 1, Price: 30.0, FinalPrice: 28.0, CategoryIDs: []int{4, 5}},
			},
			Subtotal:        28.0,
			GrandTotal:      30.0,
			ShippingAddress: Address{Country: "CA", Region: "ON", City: "Toronto", PostalCode: "M5V 2T6"},
			Customer: Customer{
				ID: 2, GroupID: 1, Email: "test2@otherexample.com",
				FirstName: "Jane", LastName: "Smith",
				CreatedAt: time.Now().AddDate(0, -1, 0), // Customer created 1 month ago
				Orders:    2, TotalSpent: 100.0,
			},
			CouponCode: "",
			CreatedAt:  time.Now(),
		}
	}

	tests := []struct {
		name      string
		condition string
		want      bool
	}{
		{
			name: "Product SKU not equals",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
				"attribute": "sku",
				"operator": "==",
				"value": "SKU001"
			}`,
			want: false,
		},
		{
			name: "Cart subtotal less than",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Address",
				"attribute": "base_subtotal",
				"operator": ">=",
				"value": "50"
			}`,
			want: false,
		},
		{
			name: "Customer group not equals",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
				"attribute": "group_id",
				"operator": "==",
				"value": "2"
			}`,
			want: false,
		},
		{
			name: "Combined condition - all false",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Combine",
				"aggregator": "all",
				"conditions": [
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Address",
						"attribute": "country_id",
						"operator": "==",
						"value": "US"
					},
					{
						"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
						"attribute": "email",
						"operator": "{}",
						"value": "example.com"
					}
				]
			}`,
			want: false,
		},
		{
			name: "Product category not in set",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Product",
				"attribute": "category_ids",
				"operator": "()",
				"value": ["1", "2", "3"]
			}`,
			want: false,
		},
		{
			name: "Customer orders count not greater than",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
				"attribute": "orders_count",
				"operator": ">",
				"value": "5"
			}`,
			want: false,
		},
		{
			name: "Customer registration date not before",
			condition: `{
				"type": "Magento\\SalesRule\\Model\\Rule\\Condition\\Customer",
				"attribute": "created_at",
				"operator": "<",
				"value": "2022-01-01"
			}`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var condition Condition
			err := json.Unmarshal([]byte(tt.condition), &condition)
			if err != nil {
				t.Fatalf("Failed to unmarshal condition: %v", err)
			}

			cart := createTestCart()
			got, err := validator.Validate(condition, cart)
			if err != nil {
				t.Fatalf("Validator.Validate() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Validator.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
