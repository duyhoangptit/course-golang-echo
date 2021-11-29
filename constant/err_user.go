package constant

import "errors"

var (
	UserConflict     = errors.New("Người dùng đã tồn tại")
	RepoConflict     = errors.New("Repository đã tồn tại")
	SignUpFail       = errors.New("Đăng ký thất bại")
	UserNotFound     = errors.New("Người dùng không tồn tại")
	PasswordInvalid  = errors.New("Mật khẩu không đúng")
	RepoNotFound     = errors.New("Repo không tồn tại")
	BookmarkNotFound = errors.New("Bookmark không tồn tại")
	RepoNotUpdated   = errors.New("Repo update thất bại")
	BookmarkConflict = errors.New("Bookmark đã tồn tại")
	BookmarkFail     = errors.New("Bookmark thất bại")
	DelBookmarkFail  = errors.New("Xóa Bookmark thất bại")
)
