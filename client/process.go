package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ProcessClient commonClient

// ProcessOptions list filters
type ProcessOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
}

// NewProcessOptions creates new options to pass to list
func NewProcessOptions() *ProcessOptions {
	return &ProcessOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ProcessOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Get the specified process
func (c *ProcessClient) Get(guid string) (*resource.Process, error) {
	var iso resource.Process
	err := c.client.get(path.Format("/v3/processes/%s", guid), &iso)
	if err != nil {
		return nil, err
	}
	return &iso, nil
}

// GetStats for the specified process
func (c *ProcessClient) GetStats(guid string) (*resource.ProcessStats, error) {
	var stats resource.ProcessStats
	err := c.client.get(path.Format("/v3/processes/%s/stats", guid), &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// List pages all processes
func (c *ProcessClient) List(opts *ProcessOptions) ([]*resource.Process, *Pager, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}

	var isos resource.ProcessList
	err := c.client.get(path.Format("/v3/processes?%s", opts.ToQueryString()), &isos)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(isos.Pagination)
	return isos.Resources, pager, nil
}

// ListAll retrieves all processes
func (c *ProcessClient) ListAll(opts *ProcessOptions) ([]*resource.Process, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}
	return AutoPage[*ProcessOptions, *resource.Process](opts, func(opts *ProcessOptions) ([]*resource.Process, *Pager, error) {
		return c.List(opts)
	})
}

// ListForApp pages all processes for the specified app
func (c *ProcessClient) ListForApp(appGUID string, opts *ProcessOptions) ([]*resource.Process, *Pager, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}

	var processes resource.ProcessList
	err := c.client.get(path.Format("/v3/apps/%s/processes?%s", appGUID, opts.ToQueryString()), &processes)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(processes.Pagination)
	return processes.Resources, pager, nil
}

// ListForAppAll retrieves all processes for the specified app
func (c *ProcessClient) ListForAppAll(appGUID string, opts *ProcessOptions) ([]*resource.Process, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}
	return AutoPage[*ProcessOptions, *resource.Process](opts, func(opts *ProcessOptions) ([]*resource.Process, *Pager, error) {
		return c.ListForApp(appGUID, opts)
	})
}

// Scale the process using the specified scaling requirements
func (c *ProcessClient) Scale(guid string, scale *resource.ProcessScale) (*resource.Process, error) {
	var process resource.Process
	_, err := c.client.post(path.Format("/v3/processes/%s/actions/scale", guid), scale, &process)
	if err != nil {
		return nil, err
	}
	return &process, nil
}

// Update the specified attributes of the process
func (c *ProcessClient) Update(guid string, r *resource.ProcessUpdate) (*resource.Process, error) {
	var process resource.Process
	_, err := c.client.patch(path.Format("/v3/processes/%s", guid), r, &process)
	if err != nil {
		return nil, err
	}
	return &process, nil
}

// Terminate an instance of a specific process. Health management will eventually restart the instance.
func (c *ProcessClient) Terminate(guid string, index int) error {
	_, err := c.client.delete(path.Format("/v3/processes/%s/instances/%d", guid, index))
	return err
}
