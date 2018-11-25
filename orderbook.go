package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

// OrderBookEntry represents an order record in the OrderBook
type OrderBookEntry struct {
	id          string
	volume      int32
	orderType   string
	price       float32
	isFulfilled bool // for Exam's Question 3 Order Fulfillment
	expiry      int  // for Exam's Question 4 Order Expiry
}

// newOrderBookEntry is the constructor function for OrderBookEntry stuct
func newOrderBookEntry(volume int32, orderType string, price float32, duration int) *OrderBookEntry {
	var id, _ = newUUID()
	orderBookEntry := OrderBookEntry{id, volume, orderType, price, false, duration}

	return &orderBookEntry
}

// OrderBook represents the list of orders
type OrderBook struct {
	entries []OrderBookEntry
}

// newOrderBook is the constructor function for OrderBook with initial data given in Exam's Question 2
func newOrderBook() *OrderBook {
	orderbook := new(OrderBook)
	orderbook.Add(*newOrderBookEntry(100, "Bid", 10., 0))
	orderbook.Add(*newOrderBookEntry(50, "Bid", 20., 0))
	orderbook.Add(*newOrderBookEntry(100, "Bid", 20., 0))
	orderbook.Add(*newOrderBookEntry(30, "Bid", 10., 0))
	orderbook.Add(*newOrderBookEntry(50, "Ask", 20., 0))
	orderbook.Add(*newOrderBookEntry(50, "Ask", 40., 0))
	orderbook.Add(*newOrderBookEntry(100, "Ask", 50., 0))
	orderbook.Add(*newOrderBookEntry(10, "Ask", 70., 0))
	orderbook.Add(*newOrderBookEntry(30, "Bid", 40., 0))
	orderbook.Add(*newOrderBookEntry(30, "Bid", 50., 0))

	return orderbook
}

// Add method adds new entry of order in the OrderBook
func (ob *OrderBook) Add(orderBookEntry OrderBookEntry) {
	ob.entries = append(ob.entries, orderBookEntry)
}

// Cancel method removes an entry of order in the OrderBook
func (ob *OrderBook) Cancel(id string) {
	for i, entry := range ob.entries {
		if id == entry.id {
			ob.entries = append(ob.entries[:i], ob.entries[i+1:]...)
		}
	}
}

// Fulfill method sets isFulfilled to true of an entry order in the OrderBook
func (ob *OrderBook) Fulfill(id string) {
	for i, entry := range ob.entries {
		if id == entry.id {
			ob.entries[i].isFulfilled = true
		}
	}
}

// View method prints the current entries in the OrderBook
func (ob *OrderBook) View() {
	if len(ob.entries) > 0 {
		fmt.Println(" ---------------------------------------------------------------------------- ")
		fmt.Printf("| %-40s | %-5s | %-10s | %-10s |\n", "Order ID", "Type", "Price", "Volume")
		fmt.Println("|------------------------------------------+-------+------------+------------|")
		for _, entry := range ob.entries {
			if !entry.isFulfilled {
				fmt.Printf("| %-40s | %-5s | %10f | %10d |\n", entry.id, entry.orderType, entry.price, entry.volume)
			}
		}
		fmt.Println(" ---------------------------------------------------------------------------- ")
	} else {
		fmt.Println("No entries on the OrderBook available")
	}
}

// match function processes an order entry to find a match in the BookOrder and fulfill it
func match(ob *OrderBook, entry OrderBookEntry) string {
	var matchOrder OrderBookEntry

	switch {
	case strings.ToUpper(entry.orderType) == "BID":
		{
			askOrders := filterOrderBook(*ob, "ASK")

			// sort askOrders by price in ascending order
			sort.Slice(askOrders, func(i, j int) bool {
				return askOrders[i].price < askOrders[j].price
			})

			for _, order := range askOrders {
				if entry.price == order.price {
					matchOrder = order
					break
				}
				if order.price < entry.price {
					matchOrder = order
				}
			}
		}
	case strings.ToUpper(entry.orderType) == "ASK":
		{
			bidOrders := filterOrderBook(*ob, "BID")

			// sort bidOrders by price in descending order
			sort.Slice(bidOrders, func(i, j int) bool {
				return bidOrders[i].price > bidOrders[j].price
			})

			for _, order := range bidOrders {
				if entry.price == order.price {
					matchOrder = order
					break
				}
				if order.price > entry.price {
					matchOrder = order
				}
			}
		}
	}

	return matchOrder.id
}

