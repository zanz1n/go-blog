package pages

type PagePropsProvider struct {
	AppName string
	Routes  []CreateRouteInfo
}

func CreateProps[T any](
	pp *PagePropsProvider,
	path string,
	pageName string,
	user *UserProps,
	data T,
) PageProps[T] {
	routes := make([]Route, len(pp.Routes))
	for i, v := range pp.Routes {
		isCurrent := false
		if len(path) > 0 {
			isCurrent = v.Href == path || v.Href == path[0:len(path)-1]
		}

		routes[i] = Route{
			IsCurrent: isCurrent,
			Name:      v.Name,
			Href:      v.Href,
		}
	}

	return PageProps[T]{
		AppName:  pp.AppName,
		PageName: pageName,
		Routes:   routes,
		User:     user,
		Data:     data,
	}
}
