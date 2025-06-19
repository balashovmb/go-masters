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
		{
			name: "forty two",
			n:    42,
			want: 165580141,
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
	router := http.NewServeMux()
	router.HandleFunc("/fibo", fibo)
	server := httptest.NewServer(router)
	defer server.Close()

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "good param",
			url:        "/fibo?N=3",
			wantStatus: http.StatusOK,
		},
		{
			name:       "bad param",
			url:        "/fibo?N=1.1",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.url)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
