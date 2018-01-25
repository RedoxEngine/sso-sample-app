# sso-sample-app
Sample implementation of receiving SSO from Redox

A production application would use some of session storage to save the JWT/Post Body, and then redirect using a secure token. 

# Setup and dependencies
[Install Go](https://golang.org/doc/install)

# Environment Variables
`SSO_SECRET` This is the Secret generated in the Redox dashboard
`PORT` The port you want the web server to run on
`BASE_URL` The url of the application. This makes redirects work, defaults to localhost:$PORT

# Commands
`go run ./sso.go` will compile the app and run
`go test` runs the tests
