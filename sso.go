package main

import (
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// RedoxIds represents patient Identifiers
type RedoxIds struct {
	ID     string `json:"id,omitempty"`
	IDType string `json:"id_type,omitempty"`
}

// RedoxCustomClaims are the custom properties in the Redox SSO JSON Web Token
type RedoxCustomClaims struct {
	Name         string     `json:"name,omitempty"`
	FirstName    string     `json:"given_name,omitempty"`
	LastName     string     `json:"family_name,omitempty"`
	MiddleName   string     `json:"middle_name,omitempty"`
	EmailAddress string     `json:"email,omitempty"`
	NPI          string     `json:"npi,omitempty"`
	Patient      []RedoxIds `json:"patient_ids,omitempty"`
	VisitNumber  string     `json:"visit_id,omitempty"`
	Facility     string     `json:"facility_id,omitempty"`
	Department   string     `json:"department_id,omitempty"`
	TimeZone     string     `json:"zoneinfo,omitempty"`
	Locale       string     `json:"locale,omitempty"`
	PhoneNumber  string     `json:"phone_number,omitempty"`
	jwt.StandardClaims
}

func main() {

	startServer()
}

func startServer() {
	r := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SSO_SECRET")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r.HandleFunc("/public", redirectHandler)
	r.Handle("/secure", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(redirectHandler)),
	))

	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}

// redirectHandler takes an HTTP request and redirects to a special place
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "https://www.youtube.com/embed/dQw4w9WgXcQ", 302)
}
