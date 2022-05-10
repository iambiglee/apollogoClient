package storage

type ChangeEvent struct {
	baseChangeEvent
	Changes map[string]*ConfigChange
}

//监听器
type ChangeListener interface {
	OnChange(event *ChangeEvent)
	//监控最新变更
	OnNewestChange(event *FullChangeEvent)
}

type ConfigChangeType int

type ConfigChange struct {
	OldValue   interface{}
	NewValue   interface{}
	ChangeType ConfigChangeType
}

type baseChangeEvent struct {
	Namespace      string
	NotificationID int64
}

type FullChangeEvent struct {
	baseChangeEvent
	Changes map[string]interface{}
}