// filterOrderBook function filters an OrderBook entries on the given orderType
func filterOrderBook(orderBook OrderBook, orderType string) []OrderBookEntry {
	filtered := make([]OrderBookEntry, 0)
	for _, entry := range orderBook.entries {
		if !entry.isFulfilled && strings.ToUpper(entry.orderType) == orderType {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// orderExists function checks if an order exists in OrderBook by id
func orderExists(orderBook OrderBook, id string) bool {
	for _, entry := range orderBook.entries {
		if entry.id == id {
			return true
		}
	}
	return false
}

// exchange function processes an orders and find a match for each order in the OrderBook
func exchange(ob *OrderBook, orders []OrderBookEntry) {
	// each order is processed one at a time
	// blocking execution without using go routine
	for _, order := range orders {
		ch := make(chan bool) // not neccessarry for this scenario
		runMatch(ob, order, ch)
	}

	// each processing of record spawned as a goroutine and executes non-blocking
	// executing this may result conflict output incase a parallel process order have the same match
	// for _, order := range orders {
	// 	ch := make(chan bool)
	// 	go runMatch(ob, order, ch)
	// }
	// time.Sleep(time.Second * 5) // sleep for 10 seconds to show that ticker is running when using goroutines
}

// runMatch function wraps match function to allow reprocessing of order entry with timeout
func runMatch(ob *OrderBook, order OrderBookEntry, ch chan bool) {
	timeout := time.Second * time.Duration(order.expiry)
	tickChan := time.NewTicker(time.Millisecond * 500).C // this ticker will keep matching the order until timeout reaches
	doneChan := make(chan bool)                          // this channel will determine the timeout

	go func() {
		time.Sleep(timeout) // this will be our timeout for Exam's Question 4
		doneChan <- true    // send signal to doneChannel
	}()

	for {
		select {
		case <-tickChan:
			matchID := match(ob, order) // find match of the current order
			if matchID != "" {
				var matchOrder OrderBookEntry
				for _, order := range ob.entries {
					if order.id == matchID {
						matchOrder = order
					}
				}
				newVolume := matchOrder.volume - order.volume
				var newOrder OrderBookEntry

				if newVolume > 0 {
					newOrder = *newOrderBookEntry(newVolume, matchOrder.orderType, matchOrder.price, matchOrder.expiry)
					ob.Add(newOrder)
				} else if newVolume < 0 {
					newVolume = int32(math.Abs(float64(newVolume)))
					newOrder = *newOrderBookEntry(newVolume, order.orderType, order.price, order.expiry)
					ob.Add(newOrder)
				}

				ob.Fulfill(matchID)

				close(ch)
				return
			}
			if !orderExists(*ob, order.id) {
				ob.Add(order)
			}
			// fmt.Println("entry:", order, "match:", matchID) // for debugging
		case <-doneChan:
			ob.Cancel(order.id)
			close(ch)
			return
		}
	}
}

// saveOrderBook function writes current OrderBook records to a file for Exam's Question 5
func saveOrderBook(ob OrderBook) {
	t := time.Now()
	date := fmt.Sprintf("%d%02d%02dT%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
	filename := "OrderBook-" + date + ".txt"

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, entry := range ob.entries {
			record := fmt.Sprintf("%s|%s|%f|%d|%t\n", entry.id, entry.orderType, entry.price, entry.volume, entry.isFulfilled)
			f.WriteString(record)
		}
	}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits;
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random);
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
