// Build tags: only build this for the openbmc build. Be sure to note the required blank line after.
// +build openbmc

package obmc

import (
	"context"
	"fmt"
	"math"
	"path"
	"time"

	"github.com/godbus/dbus"
	mapper "github.com/superchalupa/go-redfish/plugins/dbus"
	"github.com/superchalupa/go-redfish/plugins/ocp/thermal/temperatures"
)

type Optioner interface {
	ApplyOption(options ...interface{}) error
}

const (
	SensorValue     = "xyz.openbmc_project.Sensor.Value"
	SensorThreshold = "xyz.openbmc_project.Sensor.Threshold.Warning"
)

func UpdateSensorList(ctx context.Context, temps Optioner) {
	var conn *dbus.Conn
	var err error
	for {
		// do{}while(0) equivalent so that we can skip the rest on errors
		// (break), but still hit the outside loop end to check for context
		// cancellation and do our sleep.
		for ok := true; ok; ok = false {
			if conn == nil {
				conn, err = dbus.SystemBus()
				if err != nil {
					fmt.Printf("Could not connect to system bus: %s\n", err)
					break
				}
			}
			m := mapper.New(conn)
			ret, err := m.GetSubTree(ctx, "/xyz/openbmc_project/sensors/temperature", 0, "xyz.openbmc_project.Sensor.Value")
			if err != nil {
				fmt.Printf("Mapper call failed: %s\n", err.Error())
				break
			}
			if len(ret) == 0 {
				fmt.Printf("empty array?")
				break
			}
			arr_0 := ret[0]
			dict, ok := arr_0.(map[string]map[string][]string)
			if !ok {
				fmt.Printf("type assert failed: %T\n", arr_0)
				break
			}

			for p, m1 := range dict {
				for bus, _ := range m1 {
					fmt.Printf("getting thermal for bus(%s)  path(%s)\n", bus, p)
					temps.ApplyOption(
						temperatures.WithSensor(
							fmt.Sprintf("%s#%s", bus, p),
							getThermal(ctx, conn, bus, p)))
				}
			}
		}

		// sleep for 10 seconds, or until context is cancelled
		select {
		case <-ctx.Done():
			fmt.Printf("Cancelling UpdateSensorList due to context cancellation.\n")
			break
		case <-time.After(10 * time.Second):
			continue
		}
	}
}

func getThermal(ctx context.Context, conn *dbus.Conn, bus string, objectPath string) *temperatures.RedfishThermalSensor {
	busObject := conn.Object(bus, dbus.ObjectPath(objectPath))

	unit, err := busObject.GetProperty(SensorValue + ".Unit")
	if err != nil {
		fmt.Printf("Error getting .Unit property for bus(%s) path(%s): %s\n", bus, objectPath, err.Error())
		return nil
	}
	if unit.Value() != "xyz.openbmc_project.Sensor.Value.Unit.DegreesC" {
		fmt.Printf("Don't know how to handle units: %s\n", unit)
		return nil
	}

	scale, err := busObject.GetProperty(SensorValue + ".Scale")
	if err != nil {
		fmt.Printf("Error getting .Scale property for bus(%s) path(%s): %s\n", bus, objectPath, err.Error())
		return nil
	}
	s, ok := scale.Value().(int64)
	if !ok {
		fmt.Printf("Type assert of scale to int failed.\n")
		return nil
	}

	value, err := busObject.GetProperty(SensorValue + ".Value")
	if err != nil {
		fmt.Printf("Error getting .Value property for bus(%s) path(%s): %s\n", bus, objectPath, err.Error())
		return nil
	}
	v, ok := value.Value().(int64)
	if !ok {
		fmt.Printf("Type assert of value to int failed: %T\n", value.Value())
		return nil
	}

	// BOOL that says if we raise alarm if it goes above AlarmHigh
	_, err = busObject.GetProperty(SensorThreshold + ".WarningAlarmHigh")
	if err != nil {
		fmt.Printf("Error getting .WarningAlarmHigh property for bus(%s) path(%s): %s\n", bus, objectPath, err.Error())
	}

	UpperCriticalV, err := busObject.GetProperty(SensorThreshold + ".WarningHigh")
	if err != nil {
		fmt.Printf("Error getting .WarningHigh property for bus(%s) path(%s): %s\n", bus, objectPath, err.Error())
	}
	UpperCritical, ok := UpperCriticalV.Value().(int64)

	var scaleMultiplier float64 = math.Pow(10, float64(s))
	return &temperatures.RedfishThermalSensor{
		Name:                      path.Base(objectPath),
		ReadingCelsius:            float64(v) * scaleMultiplier,
		UpperThresholdNonCritical: float64(UpperCritical) * scaleMultiplier,
		UpperThresholdCritical:    float64(UpperCritical) * scaleMultiplier,
		UpperThresholdFatal:       float64(UpperCritical) * scaleMultiplier,
		//Status:                    StdStatus{State: "Enabled", Health: "OK"},
	}
}
