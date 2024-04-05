package handlers

import "mup/data"

type MupHandler struct {
	repo *data.MupRepo
}

func NewUserHandler(r *data.MupRepo) *MupHandler {
	return &MupHandler{r}
}
