// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package clevergo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouteGroup(t *testing.T) {
	tests := []struct {
		path         string
		expectedPath string
		shouldPanic  bool
	}{
		{"without-prefix-slash", "", true},
		{"/", "/", false},
		{"//", "/", false},
		{"/users", "/users", false},
		{"/users/", "/users", false},
	}

	router := NewRouter()
	for _, test := range tests {
		if test.shouldPanic {
			recv := catchPanic(func() {
				newRouteGroup(router, test.path)
			})
			if recv == nil {
				t.Error("expected a panic")
			}
			continue
		}

		route := newRouteGroup(router, test.path)
		if test.expectedPath != route.path {
			t.Errorf("expected path %q, got %q", test.expectedPath, route.path)
		}
		if test.expectedPath != route.name {
			t.Errorf("expected name %q, got %q", test.expectedPath, route.name)
		}
	}
}

func ExampleRouteGroup() {
	router := NewRouter()
	api := router.Group("/api", RouteGroupMiddleware(echoMiddleware("api")))

	v1 := api.Group("/v1", RouteGroupMiddleware(
		echoMiddleware("v1"),
		echoMiddleware("authenticate"),
	))
	v1.Get("/users/:name", func(ctx *Context) error {
		ctx.WriteString(fmt.Sprintf("user: %s", ctx.Params.String("name")))
		return nil
	}, RouteMiddleware(
		echoMiddleware("fizz1"),
		echoMiddleware("fizz2"),
	))

	v2 := api.Group("/v2", RouteGroupMiddleware(
		echoMiddleware("v2"),
		echoMiddleware("authenticate"),
	))
	v2.Get("/users/:name", func(ctx *Context) error {
		ctx.WriteString(fmt.Sprintf("user: %s", ctx.Params.String("name")))
		return nil
	}, RouteMiddleware(
		echoMiddleware("buzz1"),
		echoMiddleware("buzz2"),
	))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/foo", nil)
	router.ServeHTTP(w, req)
	fmt.Println(w.Body.String())

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v2/users/bar", nil)
	router.ServeHTTP(w, req)
	fmt.Println(w.Body.String())

	// Output:
	// api v1 authenticate fizz1 fizz2 user: foo
	// api v2 authenticate buzz1 buzz2 user: bar
}

func TestRouteGroupName(t *testing.T) {
	for _, name := range []string{"foo", "bar"} {
		g := &RouteGroup{}
		RouteGroupName(name)(g)
		if g.name != name {
			t.Errorf("expected name %s, got %s", name, g.name)
		}
	}
}
