package domain

type PageInfo struct {
	Limit   int
	Offset  int
	Total   int
	HasMore bool
}

func NewPageInfo(limit, offset, total int) PageInfo {
	if limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return PageInfo{Limit: limit, Offset: offset, Total: total, HasMore: offset+limit < total}
}

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

func NormalizeSort(v string) SortDirection {
	if v == string(SortDesc) {
		return SortDesc
	}
	return SortAsc
}
