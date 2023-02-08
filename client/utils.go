package client

import "strings"

const (
	ActionOpened  = "opened"
	ActionCreated = "created"
	ActionReopen  = "reopened"
	ActionClosed  = "closed"
)

type ErrorForbidden struct {
	err string
}

func (e ErrorForbidden) Error() string {
	return e.err
}

// GetOrgRepo return the owner and name of the project
func GetOrgRepo(orgPath string) (string, string) {
	org, repo := strings.Split(orgPath, "/")[0], strings.Split(orgPath, "/")[1]
	return org, repo
}
