package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) CreateSpace(r v3.CreateSpaceRequest) (*v3.Space, error) {
	req := c.NewRequest("POST", "/v3/spaces")
	params := map[string]interface{}{
		"name": r.Name,
		"relationships": map[string]interface{}{
			"organization": v3.ToOneRelationship{
				Data: v3.Relationship{
					GUID: r.OrgGUID,
				},
			},
		},
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	req.obj = params
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating  space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating  space %s, response code: %d", r.Name, resp.StatusCode)
	}

	var space v3.Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, errors.Wrap(err, "Error reading  space JSON")
	}

	return &space, nil
}

func (c *Client) GetSpaceByGUID(spaceGUID string) (*v3.Space, error) {
	req := c.NewRequest("GET", "/v3/spaces/"+spaceGUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting  space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting  space with GUID [%s], response code: %d", spaceGUID, resp.StatusCode)
	}

	var space v3.Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, errors.Wrap(err, "Error reading  space JSON")
	}

	return &space, nil
}

func (c *Client) DeleteSpace(spaceGUID string) error {
	req := c.NewRequest("DELETE", "/v3/spaces/"+spaceGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting  space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting  space with GUID [%s], response code: %d", spaceGUID, resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateSpace(spaceGUID string, r v3.UpdateSpaceRequest) (*v3.Space, error) {
	req := c.NewRequest("PATCH", "/v3/spaces/"+spaceGUID)
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	if len(params) > 0 {
		req.obj = params
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating  space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating  space %s, response code: %d", spaceGUID, resp.StatusCode)
	}

	var space v3.Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, errors.Wrap(err, "Error reading  space JSON")
	}

	return &space, nil
}

func (c *Client) ListSpacesByQuery(query url.Values) ([]v3.Space, error) {
	var spaces []v3.Space
	requestURL := "/v3/spaces"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  spaces")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing  spaces, response code: %d", resp.StatusCode)
		}

		var data v3.ListSpacesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  spaces")
		}

		spaces = append(spaces, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for  spaces")
		}
	}

	return spaces, nil
}

// ListSpaceUsers lists users by space GUID
func (c *Client) ListSpaceUsers(spaceGUID string) ([]v3.User, error) {
	var users []v3.User
	requestURL := "/v3/spaces/" + spaceGUID + "/users"

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  space users")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing  space users, response code: %d", resp.StatusCode)
		}

		var data v3.ListSpaceUsersResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  space users")
		}
		users = append(users, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for  space users")
		}
	}

	return users, nil
}
