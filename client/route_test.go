package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(123)
	route := g.Route().JSON
	route2 := g.Route().JSON
	domain := g.Domain().JSON
	space := g.Space().JSON
	space2 := g.Space().JSON
	org := g.Organization().JSON

	tests := []RouteTest{
		{
			Description: "Create route",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/routes",
				Output:   g.Single(route),
				Status:   http.StatusCreated,
				PostForm: `{
					"host": "a-hostname",
					"path": "/some_path",
					"port": 6666,
					"relationships": {
					  "domain": {
						"data": { "guid": "a99f869d-151a-4a80-95b7-653ada640824" }
					  },
					  "space": {
						"data": { "guid": "33d27af8-788d-4de5-8f37-fb80d517f2ed" }
					  }
					}
				  }`,
			},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewRouteCreateWithHost("a99f869d-151a-4a80-95b7-653ada640824",
					"33d27af8-788d-4de5-8f37-fb80d517f2ed",
					"a-hostname",
					"/some_path",
					6666)
				return c.Routes.Create(r)
			},
		},
		{
			Description: "Delete route",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.Delete("5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output:   g.Single(route),
				Status:   http.StatusOK,
			},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.Get("5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route include domain",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: route,
					Domains:  []string{domain},
				}),
				Status: http.StatusOK,
			},
			Expected:  route,
			Expected2: domain,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.GetIncludeDomain("5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route include space",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: route,
					Spaces:   []string{space},
				}),
				Status: http.StatusOK,
			},
			Expected:  route,
			Expected2: space,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.GetIncludeSpace("5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route include space and org",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource:      route,
					Spaces:        []string{space},
					Organizations: []string{org},
				}),
				Status: http.StatusOK,
			},
			Expected:  route,
			Expected2: space,
			Expected3: org,
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.Routes.GetIncludeSpaceAndOrg("5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "List all routes",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output:   g.Paged([]string{route}, []string{route2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(route, route2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.ListAll(nil)
			},
		},
		{
			Description: "List all routes for an app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/758c78dc-60bc-4f84-999b-247bdc2c37fe/routes",
				Output:   g.Paged([]string{route}, []string{route2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(route, route2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.ListForAppAll("758c78dc-60bc-4f84-999b-247bdc2c37fe", nil)
			},
		},
		{
			Description: "List all routes and include domains",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{route},
						Domains:   []string{domain},
					},
					testutil.PagedResult{
						Resources: []string{route2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(route, route2),
			Expected2: g.Array(domain),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.ListIncludeDomainsAll(nil)
			},
		},
		{
			Description: "List all routes and include spaces",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{route},
						Spaces:    []string{space},
					},
					testutil.PagedResult{
						Resources: []string{route2},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(route, route2),
			Expected2: g.Array(space, space2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.ListIncludeSpacesAll(nil)
			},
		},
		{
			Description: "List all routes and include spaces and orgs",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:     []string{route},
						Spaces:        []string{space},
						Organizations: []string{org},
					},
					testutil.PagedResult{
						Resources: []string{route2},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(route, route2),
			Expected2: g.Array(space, space2),
			Expected3: g.Array(org),
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.Routes.ListIncludeSpacesAndOrgsAll(nil)
			},
		},
		{
			Description: "Update route",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output:   g.Single(route),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"key": "value"}, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.RouteUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Routes.Update("5a85c020-3e3d-42a5-a475-5084c5357e82", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
