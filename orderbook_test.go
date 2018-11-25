package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestOrderBook(t *testing.T) {
	ob := newOrderBook()
	dummyOrder := newOrderBookEntry(50, "Bid", 50., 0)

	t.Run("Function-newOrderBook", func(t *testing.T) { testNewOrderBook(t) })
	t.Run("Function-newOrderBookEntry", func(t *testing.T) { testNewOrderBookEntry(t) })
	t.Run("Method-Add", func(t *testing.T) { testAdd(*ob, *dummyOrder, t) })
	t.Run("Method-Cancel", func(t *testing.T) { testCancel(*ob, *dummyOrder, t) })
	t.Run("Method-Fulfill", func(t *testing.T) { testFulfill(*ob, *dummyOrder, t) })
	t.Run("Function-match", func(t *testing.T) { testMatch(*ob, t) })
	t.Run("Function-filterOrderBook", func(t *testing.T) { testFilterOrderBook(*ob, t) })
	t.Run("Function-orderExists", func(t *testing.T) { testOrderExists(*ob, t) })
	t.Run("Function-exchange", func(t *testing.T) { testExchange(*ob, t) })
	t.Run("Function-saveOrderBook", func(t *testing.T) { testSaveOrderBook(*ob, t) })
	// t.Run("-", func(t *testing.T) {})
}

func testNewOrderBook(t *testing.T) {
	ob := newOrderBook()

	if len(ob.entries) < 1 {
		t.Error("Failed creating new OrderBook")
	}
}

func testNewOrderBookEntry(t *testing.T) {
	tests := [...]struct {
		volume    int32
		orderType string
		price     float32
		duration  int
	}{
		{50, "Bid", 50., 0},
		{50, "Ask", 30., 0},
	}

	for _, tc := range tests {
		newEntry := newOrderBookEntry(tc.volume, tc.orderType, tc.price, tc.duration)
		if newEntry.id == "" {
			t.Error("Failed creating new OrderBookEntry: ", newEntry)
		}
	}
}

func testAdd(ob OrderBook, dummyOrder OrderBookEntry, t *testing.T) {
	ob.Add(dummyOrder)
	exists := false

	for _, entry := range ob.entries {
		if entry.id == dummyOrder.id {
			exists = true
			break
		}
	}

	if !exists {
		t.Fail()
	}
}

func testCancel(ob OrderBook, dummyOrder OrderBookEntry, t *testing.T) {
	ob.Add(dummyOrder)
	ob.Cancel(dummyOrder.id)
	exists := false

	for _, entry := range ob.entries {
		if entry.id == dummyOrder.id {
			exists = true
			break
		}
	}

	if exists {
		t.Fail()
	}
}

func testFulfill(ob OrderBook, dummyOrder OrderBookEntry, t *testing.T) {
	ob.Add(dummyOrder)
	ob.Fulfill(dummyOrder.id)
	fulfilled := false

	for _, entry := range ob.entries {
		if entry.id == dummyOrder.id && entry.isFulfilled {
			fulfilled = true
			break
		}
	}

	if !fulfilled {
		t.Fail()
	}
}

func testMatch(ob OrderBook, t *testing.T) {
	tests := [...]struct {
		volume    int32
		orderType string
		price     float32
		duration  int
		withMatch bool
	}{
		{50, "Bid", 20., 0, true},
		{50, "Bid", 30., 0, true},
		{50, "Bid", 10., 0, false},
		{50, "Ask", 10., 0, true},
		{50, "Ask", 30., 0, true},
		{50, "Ask", 80., 0, false},
	}

	for _, tc := range tests {
		newEntry := newOrderBookEntry(tc.volume, tc.orderType, tc.price, tc.duration)
		matchID := match(&ob, *newEntry)

		if tc.withMatch && matchID == "" {
			t.Errorf("Failed Scenario: %s with match: %v", newEntry.orderType, newEntry)
		}
		if !tc.withMatch && matchID != "" {
			t.Errorf("Failed Scenario: %s without match: %v", newEntry.orderType, newEntry)
		}
	}
}

func testFilterOrderBook(ob OrderBook, t *testing.T) {
	tests := [...]struct {
		orderType string
	}{
		{"ASK"},
		{"BID"},
	}

	for _, tc := range tests {
		filteredOrderBook := filterOrderBook(ob, tc.orderType)

		for _, entry := range filteredOrderBook {
			if entry.isFulfilled || strings.ToUpper(entry.orderType) != tc.orderType {
				t.Error()
			}
		}
	}
}

func testOrderExists(ob OrderBook, t *testing.T) {
	tests := [...]struct {
		volume    int32
		orderType string
		price     float32
		duration  int
		expected  bool
	}{
		{50, "Bid", 20., 0, true},
		{50, "Bid", 30., 0, false},
	}

	for _, tc := range tests {
		newEntry := newOrderBookEntry(tc.volume, tc.orderType, tc.price, tc.duration)

		if tc.expected {
			ob.Add(*newEntry)
		}

		exists := orderExists(ob, newEntry.id)

		if (!exists && tc.expected) || (exists && !tc.expected) {
			t.Fail()
		}
	}
}

func testExchange(ob OrderBook, t *testing.T) {
	var orders []OrderBookEntry
	orders = append(orders, *newOrderBookEntry(50, "Bid", 20., 0))
	orders = append(orders, *newOrderBookEntry(50, "Ask", 20., 5))
	orders = append(orders, *newOrderBookEntry(20, "Bid", 40., 5))
	orders = append(orders, *newOrderBookEntry(100, "Bid", 30., 10))

	exchange(&ob, orders)
}

func testSaveOrderBook(ob OrderBook, t *testing.T) {
	time := time.Now()
	date := fmt.Sprintf("%d%02d%02dT%02d%02d", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute())
	dummyFilename := "OrderBook-" + date + ".txt"

	saveOrderBook(ob)

	if _, err := os.Stat(dummyFilename); os.IsNotExist(err) {
		t.Fail()
	}

	os.Remove(dummyFilename)
}
