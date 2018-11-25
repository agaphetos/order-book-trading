package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	ob := newOrderBook() // initialize our OrderBook
	ob.View()

	// Create Orders for Question 3 with given data
	var orders []OrderBookEntry
	orders = append(orders, *newOrderBookEntry(50, "Bid", 20., 0))
	orders = append(orders, *newOrderBookEntry(50, "Ask", 20., 5))
	orders = append(orders, *newOrderBookEntry(20, "Bid", 40., 5))
	orders = append(orders, *newOrderBookEntry(100, "Bid", 30., 10))

	// Process Exchange
	exchange(ob, orders)
	saveOrderBook(*ob)

	ob.View()
}
