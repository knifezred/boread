package dto

type TagRequest struct {
	TagName string `json:"tagName" binding:"required,max=64"`
}

type TagSearch struct {
	PageRequest
	TagName string `json:"tagName"`
}