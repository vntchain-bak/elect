package elect

import (
	"encoding/json"
	"fmt"
)

var errNotFound = "not found"

// QueryStake returns stake information of the account in json format, or an error if failed.
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

// QueryVote returns vote information of the account in json format, or an error if failed.
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

// QueryCandidates returns a witnesses list in json format, or an error if failed.
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

// QueryRestVNTBounty returns a integer of the rest vnt bounty in wei, or an error if failed.
func (e *Election) QueryRestVNTBounty() ([]byte, error) {
	restBounty, err := e.vc.RestVNTBounty(e.ctx)
	if err != nil {
		return nil, err
	}

	return []byte(restBounty.String() + " wei"), nil
}
