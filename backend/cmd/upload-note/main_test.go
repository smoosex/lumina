package main

import "testing"

func TestSignaturePath(t *testing.T) {
	tests := map[string]string{
		"":                                      "/",
		"/api":                                 "/api",
		"/api/admin/reading-notes/example":     "/api/admin/reading-notes/example",
		"/lumina/api/admin/reading-notes/post": "/api/admin/reading-notes/post",
		"/nested/lumina/api/reading-notes":     "/api/reading-notes",
		"/notes/example":                       "/notes/example",
	}

	for input, want := range tests {
		got := signaturePath(input)
		if got != want {
			t.Fatalf("signaturePath(%q) = %q, want %q", input, got, want)
		}
	}
}
