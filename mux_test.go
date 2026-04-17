package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestBasicMatch verifies that a registered route returns 200.
func TestBasicMatch(t *testing.T) {
	m := NewMux()
	m.Get("/path", func(w http.ResponseWriter, r *http.Request) {})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// TestNotFound verifies that an unregistered path returns 404.
func TestNotFound(t *testing.T) {
	m := NewMux()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

// TestMethodNotAllowed verifies that a path match with wrong method returns 405.
func TestMethodNotAllowed(t *testing.T) {
	m := NewMux()
	rec := httptest.NewRecorder()
	m.Get("/path", func(w http.ResponseWriter, r *http.Request) {})
	req, err := http.NewRequest("POST", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

// TestURLParams verifies that URL parameters are captured and accessible via URLParam.
func TestURLParams(t *testing.T) {
	m := NewMux()
	var gotID string
	m.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		gotID = URLParam(r, "id")
	})
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/users/42", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if gotID != "42" {
		t.Errorf("expected id=42, got %s", gotID)
	}
}

// TestMiddleware verifies that router-level middleware runs and can modify the response.
func TestMiddleware(t *testing.T) {
	m := NewMux()
	m.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "true")
			next.ServeHTTP(w, r)
		})
	})
	m.Get("/path", func(w http.ResponseWriter, r *http.Request) {})

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)

	if rec.Header().Get("X-Test") != "true" {
		t.Errorf("expected X-Test header to be set")
	}
}

// TestGroupPrefix verifies that routes registered inside a group are reachable at the prefixed path.
func TestGroupPrefix(t *testing.T) {
	m := NewMux()
	m.Group("/testGroup", func(r Router) {
		r.Get("/path", func(w http.ResponseWriter, req *http.Request) {})
	})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/testGroup/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// TestGroupMiddleware verifies that group middleware only fires for routes inside the group.
func TestGroupMiddleware(t *testing.T) {
	m := NewMux()
	m.Group("/testGroup", func(r Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Group", "true")
				next.ServeHTTP(w, r)
			})
		})
		r.Get("/path", func(w http.ResponseWriter, r *http.Request) {})
	})
	m.Get("/path", func(w http.ResponseWriter, r *http.Request) {})

	rec1 := httptest.NewRecorder()
	req1, err := http.NewRequest("GET", "/testGroup/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec1, req1)

	rec2 := httptest.NewRecorder()
	req2, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec2, req2)

	if rec1.Header().Get("X-Group") != "true" {
		t.Errorf("expected X-Group header on group route")
	}
	if rec2.Header().Get("X-Group") != "" {
		t.Errorf("expected X-Group header to be absent on non-group route")
	}
}

// TestWildcard verifies that a trailing * segment matches the rest of the path
// and the captured tail is accessible via URLParam(r, "*").
func TestWildcard(t *testing.T) {
	m := NewMux()
	var gotFile string
	m.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		gotFile = URLParam(r, "*")
	})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/files/docs/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if gotFile != "docs/test.txt" {
		t.Errorf("expected *=docs/test.txt, got %s", gotFile)
	}
}

// TestTrailingSlash verifies that a request with a mismatched trailing slash
// is redirected (301) to the alternate path when RedirectTrailingSlash is enabled.
func TestTrailingSlash(t *testing.T) {
	m := NewMux()
	m.Get("/users", func(w http.ResponseWriter, r *http.Request) {})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if rec.Code != http.StatusMovedPermanently || rec.Header().Get("Location") != "/users" {
		t.Errorf("expected 301, location /users, got %d location %s", rec.Code, rec.Header().Get("Location"))
	}
}

// TestMethodChaining verifies that Route returns a RouteBuilder that can register
// multiple HTTP methods on the same pattern via chaining.
func TestMethodChaining(t *testing.T) {
	m := NewMux()
	m.Route("/chain").Get(func(w http.ResponseWriter, r *http.Request) {}).Post(func(w http.ResponseWriter, r *http.Request) {})
	rec1 := httptest.NewRecorder()
	req1, err := http.NewRequest("GET", "/chain", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec1, req1)
	rec2 := httptest.NewRecorder()
	req2, err := http.NewRequest("POST", "/chain", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec2, req2)
	if rec1.Code != http.StatusOK || rec2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d and %d", rec1.Code, rec2.Code)
	}
}

// TestNamedWildcard verifies that a {name:*} wildcard segment matches the rest of the path
// and the captured tail is accessible via URLParam using the given name.
func TestNamedWildcard(t *testing.T) {
	m := NewMux()
	var gotFile string
	m.Get("/files/{path:*}", func(w http.ResponseWriter, r *http.Request) {
		gotFile = URLParam(r, "path")
	})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/files/docs/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec, req)
	if gotFile != "docs/test.txt" {
		t.Errorf("expected path=docs/test.txt, got %s", gotFile)
	}
}

// TestRouteBuilderMiddleware verifies that middleware registered via RouteBuilder.Use
// only applies to routes registered through that builder, not to other routes.
func TestRouteBuilderMiddleware(t *testing.T) {
	m := NewMux()
	m.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Outside", "true")
			next.ServeHTTP(w, r)
		})
	})
	m.Get("/test", func(w http.ResponseWriter, r *http.Request) {})
	m.Route("/chain-test").Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Chain", "true")
			next.ServeHTTP(w, r)
		})
	}).Get(func(w http.ResponseWriter, r *http.Request) {})

	rec1 := httptest.NewRecorder()
	req1, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec1, req1)

	rec2 := httptest.NewRecorder()
	req2, err := http.NewRequest("GET", "/chain-test", nil)
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(rec2, req2)

	if rec1.Header().Get("X-Test-Outside") != "true" {
		t.Errorf("expected X-Test-Outside header on /test")
	}
	if rec1.Header().Get("X-Test-Chain") != "" {
		t.Errorf("expected X-Test-Chain header to be absent on /test")
	}
	if rec2.Header().Get("X-Test-Chain") != "true" {
		t.Errorf("expected X-Test-Chain header on /chain-test")
	}
}
