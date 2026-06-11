package identity

import (
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"
)

// identityRouteRegistrarFull extends view.RouteRegistrar with HandleFunc support
// for raw http.HandlerFunc routes (e.g., JSON search endpoints). Apps that
// implement HandleFunc (e.g., service-admin chi router wrapper) can register
// JSON endpoints; apps that don't will skip.
type identityRouteRegistrarFull interface {
	view.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// identityHandleFunc is a nil-safe helper that registers an http.HandlerFunc
// route if the RouteRegistrar supports it, otherwise logs a warning and skips.
func identityHandleFunc(r view.RouteRegistrar, method, path string, handler http.HandlerFunc) {
	if handler == nil {
		return
	}
	if full, ok := r.(identityRouteRegistrarFull); ok {
		full.HandleFunc(method, path, handler)
		return
	}
	log.Printf("identity: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
}
