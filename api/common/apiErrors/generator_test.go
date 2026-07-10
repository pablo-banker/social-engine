package apiErrors

import (
	"net/http"
	"strings"
	"testing"
)

func TestRenderJSON(t *testing.T) {
	cases := []struct {
		name    string
		in      []*APIError
		want    string
		wantErr string
	}{
		{
			name: "single entry",
			in: []*APIError{
				New(101, "boom", http.StatusInternalServerError),
			},
			want: "{\n  \"DOBRO-000101\": \"boom\"\n}\n",
		},
		{
			name: "century change inserts blank line",
			in: []*APIError{
				New(101, "a", http.StatusOK),
				New(102, "b", http.StatusOK),
				New(201, "c", http.StatusOK),
			},
			want: "{\n  \"DOBRO-000101\": \"a\",\n  \"DOBRO-000102\": \"b\",\n\n  \"DOBRO-000201\": \"c\"\n}\n",
		},
		{
			name: "input order is preserved across centuries",
			in: []*APIError{
				New(201, "a", http.StatusOK),
				New(101, "b", http.StatusOK),
			},
			want: "{\n  \"DOBRO-000201\": \"a\",\n\n  \"DOBRO-000101\": \"b\"\n}\n",
		},
		{
			name: "escapes special chars in message",
			in: []*APIError{
				New(101, `msg with "quotes" and \ slashes`, http.StatusOK),
			},
			want: "{\n  \"DOBRO-000101\": \"msg with \\\"quotes\\\" and \\\\ slashes\"\n}\n",
		},
		{
			name: "empty input",
			in:   []*APIError{},
			want: "{\n}\n",
		},
		{
			name:    "nil entry returns error",
			in:      []*APIError{nil},
			wantErr: "nil entry",
		},
		{
			name: "duplicate code returns error",
			in: []*APIError{
				New(101, "a", http.StatusOK),
				New(101, "b", http.StatusOK),
			},
			wantErr: "duplicate code DOBRO-000101",
		},
		{
			name: "malformed code returns error",
			in: []*APIError{
				{Code: "DOBRO-abc", Message: "x"},
			},
			wantErr: "invalid code",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RenderJSON(tc.in)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tc.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("output mismatch\nwant:\n%s\ngot:\n%s", tc.want, string(got))
			}
		})
	}
}

func TestRenderJSONOnAllRegistry(t *testing.T) {
	got, err := RenderJSON(All)
	if err != nil {
		t.Fatalf("RenderJSON(All) returned error: %v", err)
	}
	if !strings.HasPrefix(string(got), "{\n") || !strings.HasSuffix(string(got), "}\n") {
		t.Errorf("malformed envelope")
	}
	if !strings.Contains(string(got), "\"DOBRO-000419\": \"error getting user by username\"") {
		t.Errorf("expected ErrGetUserByUsername entry to be present in rendered output")
	}
}
