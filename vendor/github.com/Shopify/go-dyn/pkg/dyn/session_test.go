package dyn

import (
	"net/http"
	"testing"
)

func TestSessionLogIn(t *testing.T) {
	c := mockClient("session/log_in.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPost, r)
		assertPath(t, "/REST/Session", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)

		assertJSON(t, "insert-customer-here", "customer_name", j)
		assertJSON(t, "insert-user-here", "user_name", j)
		assertJSON(t, "insert-password-here", "password", j)

		w.Header().Set("Content-Type", "application/json")
	})

	if err := c.LogIn("insert-customer-here", "insert-user-here", "insert-password-here"); err != nil {
		t.Error(err)
	} else {
		assertEqual(t, "insert-token-here", c.token, "session token")
	}
}

func TestSessionLogInError(t *testing.T) {
	c := mockClient("session/log_in_error.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)
	})

	if err := c.LogIn("insert-customer-here", "insert-user-here", "insert-password-here"); err != nil {
		assertEqual(t, "INVALID_DATA: login: Invalid credentials", err.Error(), "error")
	} else {
		t.Error("Expected LogIn to fail")
	}
}

func TestSessionIsActive(t *testing.T) {
	c := mockClient("session/active.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodGet, r)
		assertPath(t, "/REST/Session", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if ok, err := c.IsActive(); err != nil {
		t.Error(err)
	} else {
		assertEqual(t, true, ok, "IsActive")
	}
}

func TestSessionIsActiveError(t *testing.T) {
	c := mockClient("session/active_error.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)
	})

	if ok, err := c.IsActive(); err != nil {
		assertEqual(t, false, ok, "IsActive")
		assertEqual(t, "INVALID_DATA: login: Bad or expired credentials", err.Error(), "error")
	} else {
		t.Error("Expected IsActive to fail")
	}
}

func TestSessionKeepAlive(t *testing.T) {
	c := mockClient("session/keep_alive.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPut, r)
		assertPath(t, "/REST/Session", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if err := c.KeepAlive(); err != nil {
		t.Error(err)
	}
}

func TestSessionKeepAliveError(t *testing.T) {
	c := mockClient("session/keep_alive_error.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)
	})

	if err := c.KeepAlive(); err != nil {
		assertEqual(t, "INVALID_DATA: login: Bad or expired credentials", err.Error(), "error")
	} else {
		t.Error("Expected KeepAlive to fail")
	}
}

func TestSessionLogOut(t *testing.T) {
	c := mockClient("session/log_out.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodDelete, r)
		assertPath(t, "/REST/Session", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if err := c.LogOut(); err != nil {
		t.Error(err)
	}
}

func TestSessionLogOutError(t *testing.T) {
	c := mockClient("session/log_out_error.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)
	})

	if err := c.LogOut(); err != nil {
		assertEqual(t, "INVALID_DATA: login: Bad or expired credentials", err.Error(), "error")
	} else {
		t.Error("Expected LogOut to fail")
	}
}
