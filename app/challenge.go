package app

import (
	"encoding/json"
	"net/http"

	"app/internal/pow"
	"app/internal/util"
)

// ChallengeHandler is a handler for the /challenge endpoint.
func ChallengeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := util.NewContext(w, r)
	if err != nil {
		return
	}

	if err = util.ValidateRequestMethod(&ctx, "GET", ""); err != nil {
		return
	}

	c, err := pow.GenerateChallenge(ctx)
	if err != nil {
		ctx.WriteInternalServerError()
		return
	}
	json.NewEncoder(w).Encode(c)
}
