package main

import (
	"github.com/sirupsen/logrus"
	sdk "github.com/xanzy/go-gitlab"
)

// TODO: set botName
const botName = ""

type iClient interface {
}

func newRobot(cli iClient, gc func() (*configuration, error)) *robot {
	return &robot{cli: cli, getConfig: gc}
}

type robot struct {
	getConfig func() (*configuration, error)
	cli       iClient
}

func (bot *robot) HandleMergeRequestEvent(e *sdk.MergeEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to handle PR event, delete this function.
	return nil
}

func (bot *robot) HandleIssueEvent(e *sdk.IssueEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to handle Issue event, delete this function.
	return nil
}

func (bot *robot) HandlePushEvent(e *sdk.PushEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to handle Push event, delete this function.
	return nil
}

func (bot *robot) HandleMergeCommentEvent(e *sdk.MergeCommentEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to handle Note event of PR, delete this function.
	return nil
}

func (bot *robot) HandleIssueCommentEvent(e *sdk.IssueCommentEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to handle Note event of Issue, delete this function.
	return nil
}
