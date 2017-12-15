package main

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirectHandler(t *testing.T) {

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	redirectHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 302, resp.StatusCode, "Should return a 302 response")
	assert.Equal(t, "<a href=\"https://www.youtube.com/embed/dQw4w9WgXcQ\">Found</a>.\n\n", string(body), "Should return an anchor tag in the body")
	fmt.Println(resp.StatusCode)
}
