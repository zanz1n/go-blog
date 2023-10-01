package pages

import "strings"

type PagePropsProvider struct {
	AppName string
	Routes  []CreateRouteInfo
}

func CreateProps[T any](
	pp *PagePropsProvider,
	path string,
	user *UserProps,
	data T,
) PageProps[T] {
	pageName := ""

	routes := make([]Route, len(pp.Routes))
	for i, v := range pp.Routes {
		isCurrent := strings.Contains(path, v.Href)
		if isCurrent {
			pageName = v.Name
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
