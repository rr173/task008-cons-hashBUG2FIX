package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOwnersRejectsBlankKeys is a regression test for the batch endpoint.
// The single GET /owner handler rejects empty or whitespace-only keys with a
// 400, but POST /owners accepted them and returned a mapping. This asserts the
// batch endpoint mirrors the single-query validation: any blank key in the
// request must yield HTTP 400 instead of 200.
func TestOwnersRejectsBlankKeys(t *testing.T) {
	srv := NewService()
	if _, err := srv.AddNode("alpha", 0); err != nil {
		t.Fatalf("add node: %v", err)
	}
	ts := httptest.NewServer(buildMux(srv))
	defer ts.Close()

	cases := []struct {
		name string
		body string
	}{
		{"empty string among valid", `{"keys":["valid", "", "also-valid"]}`},
		{"pure whitespace", `{"keys":["valid", "   "]}`},
		{"only empty string", `{"keys":[""]}`},
		{"only whitespace", `{"keys":["  \t "]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Post(ts.URL+"/owners", "application/json", bytes.NewBufferString(tc.body))
			if err != nil {
				t.Fatalf("post owners: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400 (blank keys must be rejected like GET /owner)", resp.StatusCode)
			}
		})
	}
}
