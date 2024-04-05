package handlers

import "court/data"

type SSOHandler struct {
	repo *data.CourtRepo
}

// Constructor
func NewSSOHandler(r *data.CourtRepo) *SSOHandler {
	return &SSOHandler{r}
}

// TODO: Handler methods
