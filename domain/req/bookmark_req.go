package req

type BookmarkReq struct {
	RepoName string `json:"repo,omitempty" validate:"required"`
}
