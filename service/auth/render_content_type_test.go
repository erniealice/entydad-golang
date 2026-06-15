package auth

// render_content_type_test.go — relocated from app
// service-admin/internal/composition/auth_content_type_test.go (Wave B D2a-3).
// The app's auth_bridge.go renderAuthView shim was deleted when the auth chain
// moved into the entydad block; the canonical implementation is
// (m *AuthModule).renderAuthView (helpers.go). This internal test pins the same
// W1 (Phase 2a, 20260531-csp-and-auth-cycle-remediation) invariant directly on
// the canonical method: the auth-shell bypass path (login /
// select-workspace-role) must commit Content-Type: text/html BEFORE WriteHeader
// so nosniff renders it as HTML, not plain text.

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// authSnapshotWriter captures the header map at the moment WriteHeader (or the
// first Write) is called — mirroring how a real net/http connection commits
// headers to the wire. httptest.ResponseRecorder keeps a single mutable header
// map, so a handler that sets Content-Type AFTER WriteHeader would look correct
// under a recorder but ship a header-less response in production. The pyeza
// renderer DOES set text/html internally, but only AFTER renderAuthView's
// WriteHeader — so only a snapshotting writer can prove the W1 (Phase 2a)
// ordering is intact.
type authSnapshotWriter struct {
	hdr       http.Header
	committed http.Header
	status    int
	body      strings.Builder
	wroteOnce bool
}

func newAuthSnapshotWriter() *authSnapshotWriter {
	return &authSnapshotWriter{hdr: make(http.Header)}
}

func (s *authSnapshotWriter) Header() http.Header { return s.hdr }

func (s *authSnapshotWriter) WriteHeader(status int) {
	if s.wroteOnce {
		return
	}
	s.wroteOnce = true
	s.status = status
	s.committed = s.hdr.Clone()
}

func (s *authSnapshotWriter) Write(b []byte) (int, error) {
	if !s.wroteOnce {
		s.WriteHeader(http.StatusOK)
	}
	return s.body.Write(b)
}

func (s *authSnapshotWriter) committedContentType() string {
	if s.committed == nil {
		return s.hdr.Get("Content-Type")
	}
	return s.committed.Get("Content-Type")
}

// TestRenderAuthView_SetsHTMLContentType pins the W1 (Phase 2a) renderAuthView
// fix: the auth-shell bypass path (login / select-workspace-role) must commit
// Content-Type: text/html before WriteHeader so nosniff renders it as HTML, not
// plain text. Uses a real pyeza renderer built from an in-memory FS so the
// renderer's internal (post-WriteHeader) CT-set cannot mask a regression in the
// caller's ordering.
func TestRenderAuthView_SetsHTMLContentType(t *testing.T) {
	const tmplName = "auth-shell-test"
	mapFS := fstest.MapFS{
		"auth.html": &fstest.MapFile{
			Data: []byte(`{{define "` + tmplName + `"}}<html><body>auth shell</body></html>{{end}}`),
		},
	}
	renderer := pyeza.NewHTMLRendererFromFS(mapFS)
	if err := renderer.Init(); err != nil {
		t.Fatalf("renderer Init failed: %v", err)
	}

	m := &AuthModule{deps: &Deps{Renderer: renderer}}

	w := newAuthSnapshotWriter()
	r := httptest.NewRequest(http.MethodGet, "/auth/select-workspace-role", nil)

	m.renderAuthView(w, r, view.ViewResult{
		Template:   tmplName,
		StatusCode: http.StatusOK,
	})

	if got := w.committedContentType(); got != "text/html; charset=utf-8" {
		t.Errorf("committed Content-Type = %q, want %q (renderAuthView must set it BEFORE WriteHeader; under nosniff a missing CT renders the auth shell as plain text)", got, "text/html; charset=utf-8")
	}
	if body := w.body.String(); !strings.Contains(body, "auth shell") {
		t.Errorf("rendered body = %q, want it to contain the auth-shell template output", body)
	}
}
