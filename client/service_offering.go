package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServiceOfferingClient commonClient

// ServiceOfferingListOptions list filters
type ServiceOfferingListOptions struct {
	*ListOptions

	Names              Filter `qs:"names"`
	ServiceBrokerGUIDs Filter `qs:"service_broker_guids"`
	ServiceBrokerNames Filter `qs:"service_broker_names"`
	SpaceGUIDs         Filter `qs:"space_guids"`
	OrganizationGUIDs  Filter `qs:"organization_guids"`
	Available          *bool  `qs:"available"`
}

// NewServiceOfferingListOptions creates new options to pass to list
func NewServiceOfferingListOptions() *ServiceOfferingListOptions {
	return &ServiceOfferingListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceOfferingListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Delete the specified service offering
func (c *ServiceOfferingClient) Delete(guid string) error {
	_, err := c.client.delete(path.Format("/v3/service_offerings/%s", guid))
	return err
}

// Get the specified service offering
func (c *ServiceOfferingClient) Get(guid string) (*resource.ServiceOffering, error) {
	var ServiceOffering resource.ServiceOffering
	err := c.client.get(path.Format("/v3/service_offerings/%s", guid), &ServiceOffering)
	if err != nil {
		return nil, err
	}
	return &ServiceOffering, nil
}

// List pages service offerings the user has access to
func (c *ServiceOfferingClient) List(opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, *Pager, error) {
	if opts == nil {
		opts = NewServiceOfferingListOptions()
	}

	var res resource.ServiceOfferingList
	err := c.client.get(path.Format("/v3/service_offerings?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service offerings the user has access to
func (c *ServiceOfferingClient) ListAll(opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, error) {
	if opts == nil {
		opts = NewServiceOfferingListOptions()
	}
	return AutoPage[*ServiceOfferingListOptions, *resource.ServiceOffering](opts, func(opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, *Pager, error) {
		return c.List(opts)
	})
}

// Update the specified attributes of the service offering
func (c *ServiceOfferingClient) Update(guid string, r *resource.ServiceOfferingUpdate) (*resource.ServiceOffering, error) {
	var res resource.ServiceOffering
	_, err := c.client.patch(path.Format("/v3/service_offerings/%s", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
