package mdpdf

type ConvertRequest struct {
	Markdown string `json:"markdown"`
	Theme    string `json:"theme"`     // "github" | "minimal"
	PageSize string `json:"page_size"` // "A4" | "Letter"
}
