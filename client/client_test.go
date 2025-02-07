package client

//
//func TestMakeRequest(t *testing.T) {
//	serverURL := testutil.Setup(testutil.MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
//	defer testutil.Teardown()
//	c, _ := NewUserPasswordConfig(serverURL, "foo", "bar")
//	c.WithSkipTLSValidation(true)
//	client, err := New(c)
//	require.NoError(t, err)
//	req := client.NewRequest("GET", "/v2/organizations")
//	resp, err := client.DoRequest(req)
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//}
//
//func TestMakeRequestFailure(t *testing.T) {
//	serverURL := testutil.Setup(testutil.MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
//	defer testutil.Teardown()
//	c, _ := NewUserPasswordConfig(serverURL, "foo", "bar")
//	c.WithSkipTLSValidation(true)
//	client, err := New(c)
//	require.NoError(t, err)
//	req := client.NewRequest("GET", "/v2/organizations")
//	req.url = "%gh&%ij"
//	resp, err := client.DoRequest(req)
//	require.Nil(t, resp)
//	require.NotNil(t, err)
//}
//
//func TestMakeRequestWithTimeout(t *testing.T) {
//	serverURL := testutil.Setup(testutil.MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
//	defer testutil.Teardown()
//	c, _ := NewUserPasswordConfig(serverURL, "foo", "bar")
//	c.WithSkipTLSValidation(true)
//	c.WithHTTPClient(&http.Client{Timeout: 10 * time.Nanosecond})
//	client, err := New(c)
//	require.NotNil(t, err)
//	require.Nil(t, client)
//}
//
//func TestHTTPErrorHandling(t *testing.T) {
//	serverURL := testutil.Setup(testutil.MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"502 Bad Gateway"}, Status: 502}, t)
//	defer testutil.Teardown()
//	c, _ := NewUserPasswordConfig(serverURL, "foo", "bar")
//	c.WithSkipTLSValidation(true)
//	client, err := New(c)
//	require.NoError(t, err)
//	req := client.NewRequest("GET", "/v2/organizations")
//	resp, err := client.DoRequest(req)
//	require.NotNil(t, err)
//	require.NotNil(t, resp)
//
//	httpErr := err.(CloudFoundryHTTPError)
//	require.Equal(t, 502, httpErr.StatusCode)
//	require.Equal(t, "502 Bad Gateway", httpErr.Status)
//	require.Equal(t, "502 Bad Gateway", string(httpErr.Body))
//}
//
//func TestTokenRefresh(t *testing.T) {
//	serverURL := testutil.Setup(testutil.MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
//	testutil.SetupFakeUAAServer(1)
//	defer testutil.Teardown()
//
//	c, _ := NewUserPasswordConfig(serverURL, "foo", "bar")
//	client, err := New(c)
//	require.NoError(t, err)
//
//	token, err := client.GetToken()
//	require.NoError(t, err)
//	require.Equal(t, "bearer foobar2", token)
//
//	for i := 0; i < 5; i++ {
//		token, _ = client.GetToken()
//		if token == "bearer foobar3" {
//			break
//		}
//		time.Sleep(time.Second)
//	}
//	require.Equal(t, "bearer foobar3", token)
//}
//
//func TestEndpointRefresh(t *testing.T) {
//	serverURL := testutil.Setup(testutil.MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
//	testutil.SetupFakeUAAServer(0)
//	defer testutil.Teardown()
//
//	c, _ := NewUserPasswordConfig(serverURL, "foo", "bar")
//	client, err := New(c)
//	require.NoError(t, err)
//
//	//lastTokenSource := client.config.TokenSource
//	for i := 1; i < 5; i++ {
//		_, err := client.GetToken()
//		require.NoError(t, err)
//		//So(client.config.TokenSource, ShouldNotEqual, lastTokenSource)
//		//lastTokenSource = client.config.TokenSource
//	}
//}
