package costexplorer

import (
	"errors"
	"fmt"
	"github.com/anurag333/lld/models"
)

const NUM_MONTHS = 12
const NUM_DAYS = 31

type CostExplorer struct {
	CustomerDetails map[string]models.Customer // customerid-> customerObj
	PricingPlan     map[string]models.Plan     // planId -> planObj
}

func NewCostExplorer() *CostExplorer {
	return &CostExplorer{
		CustomerDetails: make(map[string]models.Customer),
		PricingPlan:     make(map[string]models.Plan),
	}
}

func (costExp *CostExplorer) AddCustomer(customer models.Customer) error {
	_, exists := costExp.CustomerDetails[customer.CustomerId]
	if exists {
		return errors.New("customer already exists")
	}
	costExp.CustomerDetails[customer.CustomerId] = customer
	return nil
}

func (costExp *CostExplorer) MonthlyCostList(customerId string) ([]float64, error) {
	customer, exists := costExp.CustomerDetails[customerId]
	if !exists {
		return nil, errors.New("Customer not registered")
	}
	monthlyExpences := make([]float64, NUM_MONTHS)

	for _, product := range customer.Products {
		subscription := product.SubscriptionObj
		startDate := subscription.StartDate
		planId := subscription.PlanId
		month := startDate.Month() - 1
		day := startDate.Day()
		for i := month + 1; i < NUM_MONTHS; i++ {
			monthlyExpences[i] += costExp.PricingPlan[planId].MonthlyCost
		}
		numDaysToBeCharged := NUM_DAYS - day + 1
		costForCurrentMonth := float64(numDaysToBeCharged) / float64(NUM_DAYS) * costExp.PricingPlan[planId].MonthlyCost
		monthlyExpences[month] += costForCurrentMonth
	}
	return monthlyExpences, nil
}

func (costExp *CostExplorer) MonthlyCostPerProductList(customerId string) (map[string][]float64, error) {
	customer, exists := costExp.CustomerDetails[customerId]
	if !exists {
		return nil, errors.New("Customer not registered")
	}
	monthlyExpences := make(map[string][]float64)

	for _, product := range customer.Products {
		subscription := product.SubscriptionObj
		startDate := subscription.StartDate
		planId := subscription.PlanId
		productName := product.Name
		month := startDate.Month() - 1
		day := startDate.Day()
		_, exists := monthlyExpences[productName]
		if !exists {
			monthlyExpences[productName] = make([]float64, NUM_MONTHS)
		}
		for i := month + 1; i < NUM_MONTHS; i++ {
			monthlyExpences[productName][i] += costExp.PricingPlan[planId].MonthlyCost
		}
		numDaysToBeCharged := NUM_DAYS - day + 1
		costForCurrentMonth := float64(numDaysToBeCharged) / float64(NUM_DAYS) * costExp.PricingPlan[planId].MonthlyCost
		monthlyExpences[productName][month] += costForCurrentMonth
	}
	return monthlyExpences, nil
}

func (costExp *CostExplorer) AnnualCost(customerId string) (float64, error) {
	monthlyExpences, err := costExp.MonthlyCostList(customerId)
	if err != nil {
		fmt.Println("error fetching annual cost : ", err.Error())
		return 0, err
	}
	var annualCost float64 = 0
	for _, cost := range monthlyExpences {
		annualCost += cost
	}
	return annualCost, nil
}
