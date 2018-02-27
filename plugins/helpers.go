package plugins

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	domain "github.com/superchalupa/go-redfish/redfishresource"
)

type Option func(*Service) error

type Service struct {
	sync.Mutex
	pluginType domain.PluginType
}

func NewService(options ...Option) *Service {
	s := &Service{}
	s.ApplyOption(options...)
	return s
}

func (c *Service) ApplyOption(options ...Option) error {
	for _, o := range options {
		err := o(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func PluginType(pt domain.PluginType) Option {
	return func(s *Service) error {
		s.pluginType = pt
		return nil
	}
}

func (s *Service) PluginType() domain.PluginType { return s.pluginType }

func RefreshProperty(
	ctx context.Context,
	s interface{},
	rrp *domain.RedfishResourceProperty,
	method string,
	meta map[string]interface{},
) {
	property, ok := meta["property"].(string)
	if ok {
		v := reflect.ValueOf(s)
		for i := 0; i < v.NumField(); i++ {
			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			tag := v.Type().Field(i).Tag.Get("property")
			if tag == property {
				rrp.Value = v.Field(i).Interface()
				return
			}
		}
	}
	fmt.Printf("Incorrect metadata in aggregate: neither 'data' nor 'property' set to something handleable")
}