package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Runs handles communication with the run related methods of the Terraform
// Enterprise API.
//
// TFE API docs: https://www.terraform.io/docs/enterprise/api/run.html
type Runs struct {
	client *Client
}

// RunStatus represents a run state.
type RunStatus string

//List all available run statuses.
const (
	RunApplied        RunStatus = "applied"
	RunApplying       RunStatus = "applying"
	RunCanceled       RunStatus = "canceled"
	RunConfirmed      RunStatus = "confirmed"
	RunDiscarded      RunStatus = "discarded"
	RunErrored        RunStatus = "errored"
	RunPending        RunStatus = "pending"
	RunPlanned        RunStatus = "planned"
	RunPlanning       RunStatus = "planning"
	RunPolicyChecked  RunStatus = "policy_checked"
	RunPolicyChecking RunStatus = "policy_checking"
	RunPolicyOverride RunStatus = "policy_override"
)

// RunSource represents a source type of a run.
type RunSource string

// List all available run sources.
const (
	RunSourceAPI                  RunSource = "tfe-api"
	RunSourceConfigurationVersion RunSource = "tfe-configuration-version"
	RunSourceUI                   RunSource = "tfe-ui"
)

// Run represents a Terraform Enterprise run.
type Run struct {
	ID               string               `jsonapi:"primary,runs"`
	Actions          *RunActions          `jsonapi:"attr,actions"`
	CreatedAt        time.Time            `jsonapi:"attr,created-at,iso8601"`
	HasChanges       bool                 `jsonapi:"attr,has-changes"`
	IsDestroy        bool                 `jsonapi:"attr,is-destroy"`
	Message          string               `jsonapi:"attr,message"`
	Permissions      *RunPermissions      `jsonapi:"attr,permissions"`
	Source           RunSource            `jsonapi:"attr,source"`
	Status           RunStatus            `jsonapi:"attr,status"`
	StatusTimestamps *RunStatusTimestamps `jsonapi:"attr,status-timestamps"`

	// Relations
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`
	Workspace            *Workspace            `jsonapi:"relation,workspace"`
}

// RunActions represents the run actions.
type RunActions struct {
	IsCancelable  bool `json:"is-cancelable"`
	IsComfirmable bool `json:"is-comfirmable"`
	IsDiscardable bool `json:"is-discardable"`
}

// RunPermissions represents the run permissions.
type RunPermissions struct {
	CanApply        bool `json:"can-apply"`
	CanCancel       bool `json:"can-cancel"`
	CanDiscard      bool `json:"can-discard"`
	CanForceExecute bool `json:"can-force-execute"`
}

// RunStatusTimestamps holds the timestamps for individual run statuses.
type RunStatusTimestamps struct {
	ErroredAt  time.Time `json:"errored-at"`
	FinishedAt time.Time `json:"finished-at"`
	QueuedAt   time.Time `json:"queued-at"`
	StartedAt  time.Time `json:"started-at"`
}

// RunListOptions represents the options for listing runs.
type RunListOptions struct {
	ListOptions
}

// List runs of the given workspace.
func (s *Runs) List(ctx context.Context, workspaceID string, options RunListOptions) ([]*Run, error) {
	if !validStringID(&workspaceID) {
		return nil, errors.New("Invalid value for workspace ID")
	}

	u := fmt.Sprintf("workspaces/%s/runs", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	result, err := s.client.do(ctx, req, []*Run{})
	if err != nil {
		return nil, err
	}

	var rs []*Run
	for _, r := range result.([]interface{}) {
		rs = append(rs, r.(*Run))
	}

	return rs, nil
}

// RunCreateOptions represents the options for creating a new run.
type RunCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,runs"`

	// Specifies if this plan is a destroy plan, which will destroy all
	// provisioned resources.
	IsDestroy *bool `jsonapi:"attr,is-destroy,omitempty"`

	// Specifies the message to be associated with this run.
	Message *string `jsonapi:"attr,message,omitempty"`

	// Specifies the configuration version to use for this run. If the
	// configuration version object is omitted, the run will be created using the
	// workspace's latest configuration version.
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`

	// Specifies the workspace where the run will be executed.
	Workspace *Workspace `jsonapi:"relation,workspace"`
}

func (o RunCreateOptions) valid() error {
	if o.Workspace == nil {
		return errors.New("Workspace is required")
	}
	return nil
}

// Create is used to create a new run.
func (s *Runs) Create(ctx context.Context, options RunCreateOptions) (*Run, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "runs", &options)
	if err != nil {
		return nil, err
	}

	r, err := s.client.do(ctx, req, &Run{})
	if err != nil {
		return nil, err
	}

	return r.(*Run), nil
}

// Read a single run by its ID.
func (s *Runs) Read(ctx context.Context, runID string) (*Run, error) {
	if !validStringID(&runID) {
		return nil, errors.New("Invalid value for run ID")
	}

	u := fmt.Sprintf("runs/%s", url.QueryEscape(runID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	r, err := s.client.do(ctx, req, &Run{})
	if err != nil {
		return nil, err
	}

	return r.(*Run), nil
}

// RunApplyOptions represents the options for applying a run.
type RunApplyOptions struct {
	// An optional comment about the run.
	Comment *string `json:"comment,omitempty"`
}

// Apply a specific run by its ID.
func (s *Runs) Apply(ctx context.Context, runID string, options RunApplyOptions) error {
	if !validStringID(&runID) {
		return errors.New("Invalid value for run ID")
	}

	u := fmt.Sprintf("runs/%s/actions/apply", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)

	return err
}

// RunCancelOptions represents the options for canceling a run.
type RunCancelOptions struct {
	// An optional explanation for why the run was canceled.
	Comment *string `json:"comment,omitempty"`
}

// Cancel a specific run by its ID.
func (s *Runs) Cancel(ctx context.Context, runID string, options RunCancelOptions) error {
	if !validStringID(&runID) {
		return errors.New("Invalid value for run ID")
	}

	u := fmt.Sprintf("runs/%s/actions/cancel", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)

	return err
}

// RunDiscardOptions represents the options for discarding a run.
type RunDiscardOptions struct {
	// An optional explanation for why the run was discarded.
	Comment *string `json:"comment,omitempty"`
}

// Discard a specific run by its ID.
func (s *Runs) Discard(ctx context.Context, runID string, options RunDiscardOptions) error {
	if !validStringID(&runID) {
		return errors.New("Invalid value for run ID")
	}

	u := fmt.Sprintf("runs/%s/actions/discard", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)

	return err
}
