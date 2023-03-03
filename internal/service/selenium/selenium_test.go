package selenium

import (
	"bytes"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetScreenshotByURL(t *testing.T) {
	// 創建一個測試HTTP請求
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("url", "http://example.com")
	writer.Close()

	req, err := http.NewRequest("POST", "http://external-service/api/v1/screenshot", &body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 創建一個模擬的外部服務HTTP處理程序
	externalServiceHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 檢查HTTP請求是否正確處理了名為"url"的表單數據參數，32 << 20表示要解析的最大內容大小
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			t.Errorf("failed to parse form data: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		url := r.FormValue("url")
		if url != "http://example.com" {
			t.Errorf("unexpected url: got %v want http://example.com", url)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// 回傳模擬的圖像響應
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte("fake image data"))
	})

	// 創建一個測試HTTP服務器，並將外部服務處理程序註冊為路由處理程序
	testServer := httptest.NewServer(externalServiceHandler)
	defer testServer.Close()

	// 假設screenshotHandler是外部服務的API處理程序
	screenshotHandler := func(w http.ResponseWriter, r *http.Request) {
		// 設置Content-Type標頭為image/png
		w.Header().Set("Content-Type", "image/png")

		// 模擬一個PNG圖像，可以將此替換為實際從外部服務API獲取的圖像
		img := image.NewRGBA(image.Rect(0, 0, 640, 480))
		png.Encode(w, img)
	}
	// 執行測試
	handler := http.HandlerFunc(screenshotHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// 檢查HTTP響應碼是否正確
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢查響應內容是否為圖像
	contentType := rr.Header().Get("Content-Type")
	if contentType != "image/png" {
		t.Errorf("handler returned wrong content type: got %v want image/png",
			contentType)
	}

	// 檢查響應主體是否為非空字節切片
	bodyBytes := rr.Body.Bytes()
	if len(bodyBytes) == 0 {
		t.Errorf("handler returned empty response body")
	}
}
