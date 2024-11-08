package models

import "time"

type Customer struct {
	CustomerId        string
	CustomersProducts []Subscription
}

type Product struct {
	ProductId            string
	Name                 string
	Description          string
	EligibleProductPlans map[string]Plan
}

type Plan struct {
	PlanId      string
	MonthlyCost float64
}

type Subscription struct {
	PlanId    string
	ProductId string
	StartDate time.Time
}
