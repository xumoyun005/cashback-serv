package models

type Pagination struct {
	Limit     int64 `json:"-" default:"10"`
	Offset    int64 `json:"-" default:"0"`
	Page      int64 `json:"page" default:"1"`
	PageSize  int64 `json:"pageSize" default:"10"`
	PageTotal int64 `json:"pageTotal"`
	ItemTotal int64 `json:"itemTotal"`
}

func (p *Pagination) Calculate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	p.Limit = p.PageSize
	p.Offset = (p.Page - 1) * p.PageSize
}
