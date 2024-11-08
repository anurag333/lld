package main

import (
	"fmt"
	"github.com/anurag333/lld/costexplorer"
	"github.com/anurag333/lld/models"
	"time"
)

func main() {
	fmt.Println("hello world")
	plan1 := models.Plan{PlanId: "BASIC", MonthlyCost: 9.99}
	plan2 := models.Plan{PlanId: "STANDARD", MonthlyCost: 19.99}
	subscription1 := models.Subscription{PlanId: "BASIC", StartDate: time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC)}
	subscription2 := models.Subscription{PlanId: "STANDARD", StartDate: time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)}

	product1 := models.Product{
		Name:            "JIRA",
		SubscriptionObj: subscription1,
	}
	product2 := models.Product{
		Name:            "CONFLUENCE",
		SubscriptionObj: subscription2,
	}
	productList := make([]models.Product, 0)
	productList = append(productList, product1)
	productList = append(productList, product2)
	customer := models.Customer{
		CustomerId: "peter",
		Products:   productList,
	}

	costExp := costexplorer.NewCostExplorer()
	costExp.CustomerDetails[customer.CustomerId] = customer
	costExp.PricingPlan[plan1.PlanId] = plan1
	costExp.PricingPlan[plan2.PlanId] = plan2

	monthlyCost, _ := costExp.MonthlyCostList("peter")
	for _, cost := range monthlyCost {
		fmt.Print(cost)
		fmt.Print(" ")
	}
	fmt.Println("")

	fmt.Println(costExp.AnnualCost("peter"))

	monthlyCostMap, _ := costExp.MonthlyCostPerProductList("peter")
	for product, monthlyCostList := range monthlyCostMap {
		fmt.Println(product)
		for _, cost := range monthlyCostList {
			fmt.Print(cost)
			fmt.Print(" ")
		}
		fmt.Println("")
	}
	fmt.Println("")

}
