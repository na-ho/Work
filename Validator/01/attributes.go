package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// getItemAttribute retrieves an attribute value from a cart item
func (cv *ConditionValidator) getItemAttribute(item Item, attribute string) (interface{}, error) {
	switch strings.ToLower(attribute) {
	case "sku":
		return item.SKU, nil
	case "price":
		return item.Price, nil
	case "final_price":
		return item.FinalPrice, nil
	case "quantity":
		return item.Quantity, nil
	case "name":
		return item.Name, nil
	case "weight":
		return item.Weight, nil
	case "category_ids":
		return item.CategoryIDs, nil
	case "created_at":
		return item.CreatedAt, nil
	case "updated_at":
		return item.UpdatedAt, nil
	default:
		// Check custom attributes
		if val, ok := item.Attributes[attribute]; ok {
			return val, nil
		}
		// Use reflection for other fields
		return cv.getStructField(item, attribute)
	}
}

// getAddressAttribute retrieves an attribute value from an address
func (cv *ConditionValidator) getAddressAttribute(address Address, attribute string) (interface{}, error) {
	switch strings.ToLower(attribute) {
	case "country":
		return address.Country, nil
	case "region":
		return address.Region, nil
	case "region_id":
		return address.RegionID, nil
	case "city":
		return address.City, nil
	case "postcode":
		return address.PostalCode, nil
	case "street":
		return address.Street, nil
	case "telephone":
		return address.Telephone, nil
	case "company":
		return address.Company, nil
	case "firstname":
		return address.FirstName, nil
	case "lastname":
		return address.LastName, nil
	case "email":
		return address.Email, nil
	default:
		// Use reflection for other fields
		return cv.getStructField(address, attribute)
	}
}

// getCustomerAttribute retrieves an attribute value from a customer
func (cv *ConditionValidator) getCustomerAttribute(customer Customer, attribute string) (interface{}, error) {
	switch strings.ToLower(attribute) {
	case "id":
		return customer.ID, nil
	case "group_id":
		return customer.GroupID, nil
	case "email":
		return customer.Email, nil
	case "firstname":
		return customer.FirstName, nil
	case "lastname":
		return customer.LastName, nil
	case "gender":
		return customer.Gender, nil
	case "dob":
		return customer.DateOfBirth, nil
	case "created_at":
		return customer.CreatedAt, nil
	case "last_login_at":
		return customer.LastLoginAt, nil
	case "orders_count":
		return customer.Orders, nil
	case "total_spent":
		return customer.TotalSpent, nil
	case "average_order_amount":
		return customer.AverageOrderAmount, nil
	case "is_subscribed":
		return customer.IsSubscribed, nil
	default:
		// Check custom attributes
		if val, ok := customer.Attributes[attribute]; ok {
			return val, nil
		}
		// Use reflection for other fields
		return cv.getStructField(customer, attribute)
	}
}

// getCartAttribute retrieves an attribute value from the cart
func (cv *ConditionValidator) getCartAttribute(cart Cart, attribute string) (interface{}, error) {
	switch strings.ToLower(attribute) {
	case "subtotal":
		return cart.Subtotal, nil
	case "grand_total":
		return cart.GrandTotal, nil
	case "coupon_code":
		return cart.CouponCode, nil
	case "created_at":
		return cart.CreatedAt, nil
	case "items_count":
		return len(cart.Items), nil
	case "total_quantity":
		var total int
		for _, item := range cart.Items {
			total += item.Quantity
		}
		return total, nil
	default:
		// Use reflection for other fields
		return cv.getStructField(cart, attribute)
	}
}

// getStructField is a helper function to get a field value using reflection
func (cv *ConditionValidator) getStructField(obj interface{}, field string) (interface{}, error) {
	r := reflect.ValueOf(obj)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	if !f.IsValid() {
		return nil, fmt.Errorf("field not found: %s", field)
	}
	return f.Interface(), nil
}

// formatValue formats the value based on its type for comparison
func (cv *ConditionValidator) formatValue(value interface{}) interface{} {
	switch v := value.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	case []int:
		strSlice := make([]string, len(v))
		for i, num := range v {
			strSlice[i] = fmt.Sprintf("%d", num)
		}
		return strings.Join(strSlice, ",")
	default:
		return v
	}
}
