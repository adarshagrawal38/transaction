package controller

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateStats(t *testing.T) {
	endPoint := NewEndpoint()
	tran := &Transaction{amount: 10.0}
	endPoint.updateStats(tran)

	tran1 := &Transaction{amount: 20.0}
	endPoint.updateStats(tran1)

	values := endPoint.cityStats["default"]

	assert.Equal(t, 15.0, values.Avg)
	assert.Equal(t, 10.0, values.Min)
	assert.Equal(t, 20.0, values.Max)
	assert.Equal(t, 2.0, values.Count)
	assert.Equal(t, 30.0, values.Sum)
}

func TestParseInput(t *testing.T) {
	input := map[string]string{}

	trans, err := parseInput(input)
	assert.Error(t, errors.New("amount feild missing"), err)
	assert.Nil(t, trans)

	input["amount"] = "abcd"
	trans, err = parseInput(input)
	assert.NotNil(t, err)
	assert.Nil(t, trans)

	input["amount"] = "123.0"
	trans, err = parseInput(input)
	assert.Error(t, errors.New("timestamp feild missing"), err)
	assert.Nil(t, trans)

	input["timeStamp"] = "2023-01-23T17:00:00.000Z"
	trans, err = parseInput(input)
	assert.Error(t, errors.New("transaction is older then 60 sec"), err)
	assert.Nil(t, trans)

}
