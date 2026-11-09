package request

type RetrieveCodeRequest struct {
	FileName string `json:"file_name" binding:"required"`
}
