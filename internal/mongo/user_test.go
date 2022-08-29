package mongo

import (
	"context"
	"reflect"
	"testing"

	"github.com/hrand1005/training-notebook/internal/app"
)

// TestInsert tests the mongo implementation of the UserStore's Insert() method.
func TestInsert(t *testing.T) {
	handle, err := New(TestURI, TestDB)
	if err != nil {
		t.Fatalf("failed to initialize test db handle: %v", err)
	}
	defer handle.Delete()
	defer handle.Close()

	validStore := newValidUserStore(handle)
	invalidStore := newInvalidUserStore(handle)

	tests := []struct {
		name    string
		store   app.UserStore
		user    *app.User
		wantID  bool
		wantErr bool
	}{
		{
			name:  "Nominal case returns ID and nil error",
			store: validStore,
			user: &app.User{
				FirstName:    "yorbus",
				LastName:     "bonk",
				Email:        "ybonk@apple.mail",
				PasswordHash: "ThIsIsApAsSwOrDhAsH666",
			},
			wantID:  true,
			wantErr: false,
		},
		{
			name:  "Invalid user store case returns empty ID and error",
			store: invalidStore,
			user: &app.User{
				FirstName:    "yorbus",
				LastName:     "bonk",
				Email:        "ybonk@apple.mail",
				PasswordHash: "ThIsIsApAsSwOrDhAsH666",
			},
			wantID:  false,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, err := tc.store.Insert(tc.user)
			gotID := id != ""
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

// TestFindByID tests the mongo implementation of the UserStore's FindByID() method.
func TestFindByID(t *testing.T) {
	handle, err := New(TestURI, TestDB)
	if err != nil {
		t.Fatalf("failed to initialize test db handle: %v", err)
	}
	defer handle.Delete()
	defer handle.Close()

	testUser := &app.User{
		FirstName: "test-first-name",
		LastName:  "test-last-name",
		Email:     "test-email@yahoo.mail",
	}

	validStore := newValidUserStore(handle)
	testUserID, err := validStore.Insert(testUser)
	if err != nil {
		t.Fatalf("failed to initialize test user: %v", err)
	}

	testUser.ID = testUserID

	invalidStore := newInvalidUserStore(handle)

	tests := []struct {
		name     string
		store    app.UserStore
		id       app.UserID
		wantUser *app.User
		wantErr  bool
	}{
		{
			name:     "Nominal success returns user found with id",
			store:    validStore,
			id:       testUserID,
			wantUser: testUser,
			wantErr:  false,
		},
		{
			name:     "Invalid id returns nil user and error",
			store:    invalidStore,
			id:       "invalid-id",
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid store returns nil user and error",
			store:    invalidStore,
			id:       testUserID,
			wantUser: nil,
			wantErr:  true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotUser, err := tc.store.FindByID(tc.id)
			gotErr := err != nil

			if !reflect.DeepEqual(tc.wantUser, gotUser) {
				t.Fatalf("want user: %#v\ngot user: %#v\nid: %v\nerr: %v", tc.wantUser, gotUser, tc.id, err)
			}

			if tc.wantErr != gotErr {
				t.Fatalf("want err: %v\ngot err: %v\nerr: %v", tc.wantErr, gotErr, err)
			}
		})
	}
}

// newValidUserStore returns a UserStore whose database operations perform as expected.
// NOTE: for testing only
func newValidUserStore(h *mongoHandle) app.UserStore {
	return NewUserStore(h)
}

// newInvalidUserStore returns a UserStore whose underlying context has been canceled, thus
// all database operations will fail.
// NOTE: for testing only
func newInvalidUserStore(h *mongoHandle) app.UserStore {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	return &userStore{
		coll: h.db.Collection(TestDB),
		ctx:  cancelCtx,
	}
}
