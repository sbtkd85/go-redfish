package domain

import (
	"context"
	"time"

	eh "github.com/looplab/eventhorizon"
)

func init() {
	eh.RegisterCommand(func() eh.Command { return &GET{} })
}

const (
	GETCommand = eh.CommandType("http:RedfishResource:GET")
)

// Static type checking for commands to prevent runtime errors due to typos
var _ = eh.Command(&GET{})

// HTTP GET Command
type GET struct {
	ID    eh.UUID `json:"id"`
	CmdID eh.UUID `json:"cmdid"`
}

func (c *GET) AggregateType() eh.AggregateType { return AggregateType }
func (c *GET) AggregateID() eh.UUID            { return c.ID }
func (c *GET) CommandType() eh.CommandType     { return GETCommand }
func (c *GET) SetAggID(id eh.UUID)             { c.ID = id }
func (c *GET) SetCmdID(id eh.UUID)             { c.CmdID = id }
func (c *GET) SetUserDetails(u string, p []string) string {
	return "checkMaster"
}
func (c *GET) Handle(ctx context.Context, a *RedfishResourceAggregate) error {
	// set up the base response data
	data := HTTPCmdProcessedData{
		CommandID:  c.CmdID,
		StatusCode: 200,
	}
	// TODO: Should be able to discern supported methods from the meta and return those

	data.Results, _ = a.ProcessMeta(ctx, "GET", map[string]interface{}{})

	// TODO: set error status code based on err from ProcessMeta
	// TODO: This is not thread safe: deep copy
	data.Headers = a.Headers

	a.eventBus.HandleEvent(ctx, eh.NewEvent(HTTPCmdProcessed, data, time.Now()))
	return nil
}
