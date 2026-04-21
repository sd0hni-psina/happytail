package models

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func NewPaginatedResponse[T any](data []T, totalItems int, params PaginationParams) PaginatedResponse[T] {
	meta := NewPaginationMeta(totalItems, params)
	return PaginatedResponse[T]{
		Data: data,
		Meta: meta,
	}
}

func NewPaginationMeta(totalItems int, params PaginationParams) PaginationMeta {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	totalPages := (totalItems + params.Limit - 1) / params.Limit
	if params.Page > totalPages && totalPages > 0 {
		params.Page = totalPages
	}

	return PaginationMeta{
		CurrentPage: params.Page,
		Limit:       params.Limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

}

func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

func ParsePagination(r *http.Request) PaginationParams {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = DefaultPage
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > MaxLimit {
		limit = DefaultLimit
	}
	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}
