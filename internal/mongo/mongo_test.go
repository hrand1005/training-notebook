package mongo

import (
	"testing"
)

const (
	TestURI = "mongodb://localhost:27017"
	TestDB  = "test"
)

// TestNew tests the main constructor for the mongo package.
func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		uri        string
		wantHandle bool
		wantErr    bool
	}{
		{
			name:       "Nominal case returns handle and nil error",
			uri:        TestURI,
			wantHandle: true, // non-nil
			wantErr:    false,
		},
		{
			name:       "Invalid URI returns error",
			uri:        "invalid-uri",
			wantHandle: false,
			wantErr:    true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handle, err := New(tc.uri, TestDB)
			gotHandle := handle != nil
			gotErr := err != nil

			if tc.wantHandle != gotHandle {
				t.Fatalf("want handle: %v\ngot handle: %v\nerr: %v", tc.wantHandle, gotHandle, err)
			}
			if tc.wantErr != gotErr {
				t.Fatalf("want err: %v\ngot err: %v\nerr: %v", tc.wantErr, gotErr, err)
			}

			if handle != nil {
				handle.Close()
			}
		})
	}
}
