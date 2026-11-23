package tcgdex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"
)

func NewTCGDex() *TCGDex {
	api := new(TCGDex)
	api.client = &http.Client{
		Timeout: 10 * time.Second,
	}
	api.baseURL = "api.tcgdex.net/v2/en"
	return api
}

func (td *TCGDex) NewGetRequest(parts ...string) (*http.Request, error) {
	parts = append([]string{td.baseURL}, parts...)
	url := path.Join(parts...)
	req, err := http.NewRequest(http.MethodGet, "https://"+url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error requesting sets from TCGDex: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	return req, err
}

func (td *TCGDex) DecodeJSON(res *http.Response, body any) error {
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&body); err != nil {
		return fmt.Errorf("JSON decoding error %w", err)
	}
	return nil
}

func (td *TCGDex) GetSeries() ([]TCGDexSerieBrief, error) {
	req, err := td.NewGetRequest("series")
	if err != nil {
		return nil, err
	}
	res, err := td.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response error: %d - %w", res.StatusCode, err)
	}

	var body []TCGDexSerieBrief
	if err := td.DecodeJSON(res, &body); err != nil {
		return nil, fmt.Errorf("JSON decoding error %w", err)
	}
	return body, nil
}

func (td *TCGDex) GetSets() ([]TCGDexSetBrief, error) {
	url := path.Join(td.baseURL, "sets")
	req, err := http.NewRequest(http.MethodGet, "https://"+url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error requesting sets from TCGDex: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := td.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request Failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response error: %d - %w", res.StatusCode, err)
	}

	decoder := json.NewDecoder(res.Body)
	var body []TCGDexSetBrief
	if err := decoder.Decode(&body); err != nil {
		return nil, fmt.Errorf("JSON Decoding error %w", err)
	}
	return body, nil
}

func (td *TCGDex) GetSet(id string) (TCGDexSet, error) {
	url := path.Join(td.baseURL, "sets", id)
	req, err := http.NewRequest(http.MethodGet, "https://"+url, nil)
	none := TCGDexSet{}
	if err != nil {
		return none, fmt.Errorf("Error requesting a set from TCGDex: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := td.client.Do(req)
	if err != nil {
		return none, fmt.Errorf("Request Failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return none, fmt.Errorf("Response error: %d :  %w", res.StatusCode, err)
	}

	var body TCGDexSet
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&body); err != nil {
		return none, fmt.Errorf("JSON Decoding error: %w", err)
	}
	return body, nil
}

func (td *TCGDex) GetCard(id string) ([]byte, error) {
	req, err := td.NewGetRequest("cards", id)
	if err != nil {
		return nil, fmt.Errorf("malformed request: %w", err)
	}
	res, err := td.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response error: %d - %w", res.StatusCode, err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not write response into buffer %w", err)
	}

	return buf.Bytes(), nil
}
