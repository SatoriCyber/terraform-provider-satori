package api

type DataAccessIdentity struct {
	IdentityType string `json:"identityType"`
	Identity     string `json:"identity"`
}

type DataAccessTimeLimit struct {
	Expiration   *interface{} `json:"expiration,omitempty"`
	ShouldExpire bool         `json:"shouldExpire"`
}

type DataAccessUnusedTimeLimit struct {
	UnusedDaysUntilRevocation int  `json:"unusedDaysUntilRevocation"`
	ShouldRevoke              bool `json:"shouldRevoke"`
}

type DataAccessSelfServiceAndRequestTimeLimit struct {
	ShouldExpire bool   `json:"shouldExpire"`
	UnitType     string `json:"unitType"`
	Units        int    `json:"units"`
}
