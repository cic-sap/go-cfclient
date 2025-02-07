package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type AppFeatureClient commonClient

// Get retrieves the named app feature
func (c *AppFeatureClient) Get(appGUID, featureName string) (*resource.AppFeature, error) {
	var a resource.AppFeature
	err := c.client.get(path.Format("/v3/apps/%s/features/%s", appGUID, featureName), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetSSH retrieves the SSH app feature
func (c *AppFeatureClient) GetSSH(appGUID string) (*resource.AppFeature, error) {
	return c.Get(appGUID, "ssh")
}

// GetRevisions retrieves the revisions app feature
func (c *AppFeatureClient) GetRevisions(appGUID string) (*resource.AppFeature, error) {
	return c.Get(appGUID, "revisions")
}

// List pages all app features
func (c *AppFeatureClient) List(appGUID string) ([]*resource.AppFeature, *Pager, error) {
	var res resource.AppFeatureList
	err := c.client.get(path.Format("/v3/apps/%s/features", appGUID), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// Update the enabled attribute of the named app feature
func (c *AppFeatureClient) Update(appGUID, featureName string, enabled bool) (*resource.AppFeature, error) {
	r := &resource.AppFeatureUpdate{
		Enabled: enabled,
	}
	var a resource.AppFeature
	_, err := c.client.patch(path.Format("/v3/apps/%s/features/%s", appGUID, featureName), r, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateSSH updated the enabled attribute of the SSH app feature
func (c *AppFeatureClient) UpdateSSH(appGUID string, enabled bool) (*resource.AppFeature, error) {
	return c.Update(appGUID, "ssh", enabled)
}

// UpdateRevisions updated the enabled attribute of the revisions app feature
func (c *AppFeatureClient) UpdateRevisions(appGUID string, enabled bool) (*resource.AppFeature, error) {
	return c.Update(appGUID, "revisions", enabled)
}
