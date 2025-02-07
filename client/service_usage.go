package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type ServiceUsageClient commonClient

// ServiceUsageOptions list filters
type ServiceUsageOptions struct {
	*ListOptions

	AfterGUID            string `qs:"after_guid"`
	GUIDs                Filter `qs:"guids"`
	ServiceInstanceTypes Filter `qs:"service_instance_types"`
	ServiceOfferingGUIDs Filter `qs:"service_offering_guids"`
}

// NewServiceUsageOptions creates new options to pass to list
func NewServiceUsageOptions() *ServiceUsageOptions {
	return &ServiceUsageOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceUsageOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Get retrieves the specified service event
func (c *ServiceUsageClient) Get(guid string) (*resource.ServiceUsage, error) {
	var a resource.ServiceUsage
	err := c.client.get(path.Format("/v3/service_usage_events/%s", guid), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// List pages all service usage events
func (c *ServiceUsageClient) List(opts *ServiceUsageOptions) ([]*resource.ServiceUsage, *Pager, error) {
	if opts == nil {
		opts = NewServiceUsageOptions()
	}
	var res resource.ServiceUsageList
	err := c.client.get(path.Format("/v3/service_usage_events?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service usage events
func (c *ServiceUsageClient) ListAll(opts *ServiceUsageOptions) ([]*resource.ServiceUsage, error) {
	if opts == nil {
		opts = NewServiceUsageOptions()
	}
	return AutoPage[*ServiceUsageOptions, *resource.ServiceUsage](opts, func(opts *ServiceUsageOptions) ([]*resource.ServiceUsage, *Pager, error) {
		return c.List(opts)
	})
}

// Purge destroys all existing events. Populates new usage events, one for each existing service instance.
// All populated events will have a created_at value of current time.
//
// There is the potential race condition if service instances are currently being created or deleted.
// The seeded usage events will have the same guid as the service instance.
func (c *ServiceUsageClient) Purge() error {
	_, err := c.client.post("/v3/service_usage_events/actions/destructively_purge_all_and_reseed", nil, nil)
	return err
}
