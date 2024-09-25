package main

import "time"

// Condition represents a generic condition structure
type Condition struct {
	Type             string      `json:"type"`
	Attribute        string      `json:"attribute"`
	Operator         string      `json:"operator"`
	Value            interface{} `json:"value"`
	IsValueProcessed bool        `json:"is_value_processed"`
	Aggregator       string      `json:"aggregator"`
	Conditions       []Condition `json:"conditions"`
}

// Cart represents a shopping cart
type Cart struct {
	Items           []Item
	Subtotal        float64
	GrandTotal      float64
	ShippingAddress Address
	BillingAddress  Address
	Customer        Customer
	CouponCode      string
	CreatedAt       time.Time
}

// Item represents a product in the cart
type Item struct {
	SKU          string
	Name         string
	Quantity     int
	Price        float64
	FinalPrice   float64
	SpecialPrice float64
	Weight       float64
	CategoryIDs  []int
	Attributes   map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Address represents a customer address
type Address struct {
	Country    string
	Region     string
	RegionID   int
	City       string
	PostalCode string
	Street     []string
	Telephone  string
	Company    string
	FirstName  string
	LastName   string
	Email      string
}

// Customer represents a customer
type Customer struct {
	ID                 int
	GroupID            int
	Email              string
	FirstName          string
	LastName           string
	Gender             string
	DateOfBirth        time.Time
	CreatedAt          time.Time
	LastLoginAt        time.Time
	Orders             int
	TotalSpent         float64
	AverageOrderAmount float64
	Addresses          []Address
	IsSubscribed       bool
	Attributes         map[string]interface{}
}

// Product represents a product (for more detailed product conditions)
type Product struct {
	ID            int
	SKU           string
	Name          string
	Price         float64
	SpecialPrice  float64
	Weight        float64
	Status        int
	Visibility    int
	TypeID        string
	CategoryIDs   []int
	Attributes    map[string]interface{}
	StockQuantity int
	InStock       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Rule represents a sales rule
type Rule struct {
	ID                  int
	Name                string
	Description         string
	FromDate            time.Time
	ToDate              time.Time
	IsActive            bool
	Conditions          Condition
	Actions             Condition
	StopRulesProcessing bool
	SortOrder           int
	SimpleAction        string
	DiscountAmount      float64
	DiscountQty         float64
	DiscountStep        int
	ApplyToShipping     bool
	TimesUsed           int
	IsRss               bool
	CouponType          int
	UseAutoGeneration   bool
	UsesPerCoupon       int
	UsesPerCustomer     int
}
