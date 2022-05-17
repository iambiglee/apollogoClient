package storage

const (
	ADDED ConfigChangeType = iota
	MODIFIED
	DELETED
)

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

func createModifyConfigChange(oldValue interface{}, newValue interface{}) *ConfigChange {
	return &ConfigChange{
		OldValue:   oldValue,
		NewValue:   newValue,
		ChangeType: MODIFIED,
	}
}

//create add config change
func createAddConfigChange(newValue interface{}) *ConfigChange {
	return &ConfigChange{
		NewValue:   newValue,
		ChangeType: ADDED,
	}
}

//create delete config change
func createDeletedConfigChange(oldValue interface{}) *ConfigChange {
	return &ConfigChange{
		OldValue:   oldValue,
		ChangeType: DELETED,
	}
}

//base on changeList create Change event
func createConfigChangeEvent(changes map[string]*ConfigChange, nameSpace string, notificationID int64) *ChangeEvent {
	c := &ChangeEvent{
		Changes: changes,
	}
	c.Namespace = nameSpace
	c.NotificationID = notificationID
	return c
}
