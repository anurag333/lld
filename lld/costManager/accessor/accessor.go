package accessor

import "github.com/anurag333/lld/costManager/models"
import "errors"

var ProductMap = make(map[string]models.Product)
var PlanMap = make(map[string]models.Plan)
var CustomerMap = make(map[string]models.Customer)

func GetCustomer(customerId string) (*models.Customer, error) {
	customer, exists := CustomerMap[customerId]
	if !exists {
		return nil, errors.New("does not exists")
	}
	return &customer, nil
}

func GetPlan(planId string) (*models.Plan, error) {
	plan, exists := PlanMap[planId]
	if !exists {
		return nil, errors.New("does not exists")
	}
	return &plan, nil
}

func GetProduct(planId string) (*models.Product, error) {
	product, exists := ProductMap[planId]
	if !exists {
		return nil, errors.New("does not exists")
	}
	return &product, nil
}
