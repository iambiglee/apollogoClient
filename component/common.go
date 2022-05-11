package component

type AbsComponent interface {
	Start()
}

// StartRefreshConfig 开始定时服务，但是好像只是一个接口，没搞懂
func StartRefreshConfig(component AbsComponent) {
	component.Start()
}
