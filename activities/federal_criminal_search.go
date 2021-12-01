package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/temporalio/background-checks/types"
)

func (a *Activities) FederalCriminalSearch(ctx context.Context, input types.FederalCriminalSearchInput) (types.FederalCriminalSearchResult, error) {
	var result types.FederalCriminalSearchResult

	requestURL, err := url.Parse("http://thirdparty:8082/federalcriminalsearch")
	if err != nil {
		return result, err
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return result, fmt.Errorf("unable to encode input: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL.String(), bytes.NewBuffer(jsonInput))
	if err != nil {
		return result, fmt.Errorf("unable to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		return result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return result, err
}