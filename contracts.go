package entydad

import "github.com/erniealice/pyeza-golang/view"

// Re-export view contracts for backwards compatibility.
// The canonical definitions now live in github.com/erniealice/pyeza-golang/view.
type View = view.View
type ViewFunc = view.ViewFunc
type ViewContext = view.ViewContext
type ViewResult = view.ViewResult
type FlashMessage = view.FlashMessage
type RouteRegistrar = view.RouteRegistrar

var (
	OK       = view.OK
	Redirect = view.Redirect
	Error    = view.Error
)
