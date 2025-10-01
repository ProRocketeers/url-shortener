package domain

type ErrorCode string

const (
	ErrorCodeSlugConflict ErrorCode = "slug_conflict"
	ErrorCodeLinkCreate   ErrorCode = "link_create_error"
	ErrorCodeLinkNotFound ErrorCode = "link_not_found"
	ErrorCodeLinkOther    ErrorCode = "link_other"
	ErrorCodeLinkExpired  ErrorCode = "link_expired"
)

type ShortLinkError struct {
	Code ErrorCode
}

func (e *ShortLinkError) Error() string {
	return string(e.Code)
}

func (e *ShortLinkError) Is(target error) bool {
	t, ok := target.(*ShortLinkError)
	return ok && e.Code == t.Code
}
