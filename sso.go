package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var templ *template.Template
var baseURL = string(os.Getenv("BASE_URL"))
var port = os.Getenv("PORT")

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

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {
	homeTemplate, err := ioutil.ReadFile("./templates/index.tmpl")
	check(err)
	templ = template.Must(template.New("home").Parse(string(homeTemplate)))
	if baseURL == "" {
		baseURL = "http://localhost:"
	}

	startServer()
}

func startServer() {

	if port == "" {
		log.Fatal("$PORT must be set")
	}

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
	r.HandleFunc("/home", homeHandler)

	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}

// redirectHandler takes an HTTP request and redirects to the main page
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user").(*jwt.Token)
	providerName := token.Claims.(jwt.MapClaims)["name"].(string)
	http.Redirect(w, r, baseURL+"/home?auth="+providerName, 302)
}

// homeHandler renders the html template with the query string parameters
func homeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	data := struct {
		Title string
	}{
		Title: params.Get("auth"),
	}

	w.WriteHeader(http.StatusOK)

	err := templ.Execute(w, data)
	check(err)
}
