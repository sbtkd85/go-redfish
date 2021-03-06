package stdcollections

import (
	"context"
	"time"

	domain "github.com/superchalupa/go-redfish/src/redfishresource"

	plugins "github.com/superchalupa/go-redfish/src/ocp"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/utils"
)

// wait in a listener for the root service to be created, then extend it
func InitService(ctx context.Context, ch eh.CommandHandler, eb eh.EventBus, ew *utils.EventWaiter) {
	sp, err := plugins.NewEventStreamProcessor(ctx, ew, plugins.SelectEventResourceCreatedByURI("/redfish/v1"))
	if err == nil {
		sp.RunOnce(func(event eh.Event) {
			rootID := event.Data().(domain.RedfishResourceCreatedData).ID
			NewService(ctx, rootID, ch)
		})
	}
}

func NewService(ctx context.Context, rootID eh.UUID, ch eh.CommandHandler) {
	// Create Computer System Collection
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: true,

			ResourceURI: "/redfish/v1/Systems",
			Type:        "#ComputerSystemCollection.ComputerSystemCollection",
			Context:     "/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name": "Computer System Collection",
			}})

	ch.HandleCommand(ctx,
		&domain.UpdateRedfishResourceProperties{
			ID: rootID,
			Properties: map[string]interface{}{
				"Systems": map[string]interface{}{"@odata.id": "/redfish/v1/Systems"},
			},
		})

	// Create Computer System Collection
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: true,

			ResourceURI: "/redfish/v1/Chassis",
			Type:        "#ChassisCollection.ChassisCollection",
			Context:     "/redfish/v1/$metadata#ChassisCollection.ChassisCollection",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name": "Chassis Collection",
			}})

	ch.HandleCommand(ctx,
		&domain.UpdateRedfishResourceProperties{
			ID: rootID,
			Properties: map[string]interface{}{
				"Chassis": map[string]interface{}{"@odata.id": "/redfish/v1/Chassis"},
			},
		})

	// Create Computer System Collection
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: true,

			ResourceURI: "/redfish/v1/Managers",
			Type:        "#ManagerCollection.ManagerCollection",
			Context:     "/redfish/v1/$metadata#ManagerCollection.ManagerCollection",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name": "Manager Collection",
			}})

	ch.HandleCommand(ctx,
		&domain.UpdateRedfishResourceProperties{
			ID: rootID,
			Properties: map[string]interface{}{
				"Managers": map[string]interface{}{"@odata.id": "/redfish/v1/Managers"},
			},
		})

	// Add Accounts collection
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: true,

			ResourceURI: "/redfish/v1/AccountService/Accounts",
			Type:        "#ManagerAccountCollection.ManagerAccountCollection",
			Context:     "/redfish/v1/$metadata#ManagerAccountCollection.ManagerAccountCollection",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name": "Accounts Collection",
			}})

	// Add Roles collection
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: true,

			ResourceURI: "/redfish/v1/AccountService/Roles",
			Type:        "#RoleCollection.RoleCollection",
			Context:     "/redfish/v1/$metadata#RoleCollection.RoleCollection",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name": "Roles Collection",
			}})

	// Create Computer System Collection
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: false,

			ResourceURI: "/redfish/v1/AccountService",
			Type:        "#AccountService.v1_0_2.AccountService",
			Context:     "/redfish/v1/$metadata#AccountService.AccountService",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{"ConfigureManager"}, // cannot create sub objects
				"PUT":    []string{"ConfigureManager"},
				"PATCH":  []string{"ConfigureManager"},
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Id":          "AccountService",
				"Name":        "Account Service",
				"Description": "Account Service",
				"Status": map[string]interface{}{
					"State":  "Enabled",
					"Health": "OK",
				},
				"ServiceEnabled":                  true,
				"AuthFailureLoggingThreshold":     3,
				"MinPasswordLength":               8,
				"AccountLockoutThreshold":         5,
				"AccountLockoutDuration":          30,
				"AccountLockoutCounterResetAfter": 30,
				"Accounts":                        map[string]string{"@odata.id": "/redfish/v1/AccountService/Accounts"},
				"Roles":                           map[string]string{"@odata.id": "/redfish/v1/AccountService/Roles"},
			}})

	ch.HandleCommand(ctx,
		&domain.UpdateRedfishResourceProperties{
			ID: rootID,
			Properties: map[string]interface{}{
				"AccountService": map[string]interface{}{"@odata.id": "/redfish/v1/AccountService"},
			},
		})

	// give time between running command to create role collection and when we add roles
	time.Sleep(10 * time.Millisecond)

	// add standard DMTF roles: Admin
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: false,

			ResourceURI: "/redfish/v1/AccountService/Roles/Admin",
			Type:        "#Role.v1_0_2.Role",
			Context:     "/redfish/v1/$metadata#Role.Role",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name":         "User Role",
				"Id":           "Admin",
				"Description":  "Admin User Role",
				"IsPredefined": true,
				"AssignedPrivileges": []string{
					"Login",
					"ConfigureManager",
					"ConfigureUsers",
					"ConfigureSelf",
					"ConfigureComponents",
				},
			}})

	// add standard DMTF roles: Operator
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: false,

			ResourceURI: "/redfish/v1/AccountService/Roles/Operator",
			Type:        "#Role.v1_0_2.Role",
			Context:     "/redfish/v1/$metadata#Role.Role",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name":         "User Role",
				"Id":           "Operator",
				"Description":  "Operator User Role",
				"IsPredefined": true,
				"AssignedPrivileges": []string{
					"Login",
					"ConfigureSelf",
					"ConfigureComponents",
				},
			}})

	// add standard DMTF roles: Read-only
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:         eh.NewUUID(),
			Collection: false,

			ResourceURI: "/redfish/v1/AccountService/Roles/ReadOnlyUser",
			Type:        "#Role.v1_0_2.Role",
			Context:     "/redfish/v1/$metadata#Role.Role",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // Read Only
				"PUT":    []string{}, // Read Only
				"PATCH":  []string{}, // Read Only
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"Name":         "User Role",
				"Id":           "ReadOnlyUser",
				"Description":  "ReadOnlyUser User Role",
				"IsPredefined": true,
				"AssignedPrivileges": []string{
					"Login",
				},
			}})

}
