package pages

type PagePropsProvider struct {
	AppName string
}

func CreateProps[T any](
	pp *PagePropsProvider,
	pageName string,
	user *UserProps,
	data T,
) PageProps[T] {
	return PageProps[T]{
		AppName:  pp.AppName,
		PageName: pageName,
		User:     user,
		Data:     data,
	}
}
