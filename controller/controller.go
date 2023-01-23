package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	amount    float64
	timeStamp time.Time
	city      string
}

type Stats struct {
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
	Count float64 `json:"count"`
}

type Loc struct {
	City string `json:"city"`
}

type Endpoint struct {
	transactions []*Transaction
	cityStats map[string]*Stats
}

func NewEndpoint() *Endpoint{
	return &Endpoint{
		transactions: []*Transaction{},
		cityStats: map[string]*Stats{},
	}
}

func (e *Endpoint) ResetLoc(c *gin.Context) {
	input := &Loc{}
	if err := c.ShouldBind(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	n := len(e.transactions)
	if n > 0 {
		e.transactions[n-1].city = ""
		delete(e.cityStats, input.City)
		c.JSON(http.StatusOK, gin.H{"msg": "City is reseted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Please add some transaction"})

}

func (e *Endpoint) SetLoc(c *gin.Context) {
	// Add loc to last user
	input := &Loc{}
	if err := c.ShouldBind(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	n := len(e.transactions)
	if n > 0 {
		e.transactions[n-1].city = input.City
		e.updateStats(e.transactions[n-1])
		c.JSON(http.StatusOK, gin.H{"msg": "City added to transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Please add some transaction"})

}

func (e *Endpoint) DeleteTrans(c *gin.Context) {
	e.transactions = []*Transaction{}
	e.cityStats = map[string]*Stats{}
	c.Status(http.StatusNoContent)
}

func (e *Endpoint) Statistics(c *gin.Context) {
	input := map[string]string{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := input["city"]
	if key == "" {
		key = "default"
	}
	data, ok := e.cityStats[key]
	if !ok {
		data = &Stats{}
	}
	c.JSON(http.StatusOK, data)
}

func (e *Endpoint) AddTransaction(c *gin.Context) {
	input, err := parseBody(c)
	if err != nil {
		return
	}
	//fmt.Println("input: ", input)

	parsedData, err := parseInput(input)
	if err != nil {
		fmt.Println("error: ", err.Error())
		if err.Error() == "transaction is older then 60 sec" {
			fmt.Println("t1")
			c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		}

		return
	}
	//fmt.Println(parsedData)
	e.transactions = append(e.transactions, parsedData)
	e.updateStats(parsedData)
	c.JSON(http.StatusCreated, gin.H{"msg": "Transaction added"})

}

func parseBody(c *gin.Context) (map[string]string, error) {
	input := map[string]string{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return input, err
	}
	return input, nil
}

func parseInput(input map[string]string) (*Transaction, error) {
	amountIn, ok := input["amount"]
	transaction := &Transaction{}
	var err error
	if !ok {
		return nil, errors.New("amount feild missing")
	}
	if transaction.amount, err = strconv.ParseFloat(amountIn, 64); err != nil {
		return nil, err
	}
	
	timeStampIn, ok2 := input["timeStamp"]
	if !ok2 {
		return nil, errors.New("timestamp feild missing")
	}
	transaction.timeStamp, err = time.Parse("2006-01-02T15:04:05.000Z", timeStampIn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Test", transaction)

	today := time.Now().UTC()

	if transaction.timeStamp.Before(today) {
		// in past
		diff := today.Sub(transaction.timeStamp)
		fmt.Println(diff)
		if diff.Seconds() > float64(60) {
			return nil, errors.New("transaction is older then 60 sec")
		}
	} else {
		// time stamp future date
		fmt.Println(transaction.timeStamp)
		fmt.Println(today)
		return nil, errors.New("transaction is in future date")
	}

	return transaction, nil
}

func (e *Endpoint) updateStats(tran *Transaction) {
	key := "default"

	if tran.city != "" {
		key = tran.city
	}

	s, ok := e.cityStats[key]
	if !ok {
		s = &Stats{}
	}

	s.Sum += tran.amount
	s.Count++
	s.Avg = s.Sum / s.Count
	if s.Max < tran.amount {
		s.Max = tran.amount
	}
	if tran.amount < s.Min || s.Min == 0 {
		s.Min = tran.amount
	}
	e.cityStats[key] = s
	fmt.Println("key", s)
	fmt.Println(e.cityStats)
}
