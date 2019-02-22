package elect

import (
	"encoding/json"
)

func (e *Election) QueryStake() ([]byte, error) {
	stake, err := e.vc.StakeAt(e.ctx, e.cfg.Sender)
	if err != nil {
		return nil, err
	}

	return json.Marshal(stake)
}

func (e *Election) QueryVote() ([]byte, error) {
	voter, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil {
		return nil, err
	}

	return json.Marshal(voter)
}

func (e *Election) QueryCandidates() ([]byte, error) {
	candidates, err := e.vc.WitnessCandidates(e.ctx)
	if err != nil {
		return nil, err
	}

	return json.Marshal(candidates)
}

func (e *Election) QueryRestVNTBounty() ([]byte, error) {
	restBounty, err := e.vc.RestVNTBounty(e.ctx)
	if err != nil {
		return nil, err
	}

	return []byte(restBounty.String() + " wei"), nil
}
