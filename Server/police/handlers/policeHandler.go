package handlers

import "police/data"

type PoliceHandler struct {
	repo *data.PoliceRepo
}

// Constructor
func NewPoliceHandler(r *data.PoliceRepo) *PoliceHandler {
	return &PoliceHandler{r}
}

// TODO: Handler methods
