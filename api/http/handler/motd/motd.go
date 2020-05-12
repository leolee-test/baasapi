package motd

import (
	"net/http"

	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/crypto"
	"github.com/baasapi/baasapi/api/http/client"
)

type motdResponse struct {
	Title   string `json:"Title"`
	Message string `json:"Message"`
	Hash    []byte `json:"Hash"`
}

func (handler *Handler) motd(w http.ResponseWriter, r *http.Request) {

	motd, err := client.Get(baasapi.MessageOfTheDayURL, 0)
	if err != nil {
		response.JSON(w, &motdResponse{Message: ""})
		return
	}

	title, err := client.Get(baasapi.MessageOfTheDayTitleURL, 0)
	if err != nil {
		response.JSON(w, &motdResponse{Message: ""})
		return
	}

	hash := crypto.HashFromBytes(motd)
	response.JSON(w, &motdResponse{Title: string(title), Message: string(motd), Hash: hash})
}
