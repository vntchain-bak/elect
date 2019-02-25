package elect

import (
	"encoding/json"
	"fmt"
)

var errNotFound = "not found"

func (e *Election) QueryStake() ([]byte, error) {
	stake, err := e.vc.StakeAt(e.ctx, e.cfg.Sender)
	if err != nil {
		if err.Error() == errNotFound {
			return nil, fmt.Errorf("no stake information for account: %s", e.cfg.Sender.String())
		}
		return nil, err
	}

	return json.Marshal(stake)
}

func (e *Election) QueryVote() ([]byte, error) {
	voter, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil {
		if err.Error() == errNotFound {
			return nil, fmt.Errorf("no vote information for account: %s", e.cfg.Sender.String())
		}
		return nil, err
	}

	return json.Marshal(voter)
}

func (e *Election) QueryCandidates() ([]byte, error) {
	candidates, err := e.vc.WitnessCandidates(e.ctx)
	if err != nil {
		if err.Error() == errNotFound {
			return nil, fmt.Errorf("witness candidate list is empty")
		}
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
