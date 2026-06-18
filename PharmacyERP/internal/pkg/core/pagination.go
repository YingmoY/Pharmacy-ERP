package core

// PageQuery is the standard list query struct.
type PageQuery struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

func (q *PageQuery) GetOffset() int {
	return (q.Page - 1) * q.PageSize
}

// PageResult is the standard paginated response struct.
type PageResult struct {
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	List       interface{} `json:"list"`
}

func NewPageResult(total int64, page, pageSize int, list interface{}) PageResult {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	return PageResult{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		List:       list,
	}
}
