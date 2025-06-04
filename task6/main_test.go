package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_calcFibo(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{
			name: "zero",
			n:    0,
			want: 0,
		},
		{
			name: "one",
			n:    1,
			want: 1,
		},
		{
			name: "two",
			n:    2,
			want: 1,
		},
		{
			name: "negative",
			n:    -1,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcFibo(tt.n); got != tt.want {
				t.Errorf("calcFibo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fibo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fibo))
	defer server.Close()
	tests := []struct {
		name       string
		params     string
		wantStatus int
	}{
		{name: "good param",
			params:     "?N=3",
			wantStatus: 200,
		},
		{name: "bad param",
			params:     "?N=1.1",
			wantStatus: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := http.Get(server.URL + tt.params)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
