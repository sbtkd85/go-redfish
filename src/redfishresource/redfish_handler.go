package domain

import (
	"context"
	"encoding/json"
	"net/http"

	eh "github.com/looplab/eventhorizon"
)

type CmdIDSetter interface {
	SetCmdID(eh.UUID)
}
type AggIDSetter interface {
	SetAggID(eh.UUID)
}
type UserDetailsSetter interface {
	SetUserDetails(string, []string) string
	// return codes:
	//      "checkMaster" - command check passed, but also check master
	//      "authorized"  - command check passed, go right ahead, no master check
	//      "unauthorized" - failed auth check in command, no need to check further

}
type HTTPParser interface {
	ParseHTTPRequest(*http.Request) error
}

func NewRedfishHandler(dobjs *DomainObjects, u string, p []string) *RedfishHandler {
	return &RedfishHandler{UserName: u, Privileges: p, d: dobjs}
}

type RedfishHandler struct {
	UserName   string
	Privileges []string
	d          *DomainObjects
}

func (rh *RedfishHandler) IsAuthorized(requiredPrivs []string) (authorized string) {
	authorized = "unauthorized"
	if requiredPrivs == nil {
		requiredPrivs = []string{}
	}
outer:
	for _, p := range rh.Privileges {
		for _, q := range requiredPrivs {
			if p == q {
				authorized = "authorized"
				break outer
			}
		}
	}
	return
}

// TODO: need to write middleware that would allow different types of encoding on output
func (rh *RedfishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All operations have to be on URLs that exist, so look it up in the tree
	aggID, ok := rh.d.GetAggregateIDOK(r.URL.Path)
	if !ok {
		http.Error(w, "Could not find URL: "+r.URL.Path, http.StatusNotFound)
		return
	}

	search := []eh.CommandType{}

	// load the aggregate for the URL we are operating on
	agg, err := rh.d.AggregateStore.Load(r.Context(), AggregateType, aggID)
	// type assertion to get real aggregate
	redfishResource, ok := agg.(*RedfishResourceAggregate)
	if ok {
		// prepend the plugins to the search path
		search = append(search, eh.CommandType(redfishResource.ResourceURI+":"+r.Method))
		search = append(search, eh.CommandType(redfishResource.GetProperty("@odata.type").(string)+":"+r.Method))
		search = append(search, eh.CommandType(redfishResource.GetProperty("@odata.context").(string)+":"+r.Method))
		search = append(search, eh.CommandType(redfishResource.Plugin+":"+r.Method))
	}
	search = append(search, eh.CommandType("http:RedfishResource:"+r.Method))

	// search through the commands until we find one that exists
	var cmd eh.Command
	for _, cmdType := range search {
		cmd, err = eh.CreateCommand(cmdType)
		if err == nil {
			break
		}
	}

	// with a proper error if we couldnt create a command of any kind
	if cmd == nil {
		http.Error(w, "could not create command", http.StatusBadRequest)
		return
	}

	// Each command needs a unique UUID. We'll use that to listen for the HTTPProcessed Event, which should have a matching UUID.
	cmdId := eh.NewUUID()

	// some optional interfaces that the commands might implement
	if t, ok := cmd.(CmdIDSetter); ok {
		t.SetCmdID(cmdId)
	}
	if t, ok := cmd.(AggIDSetter); ok {
		t.SetAggID(aggID)
	}

	// Choices: command can process Authorization, or we can process authorization, or both
	// If command implements UserDetailsSetter interface, we'll go ahead and call that.
	// Return code from command determines if we also check privs here
	authAction := "checkMaster"
	var implementsAuthorization bool
	if t, implementsAuthorization := cmd.(UserDetailsSetter); implementsAuthorization {
		authAction = t.SetUserDetails(rh.UserName, rh.Privileges)
	}
	// if command does not implement userdetails setter, we always check privs here
	if !implementsAuthorization || authAction == "checkMaster" {
		privsToCheck := redfishResource.PrivilegeMap[r.Method]

		// convert Privileges from []interface{} to []string (way more code than there should be for something this simple)
		var t []string
		switch privs := privsToCheck.(type) {
		case []string:
			t = append(t, privs...)
		case []interface{}:
			for _, v := range privs {
				if a, ok := v.(string); ok {
					t = append(t, a)
				}
			}
		default:
		}

		authAction = rh.IsAuthorized(t)
	}

	if authAction != "authorized" {
		http.Error(w, "Not authorized to access this resource: ", http.StatusUnauthorized)
		return
	}

	// to avoid races, set up our listener first
	l, err := rh.d.EventWaiter.Listen(r.Context(), func(event eh.Event) bool {
		if event.EventType() != HTTPCmdProcessed {
			return false
		}
		if data, ok := event.Data().(HTTPCmdProcessedData); ok {
			if data.CommandID == cmdId {
				return true
			}
		}
		return false
	})
	if err != nil {
		http.Error(w, "could not create waiter"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer l.Close()

	// don't run parse until after privilege checks have been done
	if t, ok := cmd.(HTTPParser); ok {
		err := t.ParseHTTPRequest(r)
		if err != nil {
			http.Error(w, "Problems parsing http request: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	ctx := context.Background()
	if err := rh.d.CommandHandler.HandleCommand(ctx, cmd); err != nil {
		http.Error(w, "redfish handler could not handle command (type: "+string(cmd.CommandType())+"): "+err.Error(), http.StatusBadRequest)
		return
	}

	event, err := l.Wait(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, ok := event.Data().(HTTPCmdProcessedData)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set headers first
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains") // for A+ SSL Labs score
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("OData-Version", "4.0")
	w.Header().Set("Server", "go-redfish")

	for k, v := range data.Headers {
		w.Header().Add(k, v)
	}

	// and then encode response
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data.Results)
	return
}
