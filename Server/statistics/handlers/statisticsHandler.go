package handlers

import (
	"statistics/data"
)

type StatisticsHandler struct {
	repo              *data.StatisticsRepo
}

func NewStatisticsHandler(r *data.StatisticsRepo) *StatisticsHandler {
	return &StatisticsHandler{r}
}

// TODO Handler methods
