package costManager

import (
	"github.com/anurag333/lld/costManager/accessor"
)

func MonthlyCostList(customerId string) ([]float64, error) {
	customer, err := accessor.GetCustomer(customerId)
	if err != nil {
		return nil, err
	}
	monthlyCosts := make([]float64, 12)
	for _, subs := range customer.CustomersProducts {
		plan, err := accessor.GetPlan(subs.PlanId)
		if err != nil {
			return nil, err
		}
		_, month, _ := subs.StartDate.Date()
		for ; month <= 12; month++ {
			monthlyCosts[month] += plan.MonthlyCost
		}
	}
	return monthlyCosts, nil
}
