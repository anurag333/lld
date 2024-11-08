package models

import "time"

type Plan struct {
	PlanId      string
	MonthlyCost float64
}

type Subscription struct {
	PlanId    string
	StartDate time.Time
}

type Product struct {
	Name            string
	SubscriptionObj Subscription
}

type Customer struct {
	CustomerId string
	Products   []Product
}
