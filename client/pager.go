package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type Pager struct {
	NextPageURL     string
	PreviousPageURL string

	nextPageQSReader     *path.QuerystringReader
	previousPageQSReader *path.QuerystringReader
}

func NewPager(pagination resource.Pagination) *Pager {
	return &Pager{
		NextPageURL:     pagination.Next.Href,
		PreviousPageURL: pagination.Previous.Href,
	}
}

func (p *Pager) HasNextPage() bool {
	q, err := path.NewQuerystringReader(p.NextPageURL)
	if err != nil {
		return false
	}
	p.nextPageQSReader = q
	return true
}

func (p *Pager) NextPage(opts ListOptioner) {
	if !p.HasNextPage() {
		return
	}
	page := p.nextPageQSReader.Int(PageField)
	perPage := p.nextPageQSReader.Int(PerPageField)
	opts.CurrentPage(page, perPage)
}

func (p *Pager) HasPreviousPage() bool {
	q, err := path.NewQuerystringReader(p.PreviousPageURL)
	if err != nil {
		return false
	}
	p.previousPageQSReader = q
	return true
}

func (p *Pager) PreviousPage(opts ListOptioner) {
	if !p.HasPreviousPage() {
		return
	}
	page := p.previousPageQSReader.Int(PageField)
	perPage := p.previousPageQSReader.Int(PerPageField)
	opts.CurrentPage(page, perPage)
}

type ListFunc[T ListOptioner, R any] func(opts T) ([]R, *Pager, error)

func AutoPage[T ListOptioner, R any](opts T, list ListFunc[T, R]) ([]R, error) {
	var all []R
	for {
		page, pager, err := list(opts)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, nil
}
