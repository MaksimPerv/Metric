# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).



package main

import (
"github.com/stretchr/testify/assert"
"io"
"net/http"
"net/http/httptest"
"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
req, err := http.NewRequest(method, ts.URL+path, nil)
req.Header.Add("Content-Type", "text/plain")
assert.NoError(t, err)
resp, err := ts.Client().Do(req)
assert.NoError(t, err)
respBody, err := io.ReadAll(resp.Body)
resp.Body.Close()
assert.NoError(t, err)
return resp, string(respBody)
}

func TestServer(t *testing.T) {
type want struct {
statusCode int
statusType string
body       string
}
ts := httptest.NewServer(run())
defer ts.Close()

	tests := []struct {
		url    string
		method string
		want   want
	}{
		{
			url:    "/update/gauge/asd/123",
			method: "POST",
			want: want{
				statusCode: 200,
				statusType: "text/plain; charset=utf-8",
				body:       "Metric updated\n",
			},
		},
		{
			url:    "/update/gauge/asd/123",
			method: "GET",
			want: want{
				statusCode: 405,
			},
		},
	}

	for _, test := range tests {
		resp, respBody := testRequest(t, ts, test.method, test.url)
		assert.Equal(t, resp.StatusCode, test.want.statusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), test.want.statusType)
		assert.Equal(t, respBody, test.want.body)
	}
}
