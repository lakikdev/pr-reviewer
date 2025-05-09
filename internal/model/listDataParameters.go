package model

// swagger:model ListUserResponse
type ListDataParameters struct {
	Page   PageData   `json:"page"`
	Sort   []SortData `json:"sort"`
	Filter FilterData `json:"filter"`
}

type PageData struct {
	Number int `json:"number"`
	Size   int `json:"size"`
}

type SortData struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}

type FilterData struct {
	Items             []FilterItemData `json:"items"`
	QuickFilterValues []string         `json:"quickFilterValues"`
	InDateRange       *DateRange       `json:"inDateRange"`
}

type DateRange struct {
	StartAt FilterItemData `json:"startAt"`
	EndAt   FilterItemData `json:"endAt"`
}

type FilterItemData struct {
	Field         string `json:"columnField"`
	OperatorValue string `json:"operatorValue"`
	Value         string `json:"value"`
}

func (l *ListDataParameters) Verify() error {
	if l.Page.Size <= 0 {
		l.Page.Size = 50
	}

	if l.Sort == nil {
		l.Sort = make([]SortData, 0)
	}

	if l.Filter.Items == nil {
		l.Filter.Items = make([]FilterItemData, 0)
	}

	if l.Filter.QuickFilterValues == nil {
		l.Filter.QuickFilterValues = make([]string, 0)
	}

	return nil
}
