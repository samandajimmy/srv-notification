package constant

const (
	NumCharSet           = "0123456789"
	AlphaCharSet         = "abcdefghijklmnopqrstuvwxyz"
	AlphaUpperCharSet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaNumCharSet      = AlphaCharSet + NumCharSet
	AlphaNumUpperCharSet = AlphaUpperCharSet + NumCharSet
	AlphaNumRandomSet    = AlphaCharSet + AlphaUpperCharSet + NumCharSet
	SlugRandomSet        = AlphaNumRandomSet + "-"
	StatusInactive       = "Inactive"
)

const (
	SubjectKey     = "subject"
	BuildHashKey   = "build_hash"
	ApplicationKey = "app"
)

const (
	SubjectIDHeader   = "x-subject-id"
	SubjectNameHeader = "x-subject-name"
	SubjectRoleHeader = "x-subject-role"
)

const (
	DefaultConfig = "DEFAULT_CONFIG"
)
