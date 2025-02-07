package resource

import "time"

type Route struct {
	Guid         string             `json:"guid"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Host         string             `json:"host"`
	Path         string             `json:"path"`
	Url          string             `json:"url"`
	Protocol     string             `json:"protocol"`
	Port         int                `json:"port"`
	Destinations []RouteDestination `json:"destinations"`

	Metadata      Metadata                   `json:"metadata"`
	Relationships ToSpaceDomainRelationships `json:"relationships"`
	Links         map[string]Link            `json:"links"`
}

type RouteCreate struct {
	Relationships RouteRelationships `json:"relationships"`
	Host          *string            `json:"host,omitempty"`
	Path          *string            `json:"path,omitempty"`
	Port          *int               `json:"port,omitempty"`
	Metadata      *Metadata          `json:"metadata,omitempty"`
}

type RouteUpdate struct {
	Metadata *Metadata `json:"metadata"`
}

type RouteList struct {
	Pagination Pagination     `json:"pagination"`
	Resources  []*Route       `json:"resources"`
	Included   *RouteIncluded `json:"included"`
}

type RouteDestination struct {
	GUID     string              `json:"guid"`
	Route    RouteDestinationApp `json:"app"`
	Weight   *int                `json:"weight"`
	Port     int                 `json:"port"`
	Protocol string              `json:"protocol"`
}

type RouteDestinationApp struct {
	GUID    string                     `json:"guid"`
	Process RouteDestinationAppProcess `json:"process"`
}

type RouteDestinationAppProcess struct {
	Type string `json:"type"`
}

type ToSpaceDomainRelationships struct {
	Space  ToOneRelationship `json:"space"`
	Domain ToOneRelationship `json:"domain"`
}

type RouteRelationships struct {
	Space  ToOneRelationship `json:"space"`
	Domain ToOneRelationship `json:"domain"`
}

type RouteWithIncluded struct {
	Route
	Included *RouteIncluded `json:"included"`
}

type RouteIncluded struct {
	Organizations []*Organization `json:"organizations"`
	Spaces        []*Space        `json:"spaces"`
	Domains       []*Domain       `json:"domains"`
}

// RouteIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type RouteIncludeType int

const (
	RouteIncludeNone RouteIncludeType = iota
	RouteIncludeSpace
	RouteIncludeSpaceOrganization
	RouteIncludeDomain
)

func (a RouteIncludeType) String() string {
	switch a {
	case RouteIncludeSpace:
		return "space"
	case RouteIncludeSpaceOrganization:
		return "space.organization"
	case RouteIncludeDomain:
		return "domain"
	}
	return ""
}

func NewRouteCreate(domainGUID, spaceGUID string) *RouteCreate {
	return &RouteCreate{
		Relationships: RouteRelationships{
			Space: ToOneRelationship{
				Data: &Relationship{
					GUID: spaceGUID,
				},
			},
			Domain: ToOneRelationship{
				Data: &Relationship{
					GUID: domainGUID,
				},
			},
		},
	}
}

func NewRouteCreateWithHost(domainGUID, spaceGUID, host, path string, port int) *RouteCreate {
	rc := NewRouteCreate(domainGUID, spaceGUID)
	rc.Host = &host
	rc.Path = &path
	rc.Port = &port
	return rc
}
