package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
)

func testOptions() *Options {
	o := NewOptions()
	o.Upstreams = append(o.Upstreams, "http://127.0.0.1:8080/")
	o.CookieSecret = "foobar"
	o.ClientID = "bazquux"
	o.ClientSecret = "xyzzyplugh"
	return o
}

func errorMsg(msgs []string) string {
	result := make([]string, 0)
	result = append(result, "Invalid configuration:")
	result = append(result, msgs...)
	return strings.Join(result, "\n  ")
}

func TestNewOptions(t *testing.T) {
	o := NewOptions()
	err := o.Validate()
	assert.NotEqual(t, nil, err)

	expected := errorMsg([]string{
		"missing setting: upstream",
		"missing setting: cookie-secret",
		"missing setting: client-id",
		"missing setting: client-secret"})
	assert.Equal(t, expected, err.Error())
}

func TestInitializedOptions(t *testing.T) {
	o := testOptions()
	assert.Equal(t, nil, o.Validate())
}

// Note that it's not worth testing nonparseable URLs, since url.Parse()
// seems to parse damn near anything.
func TestRedirectUrl(t *testing.T) {
	o := testOptions()
	o.RedirectUrl = "https://myhost.com/oauth2/callback"
	assert.Equal(t, nil, o.Validate())
	expected := &url.URL{
		Scheme: "https", Host: "myhost.com", Path: "/oauth2/callback"}
	assert.Equal(t, expected, o.redirectUrl)
}

func TestProxyUrls(t *testing.T) {
	o := testOptions()
	o.Upstreams = append(o.Upstreams, "http://127.0.0.1:8081")
	assert.Equal(t, nil, o.Validate())
	expected := []*url.URL{
		&url.URL{Scheme: "http", Host: "127.0.0.1:8080", Path: "/"},
		// note the '/' was added
		&url.URL{Scheme: "http", Host: "127.0.0.1:8081", Path: "/"},
	}
	assert.Equal(t, expected, o.proxyUrls)
}

func TestCompiledRegex(t *testing.T) {
	o := testOptions()
	regexps := []string{"/foo/.*", "/ba[rz]/quux"}
	o.SkipAuthRegex = regexps
	assert.Equal(t, nil, o.Validate())
	actual := make([]string, 0)
	for _, regex := range o.CompiledRegex {
		actual = append(actual, regex.String())
	}
	assert.Equal(t, regexps, actual)
}

func TestCompiledRegexError(t *testing.T) {
	o := testOptions()
	o.SkipAuthRegex = []string{"(foobaz", "barquux)"}
	err := o.Validate()
	assert.NotEqual(t, nil, err)

	expected := errorMsg([]string{
		"error compiling regex=\"(foobaz\" error parsing regexp: " +
			"missing closing ): `(foobaz`",
		"error compiling regex=\"barquux)\" error parsing regexp: " +
			"unexpected ): `barquux)`"})
	assert.Equal(t, expected, err.Error())
}
