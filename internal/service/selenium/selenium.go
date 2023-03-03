package selenium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"time"
)

type SeleniumClient struct {
	URL string
}

type GetScreenshotByURLPayload struct {
	URL string
}

type GetScreenshotByURLResponse struct {
	FileName string `json:"filename"`
}

func CreateSeleniumClient(url string) *SeleniumClient {
	return &SeleniumClient{
		URL: url,
	}
}

func (sc *SeleniumClient) GetURL() string {
	return sc.URL
}

func (sc *SeleniumClient) GetScreenshotByURL(method string, path string, pl *GetScreenshotByURLPayload) (*GetScreenshotByURLResponse, error) {
	if method == "" {
		method = "POST"
	}
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("url", pl.URL)

	err := writer.Close()
	if err != nil {
		panic(err)

		return nil, err
	}
	client := &http.Client{Timeout: time.Second * 10} // 設置超時時間
	url := fmt.Sprintf("%s/%s", sc.URL, path)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// 網路超時
			return nil, fmt.Errorf("net timeout: %v", err)
		}
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var response GetScreenshotByURLResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		panic(err)

		return nil, err
	}

	return &response, nil
}
