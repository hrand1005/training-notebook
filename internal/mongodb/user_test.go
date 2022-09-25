package mongodb

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/hrand1005/training-notebook/internal/app"
)

// TestInsertUser tests the mongo implementation of the UserStore's Insert() method.
func TestInsertUser(t *testing.T) {
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
			gotID := id != app.InvalidUserID
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

// TestUpdateByID tests the mongo implementation of the UserStore's UpdateBYID() method.
func TestUpdateByID(t *testing.T) {
	handle, err := New(TestURI, TestDB)
	if err != nil {
		t.Fatalf("failed to initialize test db handle: %v", err)
	}
	defer handle.Delete()
	defer handle.Close()

	validStore := newValidUserStore(handle)
	invalidStore := newInvalidUserStore(handle)

	tests := []struct {
		name        string
		store       app.UserStore
		userIDValid bool
		user        *app.User
		wantErr     error
	}{
		{
			name:        "Nominal case returns nil error",
			store:       validStore,
			userIDValid: true,
			user: &app.User{
				FirstName: "yorbus",
			},
			wantErr: nil,
		},
		{
			name:        "Invalid user id case returns error",
			store:       validStore,
			userIDValid: false,
			user: &app.User{
				FirstName: "yorbus",
			},
			wantErr: app.ErrNotFound,
		},
		{
			name:        "Invalid user store case returns error",
			store:       invalidStore,
			userIDValid: true,
			user: &app.User{
				FirstName: "yorbus",
			},
			wantErr: app.ErrServiceFailure,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// add empty user
			id, _ := tc.store.Insert(&app.User{})

			if !tc.userIDValid {
				id = app.InvalidUserID
			}

			// attempt to update the inserted user
			gotErr := tc.store.UpdateByID(id, tc.user)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("want err: %v\ngot err: %v", tc.wantErr, gotErr)
			}
		})
	}
}

// TestDeleteByID tests the mongo implementation of the UserStore's DeleteByID() method.
func TestDeleteByID(t *testing.T) {
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
		name    string
		store   app.UserStore
		userID  app.UserID
		wantErr error
	}{
		{
			name:    "Nominal case returns nil error",
			store:   validStore,
			userID:  testUserID,
			wantErr: nil,
		},
		{
			name:    "Invalid user id case returns error",
			store:   validStore,
			userID:  "invalid-id",
			wantErr: app.ErrNotFound,
		},
		{
			name:    "Invalid user store case returns error",
			store:   invalidStore,
			userID:  testUserID,
			wantErr: app.ErrServiceFailure,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := tc.store.DeleteByID(tc.userID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("want err: %v\ngot err: %v", tc.wantErr, gotErr)
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
