package circuitbreaker

func ExampleNewBreakerGroup() {
	rateBreaker := NewRateBreaker(0.5, 10)
	group := NewBreakerGroup()
	group.Add("rateBreaker", rateBreaker)
	breaker := group.Get("rateBreaker")
	if breaker == nil {
		return
	}

}
