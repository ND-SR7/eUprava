package handlers

import "sso/data"

type SSOHandler struct {
	repo *data.SSORepo
}

// Constructor
func NewSSOHandler(r *data.SSORepo) *SSOHandler {
	return &SSOHandler{r}
}

// TODO: Handler methods
