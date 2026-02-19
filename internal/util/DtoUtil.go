package util

import "time"

type Page[E any] struct {
	HasNextPage bool   `json:"hasNextPage"`
	NextCursor  string `json:"nextCursor"`
	List        *[]E   `json:"list"`
}

func NewPage[E any](list []E, cursor string, limit int) *Page[E] {

	if len(list) < limit {
		page := Page[E]{
			List:        &list,
			NextCursor:  "",
			HasNextPage: false,
		}
		return &page
	}

	page := Page[E]{
		List:        &list,
		NextCursor:  cursor,
		HasNextPage: true,
	}

	return &page
}

type CardStats struct {
	Hist    any     `json:"hist"`
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
}

func NewCardPercentage(hist []GraphData) *CardStats {
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[0].Value
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].Value
		previous = hist[lenght-2].Value
		trend = CalculatePercentageTrend(current, previous)
	}

	card := &CardStats{
		Current: current,
		Trend:   trend,
		Hist:    hist,
	}

	return card
}

func NewCardInt(hist []GraphData) *CardStats {
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[0].Value
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].Value
		previous = hist[lenght-2].Value
		trend = current - previous
	}

	card := &CardStats{
		Current: current,
		Trend:   trend,
		Hist:    hist,
	}

	return card
}

type GraphData struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}
