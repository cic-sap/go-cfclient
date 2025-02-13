package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestUsers(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(16451)
	user := g.User().JSON
	user2 := g.User().JSON

	tests := []RouteTest{
		{
			Description: "Create user",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/users",
				Output:   g.Single(user),
				Status:   http.StatusCreated,
				PostForm: `{ "guid": "3ebeaa8b-fd55-4724-a764-9f2231d8f7db" }`,
			},
			Expected: user,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.UserCreate{
					GUID: "3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				}
				return c.Users.Create(r)
			},
		},
		{
			Description: "Delete user",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/users/3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Users.Delete("3ebeaa8b-fd55-4724-a764-9f2231d8f7db")
			},
		},
		{
			Description: "Get user",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/users/3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				Output:   g.Single(user),
				Status:   http.StatusOK},
			Expected: user,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Users.Get("3ebeaa8b-fd55-4724-a764-9f2231d8f7db")
			},
		},
		{
			Description: "List all users",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/users",
				Output:   g.Paged([]string{user}, []string{user2}),
				Status:   http.StatusOK},
			Expected: g.Array(user, user2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Users.ListAll(nil)
			},
		},
		{
			Description: "Update user",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/users/3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				Output:   g.Single(user),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: user,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.UserUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Users.Update("3ebeaa8b-fd55-4724-a764-9f2231d8f7db", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
