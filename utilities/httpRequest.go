package utilities

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func ReadJsonHttpBody(r *http.Request, i interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New("Unable to read http request body: " + err.Error())
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return errors.New("Unable to deserialize request body: " + err.Error())
	}
	return nil
}

func ReadAuthenticationHeader(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) == 0 {
		return "", errors.New("Could not get bearer token")
	}
	reqToken = splitToken[1]
	return reqToken, nil
}

func ReadStringQueryParameter(r *http.Request, key string) (string, error) {
	s, ok := r.URL.Query()[key]
	if s[0] == "" || !ok {
		msg := key + " query parameter missing"
		return s[0], errors.New(msg)
	}
	return s[0], nil
}

func ReadIntQueryParameter(r *http.Request, key string) (int, error) {
	s, ok := r.URL.Query()[key]
	var si int
	if !ok || s[0] == "" {
		msg := key + " query parameter missing"
		return si, errors.New(msg)
	}
	si, err := strconv.Atoi(s[0])
	if err != nil {
		return si, errors.New(err.Error())
	}
	return si, nil
}

func ReadBooleanQueryParameter(r *http.Request, key string) (bool, error) {
	v, ok := r.URL.Query()[key]
	var vb bool
	if !ok || v[0] == "" {
		msg := key + " query parameter missing"
		return vb, errors.New(msg)
	}
	vb, err := strconv.ParseBool(v[0])
	if err != nil {
		return vb, errors.New(err.Error())
	}
	return vb, nil
}
