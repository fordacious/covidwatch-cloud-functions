package app

import (
	"log"
	"net/http"

	"app/internal/pow"
	"app/internal/util"
)

// reportBody is the body of a POST request to the /report endpoint.
type reportBody struct {
	Challenge pow.ChallengeSolution `json:"challenge"`
	Report    report                `json:"report"`
}

func (r *reportBody) ChallengeSolution() pow.ChallengeSolution {
	return r.Challenge
}

type report struct {
	Data []byte `json:"data"`
}

// ReportHandler is a handler for the /report endpoint.
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := util.NewContext(w, r)
	if err != nil {
		log.Print(err)
		return
	}

	if err = util.ValidateRequestMethod(&ctx, "POST", ""); err != nil {
		log.Print(err)
		return
	}

	var report reportBody
	if err := pow.UnmarshalAndValidateSolution(&ctx, &report); err != nil {
		log.Print(err)
		ctx.WriteStatusError(err)
		return
	}
}
