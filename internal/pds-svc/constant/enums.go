package constant

/// Control Status

type ControlStatus = int8

const (
	Enabled = ControlStatus(iota + 1)
	Disabled
)

/// Asset Type

type AssetType = int8

const (
	ThumbnailAsset = AssetType(iota + 1)
)

type NotificationStatus = string

const (
	NotificationStatusFailed  = "Failed"
	NotificationStatusSuccess = "Success"
)

type NotificationType = string

const (
	NotificationEmail = "Email"
	NotificationFCM   = "FCM"
)
