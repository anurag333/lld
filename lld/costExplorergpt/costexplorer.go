package costExplorergpt

import (
	"fmt"
	"time"
)

// Plan represents a pricing plan
type Plan struct {
	PlanID      string
	MonthlyCost float64
}

// Subscription represents a customer's subscription to a product
type Subscription struct {
	PlanID    string
	StartDate time.Time
}

// Product represents a product that a customer can subscribe to
type Product struct {
	Name         string
	Subscription Subscription
}

// Customer represents a customer with a product subscription
type Customer struct {
	CustomerID string
	Product    Product
}

// CostExplorer provides methods to calculate monthly and annual costs for a customer's subscription
type CostExplorer struct {
	PricingPlans map[string]Plan
}

// NewCostExplorer initializes CostExplorer with the pricing plans
func NewCostExplorer(pricingPlans []Plan) *CostExplorer {
	pricingMap := make(map[string]Plan)
	for _, plan := range pricingPlans {
		pricingMap[plan.PlanID] = plan
	}
	return &CostExplorer{PricingPlans: pricingMap}
}

// monthlyCostList calculates the monthly costs for a customer's subscription within a given year
func (ce *CostExplorer) MonthlyCostList(customer Customer, year int) []float64 {
	monthlyCosts := make([]float64, 12)
	subscription := customer.Product.Subscription
	startMonth := 0

	// Get the plan and subscription start month
	plan, planExists := ce.PricingPlans[subscription.PlanID]
	if !planExists {
		fmt.Println("Plan not found for subscription")
		return monthlyCosts
	}

	// Calculate starting month of billing
	if subscription.StartDate.Year() <= year && subscription.StartDate.Month() <= 12 {
		startMonth = int(subscription.StartDate.Month()) - 1
	}

	// Fill in monthly costs from start month
	for i := startMonth; i < 12; i++ {
		monthlyCosts[i] = plan.MonthlyCost
	}
	return monthlyCosts
}

// annualCost calculates the total annual cost for a customer's subscription within a given year
func (ce *CostExplorer) AnnualCost(customer Customer, year int) float64 {
	monthlyCosts := ce.MonthlyCostList(customer, year)
	var totalCost float64

	// Sum up the monthly costs to get the annual cost
	for _, cost := range monthlyCosts {
		totalCost += cost
	}
	return totalCost
}
