package mongodb

import (
	"context"
	// "errors"
	// "reflect"
	"testing"

	"github.com/hrand1005/training-notebook/internal/app"
)

// TestInsertSet tests the mongo implementation of the SetStore's Insert() method.
func TestInsertSet(t *testing.T) {
	handle, err := New(TestURI, TestDB)
	if err != nil {
		t.Fatalf("failed to initialize test db handle: %v", err)
	}
	defer handle.Delete()
	defer handle.Close()

	validStore := newValidSetStore(handle)
	invalidStore := newInvalidSetStore(handle)

	tests := []struct {
		name    string
		store   app.SetStore
		set     *app.Set
		wantID  bool
		wantErr bool
	}{
		{
			name:  "Nominal case returns ID and nil error",
			store: validStore,
			set: &app.Set{
				Movement:  "Squat",
				Intensity: 80,
				Volume:    5,
			},
			wantID:  true,
			wantErr: false,
		},
		{
			name:  "Invalid set store case returns empty ID and error",
			store: invalidStore,
			set: &app.Set{
				Movement:  "Squat",
				Intensity: 80,
				Volume:    5,
			},
			wantID:  false,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, err := tc.store.Insert(tc.set)
			gotID := id != app.InvalidSetID
			gotErr := err != nil

			if tc.wantID != gotID {
				t.Fatalf("want id: %v\ngot id: %v\nid: %v", tc.wantID, gotID, id)
			}

			if tc.wantErr != gotErr {
				t.Fatalf("want err: %v\ngot err: %v\nerr: %v", tc.wantErr, gotErr, err)
			}
		})
	}
}

// newValidSetStore returns a SetStore whose database operations perform as expected.
// NOTE: for testing only
func newValidSetStore(h *mongoHandle) app.SetStore {
	return NewSetStore(h)
}

// newInvalidSetStore returns a SetStore whose underlying context has been canceled, thus
// all database operations will fail.
// NOTE: for testing only
func newInvalidSetStore(h *mongoHandle) app.SetStore {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	return &setStore{
		coll: h.db.Collection(TestDB),
		ctx:  cancelCtx,
	}
}
