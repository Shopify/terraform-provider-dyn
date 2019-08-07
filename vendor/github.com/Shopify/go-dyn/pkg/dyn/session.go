package dyn

type sessionLogInRequest struct {
	CustomerName string `json:"customer_name"`
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
}

type sessionLoginResponseData struct {
	Token   string `json:"token"`
	Version string `json:"version"`
}

type sessionLogInResponse struct {
	responseHeader
	sessionLoginResponseData `json:"data"`
}

// LogIn establishes an API session.
func (c *Client) LogIn(customerName, userName, password string) error {
	req := sessionLogInRequest{
		CustomerName: customerName,
		UserName:     userName,
		Password:     password,
	}

	var resp sessionLogInResponse

	if err := c.post("Session", req, &resp); err != nil {
		return err
	}

	c.token = resp.Token

	return nil
}

// KeepAlive keeps an API session alive.
func (c *Client) KeepAlive() error {
	return c.put("Session", nil, nil)
}

// IsActive verifies that an API session is alive.
func (c *Client) IsActive() (bool, error) {
	if err := c.get("Session", nil, nil); err != nil {
		return false, err
	}

	return true, nil
}

// LogOut ends an API session.
func (c *Client) LogOut() error {
	if err := c.delete("Session", nil); err != nil {
		return err
	}

	c.token = ""

	return nil
}
