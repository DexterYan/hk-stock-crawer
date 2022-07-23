package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	t.Run("Successful Get Stock List", func(t *testing.T) {
		stockList := getStockList()
		assert.Greater(t, len(stockList), 0, "The Stock List should above 0")
	})

	t.Run("Successful Get Stock Summary", func(t *testing.T) {
		stockCode := "00700"
		stock := getStockCurrentSummary(stockCode)
		assert.NotEmpty(t, stock.NameEn, "The stock name should not be empty")
	})
}
