# Order Book Trading

A sample `Order Book Trading` built with `go (golang)`.

## Table of Contents
* [Getting Started](#getting-started)
* [Development Tools](#development-tools)
* [API](#api)
* [Model](#model-struct)
* [Unit Tests](#unit-tests)
* [Author](#author)

## Getting Started

Clone or download as ZIP this project.
```sh
$ git clone https://github.com/agaphetos/order-book-trading.git
```

You can run `main.go` as sample execution of this project. Run using `go`
```sh
$ go main.go
```

### Development Tools

Things you need to make it running

* [go (golang)](https://golang.org/)
* Any TextEditor ([Atom](https://atom.io/), [Sublime](https://www.sublimetext.com/), [VSCode](https://code.visualstudio.com/))
* Any Terminal ([iTerm2](https://www.iterm2.com/), [Hyper](https://hyper.is/))

## API

### **newOrderBookEntry(volume `int32`, orderType `string`, price `float32`, duration `int`) `*OrderBookEntry`**

- constructor function for `OrderBookEntry` stuct
- `returns` an `OrderBookEntry` object

| Parameter | Data Type |
|-----------|-----------|
| volume    | `int32`   |
| orderType | `string`  |
| price     | `float32` |
| duration  | `int`     |

### **newOrderBook() `*OrderBookEntry`**

- constructor function for `OrderBook` with initial dummy data

### **Add(orderBookEntry `OrderBookEntry`)**

- `OrderBook` method to add an entry to an `OrderBook`

| Parameter      | Data Type        |
|----------------|------------------|
| orderBookEntry | `OrderBookEntry` |

### **Cancel(id `string`)**

- `OrderBook` method to cancel or remove an entry to an `OrderBook`

| Parameter | Data Type |
|-----------|-----------|
| id        | `string`  |

### **FulFill(id `string`)**

- `OrderBook` method that sets isFulfilled to true to a given order entry to an `OrderBook`

| Parameter | Data Type |
|-----------|-----------|
| id        | `string`  |

### **(ob `*OrderBook`) View()**

- `OrderBook` method that prints the current entries in the `OrderBook`

### **match(ob `*OrderBook`, entry `OrderBookEntry`) `string`**

- function that matches a given `OrderBookEntry` to a record in an `OrderBook`
- `returns` a string `OrderBookEntry.id` of the matched record

| Parameter      | Data Type        |
|----------------|------------------|
| ob             | `*OrderBook`     |
| orderBookEntry | `OrderBookEntry` |

### **filterOrderBook(orderBook `OrderBook`, orderType `string`) `[]OrderBookEntry`**

- function that filters a given `OrderBook` by `orderType`
- `returns` new `[]OrderBookEntry`

| Parameter | Data Type    |
|-----------|--------------|
| orderBook | `*OrderBook` |
| orderType | `string`     |

### **orderExists(orderBook `OrderBook`, id `string`) `bool`**

- function that checks if an `order` exists in `OrderBook` by `id`
- `returns` `bool`

| Parameter | Data Type    |
|-----------|--------------|
| orderBook | `*OrderBook` |
| id        | `string`     |

## **exchange(ob `*OrderBook`, orders `[]OrderBookEntry`)**

- function processes orders and find a match for each order in the `OrderBook`

| Parameter | Data Type          |
|-----------|--------------------|
| ob        | `*OrderBook`       |
| orders    | `[]OrderBookEntry` |

## **runMatch(ob `*OrderBook`, order `OrderBookEntry`, ch `chan` `bool`)**

- function that wraps `match` function to allow reprocessing of order entry with `timeout`

| Parameter | Data Type        |
|-----------|------------------|
| ob        | `*OrderBook`     |
| order     | `OrderBookEntry` |
| ch        | `ch` `bool`      |

## **saveOrderBook(ob OrderBook)**

- function that writes current OrderBook records to a file

| Parameter | Data Type        |
|-----------|------------------|
| ob        | `OrderBook`     |

## Model (`struct`)

**OrderBook** - a collection of limit orders (`[]OrderBookEntry`)

**OrderBookEntry** - represents an `entry` of an `OrderBook`

| Property    | Data Type | Description               | Value                   |
|-------------|-----------|---------------------------|-------------------------|
| id          | `string`  | identifier of order entry | `UUID`                  |
| volume      | `int32`   | volume of order entry     | `32-bit signed integer` |
| orderType   | `string`  | type of order entry       | `"Ask" or "Bid"`        |
| price       | `float32` | price of order entry      | `32-bit signed float`   |
| isFulfilled | `bool`    | price of order entry      | `32-bit signed float`   |
| expiry      | `int`     | expiration of order entry | `signed int`            |

## Unit Tests

`orderbook_test.go` contains unit tests for the functions of `orderbook.go`

### Code Coverage
- `80.6%`: with `main.go` included
- `88.5%`: with `main.go` not included

### Running Test

- output `code coverage` and `verbose`
```sh
$ go test -cover -v
```

- generating `code coverage report`
```sh
$ go test -coverprofile=coverage.out
```

- view functions coverage using `code coverage report`
```sh
$ go tool cover -func=coverage.out
```

- view html using `code coverage report`
```sh
$ go tool cover -html=coverage.out
```

## Author

* **James Levin Calado** - *Initial work* - [agaphetos](https://github.com/agaphetos)