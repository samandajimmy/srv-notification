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
	NotificationStatusFailed  = NotificationStatus("Failed")
	NotificationStatusSuccess = NotificationStatus("Success")
)

type NotificationType = string

const (
	NotificationEmail = NotificationType("Email")
	NotificationFCM   = NotificationType("FCM")
)
