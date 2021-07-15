package api

import "time"

type DataAccessIdentity struct {
	IdentityType string `json:"identityType"`
	Identity     string `json:"identity"`
}

type DataAccessTimeLimit struct {
	Expiration   *time.Time `json:"expiration,omitempty"`
	ShouldExpire bool       `json:"shouldExpire"`
}

type DataAccessUnusedTimeLimit struct {
	UnusedDaysUntilRevocation int32 `json:"unusedDaysUntilRevocation"`
	ShouldRevoke              bool  `json:"shouldRevoke"`
}
