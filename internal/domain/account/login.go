package account

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/lib/system"
)

//=================================================================================================
// LOGIN COMMAND
//=================================================================================================

type LoginCommand struct {
	Server  string
	Runtime system.Runtime
	Keyring system.Keyring
	Timeout time.Duration
}

func Login(cmd *LoginCommand) (*User, error) {
	return cmd.execute()
}

//=================================================================================================
// PRIVATE IMPLEMENTATION
//=================================================================================================

type loginServer struct {
	Port       int
	JwtChannel chan string
	ErrChannel chan error
	Stop       func() error
}

type loginTimer struct {
	Context context.Context
	Cancel  context.CancelFunc
}

//-------------------------------------------------------------------------------------------------

func (cmd *LoginCommand) execute() (*User, error) {

	if cmd.Server == "" {
		return nil, fmt.Errorf("missing server")
	} else if cmd.Runtime == nil {
		return nil, fmt.Errorf("missing runtime")
	} else if cmd.Keyring == nil {
		return nil, fmt.Errorf("missing keyring")
	}

	if cmd.Timeout == 0 {
		cmd.Timeout = 2 * time.Minute
	}

	jwt, ok := cmd.Keyring.Get(httpx.ParamJWT)
	if ok {
		user, err := cmd.validate(jwt)
		if err == nil {
			return user, nil
		} else {
			cmd.Keyring.Del(httpx.ParamJWT) // delete invalid JWT and continue
		}
	}

	timer := cmd.startTimer()
	defer timer.Cancel()

	server, err := cmd.startServer()
	if err != nil {
		return nil, err
	}
	defer server.Stop()

	cmd.launchBrowser(server.Port)

	jwt, err = cmd.waitForLogin(server, timer)
	if err != nil {
		return nil, err
	}

	user, err := cmd.validate(jwt)
	if err != nil {
		return nil, err
	}

	err = cmd.Keyring.Set(httpx.ParamJWT, jwt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//-------------------------------------------------------------------------------------------------

func (cmd *LoginCommand) startServer() (*loginServer, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	jwtChannel := make(chan string, 1)
	errChannel := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		jwt := r.FormValue("jwt")
		if jwt == "" {
			http.Error(w, "Missing JWT", http.StatusBadRequest)
			errChannel <- fmt.Errorf("missing jwt in callback")
			return
		}
		jwtChannel <- jwt
		w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeHTML)
		fmt.Fprintln(w, loginSuccessPage)
	})

	port := listener.Addr().(*net.TCPAddr).Port

	server := &http.Server{Handler: mux}
	go server.Serve(listener)

	stop := func() error {
		err := listener.Close()
		if err != nil {
			return err
		}
		err = server.Shutdown(context.Background())
		if err != nil {
			return err
		}
		return nil
	}

	return &loginServer{
		Port:       port,
		JwtChannel: jwtChannel,
		ErrChannel: errChannel,
		Stop:       stop,
	}, nil
}

//-------------------------------------------------------------------------------------------------

const loginSuccessPage string = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
	<title>Login Complete</title>
	<style>
	  body {
		  font-family: sans-serif;
			text-align: center;
			padding-top: 50px;
		}
	</style>
	<script>
	history.replaceState(null, '', location.pathname)
	</script>
</head>
<body>
  <h2>Login Successful</h2>
	<p>You can now close this window</p>
</body>
</html>`

//-------------------------------------------------------------------------------------------------

func (cmd *LoginCommand) startTimer() *loginTimer {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
	return &loginTimer{
		Context: ctx,
		Cancel:  cancel,
	}
}

//-------------------------------------------------------------------------------------------------

func (cmd *LoginCommand) launchBrowser(port int) {
	url, _ := url.Parse(cmd.Server)
	url.Path = "login"
	q := url.Query()
	q.Set(httpx.ParamCLI, "true")
	q.Set(httpx.ParamOrigin, fmt.Sprintf("http://127.0.0.1:%d/callback", port))
	url.RawQuery = q.Encode()
	cmd.Runtime.Open(url.String())
}

//-------------------------------------------------------------------------------------------------

func (cmd *LoginCommand) waitForLogin(server *loginServer, timer *loginTimer) (string, error) {
	select {
	case jwt := <-server.JwtChannel:
		return jwt, nil
	case err := <-server.ErrChannel:
		return "", err
	case <-timer.Context.Done():
		return "", fmt.Errorf("login timed out")
	}
}

//-------------------------------------------------------------------------------------------------

func (cmd *LoginCommand) validate(jwt string) (*User, error) {
	url, _ := url.Parse(cmd.Server)
	url.Path = "api/account/me"

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unexpected request: %s", err)
	}
	req.Header.Set(httpx.HeaderAuthorization, fmt.Sprintf("Bearer %s", jwt))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unexpected response: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("unexpected JSON response: %s", err)
	}

	return &user, nil
}

//-------------------------------------------------------------------------------------------------
