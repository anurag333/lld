package main

import (
	"fmt"
	"github.com/anurag333/lld/costExplorergpt"
	"github.com/anurag333/lld/middlewarerouter"
	rate_limiter "github.com/anurag333/lld/rate-limiter"
	"github.com/anurag333/lld/snakeandladder"
	"time"
)

func main() {
	//testAllRatelimiters()
	//testCostExplorer()
	//testSnakeandLadder()
	testRouter()
}

func testRouter() {
	// Create a new router instance
	router := middlewarerouter.NewRouter()

	// Add routes
	router.AddRoute("/bar", "result")
	router.AddRoute("/foo/bar", "result2")

	// Call routes
	if result, err := router.CallRoute("/bar"); err == nil {
		fmt.Println(result) // Output: result
	} else {
		fmt.Println(err)
	}

	if result, err := router.CallRoute("/foo/bar"); err == nil {
		fmt.Println(result) // Output: result2
	} else {
		fmt.Println(err)
	}

	if result, err := router.CallRoute("/notfound"); err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err) // Output: route not found: /notfound
	}
}

func testSnakeandLadder() {
	players := []snakeandladder.Player{
		{ID: "P1", Name: "Alice"},
		{ID: "P2", Name: "Bob"},
	}
	snakes := []snakeandladder.Snake{
		{Start: 17, End: 7},
		{Start: 54, End: 34},
	}
	ladders := []snakeandladder.Ladder{
		{Start: 3, End: 22},
		{Start: 15, End: 44},
	}

	game := snakeandladder.NewSnakeAndLadderService(100)
	game.SetPlayers(players)
	game.SetSnakes(snakes)
	game.SetLadders(ladders)
	game.StartGame()
}
func testCostExplorer() {
	pricingPlans := []costExplorergpt.Plan{
		{PlanID: "BASIC", MonthlyCost: 9.99},
		{PlanID: "STANDARD", MonthlyCost: 49.99},
		{PlanID: "PREMIUM", MonthlyCost: 249.99},
	}

	// Create a new CostExplorer with the pricing plans
	costExplorer := costExplorergpt.NewCostExplorer(pricingPlans)

	// Sample customer subscription
	startDate, _ := time.Parse("2006-01-02", "2021-01-01")
	customer := costExplorergpt.Customer{
		CustomerID: "c1",
		Product: costExplorergpt.Product{
			Name: "Jira",
			Subscription: costExplorergpt.Subscription{
				PlanID:    "BASIC",
				StartDate: startDate,
			},
		},
	}

	// Calculate monthly costs and annual cost for the customer
	year := 2021
	monthlyCosts := costExplorer.MonthlyCostList(customer, year)
	annualCost := costExplorer.AnnualCost(customer, year)

	// Print results
	fmt.Printf("Monthly Costs for %d: %v\n", year, monthlyCosts)
	fmt.Printf("Annual Cost for %d: %.2f\n", year, annualCost)
}
func testAllRatelimiters() {
	rl_fixed := rate_limiter.NewFixedWindowRatelimiter(
		rate_limiter.WithMaxRequests(5),
		rate_limiter.WithWindowSize(time.Microsecond*10),
	)
	rl_sliding := rate_limiter.NewSlidingWindowRatelimiter(5, time.Millisecond*20)
	rl_token := rate_limiter.NewTokenBucketRatelimiter(5, time.Microsecond*10)
	testRateLimit(rl_token, "token")
	testRateLimit(rl_fixed, "fixed")
	testRateLimit(rl_sliding, "sliding")

	customerIDs := []string{"customer1", "customer2", "customer3"}

	rl_leaky := rate_limiter.NewLeakyBucketRatelimiter(5, 2)
	for i := 1; i <= 10; i++ {
		for _, customerID := range customerIDs {
			if rl_leaky.AllowRequest(customerID) {
				fmt.Printf("Request %d allowed for %s\n", i, customerID)
			} else {
				fmt.Printf("Request %d denied for %s - bucket is full\n", i, customerID)
			}
		}
		time.Sleep(300 * time.Millisecond) // Simulate some delay between requests
	}

	// stop all leaky buckets before exiting
	rl_leaky.StopAll()
}

func testRateLimit(rl rate_limiter.RateLimiter, algo string) {
	fmt.Println(algo, " Window RL")
	for i := 0; i < 100; i++ {
		if rl.AllowRequest(algo) {
			fmt.Println(algo)
		} else {
			fmt.Println("rejected")
		}
	}
}
