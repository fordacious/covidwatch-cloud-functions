package app

import (
	"encoding/json"
	"log"
	"net/http"

	"app/internal/pow"
	"app/internal/util"
)

// ChallengeHandler is a handler for the /challenge endpoint.
func ChallengeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := util.NewContext(w, r)
	if err != nil {
		log.Print(err)
		return
	}

	if err = util.ValidateRequestMethod(&ctx, "GET", ""); err != nil {
		log.Print(err)
		return
	}

	c, err := pow.GenerateChallenge(ctx)
	if err != nil {
		ctx.WriteInternalServerError()
		log.Print(err)
		return
	}
	json.NewEncoder(w).Encode(c)
}
