package client

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type DeploymentClient commonClient

// DeploymentListOptions list filters
type DeploymentListOptions struct {
	*ListOptions

	AppGUIDs      Filter `qs:"app_guids"`
	States        Filter `qs:"states"`
	StatusReasons Filter `qs:"status_reasons"`
	StatusValues  Filter `qs:"status_values"`
}

// NewDeploymentListOptions creates new options to pass to list
func NewDeploymentListOptions() *DeploymentListOptions {
	return &DeploymentListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o DeploymentListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Cancel the ongoing deployment
func (c *DeploymentClient) Cancel(guid string) error {
	_, err := c.client.post(path.Format("/v3/deployments/%s/actions/cancel", guid), nil, nil)
	return err
}

// Create a new deployment
func (c *DeploymentClient) Create(r *resource.DeploymentCreate) (*resource.Deployment, error) {
	// validate the params
	if r.Droplet != nil && r.Revision != nil {
		return nil, errors.New("droplet and revision cannot both be set")
	}

	var d resource.Deployment
	_, err := c.client.post("/v3/deployments", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Get the specified deployment
func (c *DeploymentClient) Get(guid string) (*resource.Deployment, error) {
	var d resource.Deployment
	err := c.client.get(path.Format("/v3/deployments/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages deployments the user has access to
func (c *DeploymentClient) List(opts *DeploymentListOptions) ([]*resource.Deployment, *Pager, error) {
	if opts == nil {
		opts = NewDeploymentListOptions()
	}
	var res resource.DeploymentList
	err := c.client.get(path.Format("/v3/deployments?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all deployments the user has access to
func (c *DeploymentClient) ListAll(opts *DeploymentListOptions) ([]*resource.Deployment, error) {
	if opts == nil {
		opts = NewDeploymentListOptions()
	}
	return AutoPage[*DeploymentListOptions, *resource.Deployment](opts, func(opts *DeploymentListOptions) ([]*resource.Deployment, *Pager, error) {
		return c.List(opts)
	})
}

// Update the specified attributes of the deployment
func (c *DeploymentClient) Update(guid string, r *resource.DeploymentUpdate) (*resource.Deployment, error) {
	var d resource.Deployment
	_, err := c.client.patch(path.Format("/v3/deployments/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
