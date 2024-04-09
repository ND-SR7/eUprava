package handlers

import "police/data"

type SSOHandler struct {
	repo *data.PoliceRepo
}

// Constructor
func NewSSOHandler(r *data.PoliceRepo) *SSOHandler {
	return &SSOHandler{r}
}

// TODO: Handler methods
