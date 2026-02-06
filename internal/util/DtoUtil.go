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

func NewCardStats(hist any, trend float64, current float64) *CardStats {
	return &CardStats{
		Hist:    hist,
		Trend:   trend,
		Current: current,
	}
}

type GraphData struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}
