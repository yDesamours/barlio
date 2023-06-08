package metadata

type Metadata struct {
	TotalResult int `json:"totalresult,omitempty"`
	TotalPage   int `json:"totalpage,omitempty"`
	PageSize    int `json:"pagesize,omitempty"`
	PageNumber  int `json:"pagenumber,omitempty"`
}
