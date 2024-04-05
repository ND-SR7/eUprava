package handlers

import "mup/data"

type MupHandler struct {
	repo *data.MUPRepo
}

func NewUserHandler(r *data.MUPRepo) *MupHandler {
	return &MupHandler{r}
}
