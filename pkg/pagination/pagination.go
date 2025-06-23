package pagination

const globalDefaultPerPage = 20

// Pagination 用來表示分頁
type Pagination struct {
	Page       uint32 `query:"page" form:"page" json:"page" description:"目前頁面"`
	PerPage    uint32 `query:"perPage" form:"perPage" json:"perPage" description:"每頁顯示多少筆"`
	TotalCount uint32 `query:"totalCount" form:"totalCount" json:"totalCount" description:"總筆數"`
	TotalPage  uint32 `query:"totalPage" form:"totalPage" json:"totalPage" description:"總頁數"`
}

// Count 用來儲存sql count 結果
type Count struct {
	Count uint32 `json:"count" description:"sql count result"`
}

// SetTotalCountAndPage 用來計算總數和分頁
func (p *Pagination) SetTotalCountAndPage(total uint32) {
	p.CheckOrSetDefault()
	p.TotalCount = total

	quotient := p.TotalCount / p.PerPage
	remainder := p.TotalCount % p.PerPage
	if remainder > 0 {
		quotient++
	}
	p.TotalPage = quotient
}

// CheckOrSetDefault 檢查Page值若未設置則設置預設值
func (p *Pagination) CheckOrSetDefault(params ...uint32) *Pagination {
	var defaultPerPage uint32
	if len(params) >= 1 {
		defaultPerPage = params[0]
	}

	if defaultPerPage <= 0 {
		defaultPerPage = globalDefaultPerPage
	}

	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = defaultPerPage
	}
	return p
}

// LimitAndOffset return limit and offset
func (p *Pagination) LimitAndOffset() (int, int) {
	if p.PerPage >= 5000 {
		p.PerPage = 5000
	}
	return int(p.PerPage), int(p.Offset())
}

// Offset 計算 offset 的值
func (p *Pagination) Offset() uint32 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}
