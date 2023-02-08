package framework

import (
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"

	"github.com/opensourceways/robot-gitlab-lib/client"
)

const (
	logFieldOrg    = "org"
	logFieldRepo   = "repo"
	logFieldURL    = "url"
	logFieldAction = "action"

	noteableTypeIssue        = "Issue"
	noteableTypeMergeRequest = "MergeRequest"
)

// IssueEventHandler defines the handler for a gitlab.IssuesEvent.
type IssueEventHandler interface {
	HandleIssueEvent(e *gitlab.IssueEvent, log *logrus.Entry) error
}

// IssueCommentHandler defines the handler for a gitlab.IssueCommentEvent.
type IssueCommentHandler interface {
	HandleIssueCommentEvent(e *gitlab.IssueCommentEvent, log *logrus.Entry) error
}

// MergeRequestEventHandler defines the handler for a gitlab.MergeEvent on a merge request.
type MergeRequestEventHandler interface {
	HandleMergeRequestEvent(e *gitlab.MergeEvent, log *logrus.Entry) error
}

// MergeCommentEventHandler defines the handler for a gitlab.MergeCommentEvent.
type MergeCommentEventHandler interface {
	HandleMergeCommentEvent(e *gitlab.MergeCommentEvent, log *logrus.Entry) error
}

// PushEventHandler defines the handler for a gitlab.PushEvent.
type PushEventHandler interface {
	HandlePushEvent(e *gitlab.PushEvent, log *logrus.Entry) error
}

type handlers struct {
	pushEventHandler PushEventHandler

	issueEventHandler   IssueEventHandler
	issueCommentHandler IssueCommentHandler

	mergeRequestEventHandler MergeRequestEventHandler
	mergeCommentEventHandler MergeCommentEventHandler
}

func (h *handlers) registerHandler(robot interface{}) {
	if v, ok := robot.(IssueEventHandler); ok {
		h.issueEventHandler = v
	}

	if v, ok := robot.(MergeRequestEventHandler); ok {
		h.mergeRequestEventHandler = v
	}

	if v, ok := robot.(PushEventHandler); ok {
		h.pushEventHandler = v
	}

	if v, ok := robot.(IssueCommentHandler); ok {
		h.issueCommentHandler = v
	}

	if v, ok := robot.(MergeCommentEventHandler); ok {
		h.mergeCommentEventHandler = v
	}
}

func (h *handlers) getHandler() (r map[string]func([]byte, *logrus.Entry)) {
	r = make(map[string]func([]byte, *logrus.Entry))

	if h.issueEventHandler != nil {
		r[string(gitlab.EventTypeIssue)] = h.handleIssueEvent
	}

	if h.mergeRequestEventHandler != nil {
		r[string(gitlab.EventTypeMergeRequest)] = h.handleMergeRequestEvent
	}

	if h.pushEventHandler != nil {
		r[string(gitlab.EventTypePush)] = h.handlePushEvent
	}

	if h.issueCommentHandler != nil {
		r[noteableTypeIssue] = h.handleIssueCommentEvent
	}

	if h.mergeCommentEventHandler != nil {
		r[noteableTypeMergeRequest] = h.handleMergeCommentEvent
	}

	return
}

func (h *handlers) handleIssueEvent(payload []byte, l *logrus.Entry) {
	e := new(gitlab.IssueEvent)

	if err := json.Unmarshal(payload, e); err != nil {
		l.Errorf("convert to issueEvent err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.ObjectAttributes.URL,
		logFieldAction: e.ObjectAttributes.Action,
	})

	if err := h.issueEventHandler.HandleIssueEvent(e, l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (h *handlers) handleMergeRequestEvent(payload []byte, l *logrus.Entry) {
	e := new(gitlab.MergeEvent)

	if err := json.Unmarshal(payload, e); err != nil {
		l.Errorf("convert to mergeEvent err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.ObjectAttributes.URL,
		logFieldAction: e.ObjectAttributes.Action,
	})

	if err := h.mergeRequestEventHandler.HandleMergeRequestEvent(e, l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (h *handlers) handlePushEvent(payload []byte, l *logrus.Entry) {
	e := new(gitlab.PushEvent)

	if err := json.Unmarshal(payload, e); err != nil {
		l.Errorf("convert to pushEvent err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldOrg:  strings.Split(e.Project.PathWithNamespace, "/")[0],
		logFieldRepo: e.Repository.Name,
		"ref":        e.Ref,
		"head":       e.After,
	})

	if err := h.pushEventHandler.HandlePushEvent(e, l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (h *handlers) handleIssueCommentEvent(payload []byte, l *logrus.Entry) {
	e := new(gitlab.IssueCommentEvent)

	if err := json.Unmarshal(payload, e); err != nil {
		l.Errorf("convert to issueCommentEvent err: ", err.Error())

		return
	}

	l = l.WithFields(logrus.Fields{
		logFieldURL:    e.Issue.URL,
		logFieldAction: e.Issue.State,
		"commenter":    client.GetIssueCommentAuthor(e),
	})

	if err := h.issueCommentHandler.HandleIssueCommentEvent(e, l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}

func (h *handlers) handleMergeCommentEvent(payload []byte, l *logrus.Entry) {
	e := new(gitlab.MergeCommentEvent)

	if err := json.Unmarshal(payload, e); err != nil {
		l.Errorf("convert to mergeCommentEvent err: ", err.Error())

		return
	}

	org, repo := client.GetOrgRepo(e.Project.PathWithNamespace)
	l = l.WithFields(logrus.Fields{
		logFieldOrg:  org,
		logFieldRepo: repo,
		"url":        e.MergeRequest.LastCommit.URL,
		"commenter":  client.GetMRCommentAuthor(e),
	})

	if err := h.mergeCommentEventHandler.HandleMergeCommentEvent(e, l); err != nil {
		l.WithError(err).Error()
	} else {
		l.Info()
	}
}
