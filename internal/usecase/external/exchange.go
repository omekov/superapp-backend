package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//	{
//		"disclaimer": "Usage subject to terms: https://openexchangerates.org/terms",
//		"license": "https://openexchangerates.org/license",
//		"timestamp": 1726225173,
//		"base": "USD",
//		"rates": {
//		"AED": 3.67298,
//		"CNY": 7.0944,
//		"KZT": 479.662757,
//		"RUB": 90.795092
//		}
//	}
type OpenExchangeRatesResponse struct {
	Disclaimer string   `json:"disclaimer"`
	License    string   `json:"license"`
	Timestamp  int      `json:"timestamp"`
	Base       int      `json:"base"`
	Rates      Currency `json:"rate"`
}

type Currency struct {
	AED float64 `json:"AED"`
	CNY float64 `json:"CNY"`
	KZT float64 `json:"KZT"`
	RUB float64 `json:"RUB"`
}

// {
// "error": true,
// "status": 403,
// "message": "not_allowed",
// "description": "Changing the API `base` currency is available for Developer, Enterprise and Unlimited plan clients. Please upgrade, or contact support@openexchangerates.org with any questions."
// }
type OpenExchangeRatesErrorResponse struct {
	Error       bool   `json:"error"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (c Client) GetCurrency(ctx context.Context) (OpenExchangeRatesResponse, error) {
	var bodyResp OpenExchangeRatesResponse
	resp, err := c.httpClient.Get(c.exchangeURL)
	if err != nil {
		return bodyResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var bodyErrResp OpenExchangeRatesErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&bodyErrResp); err != nil {
			return bodyResp, fmt.Errorf("json.NewDecoder.Decode status: %d, err: %s", resp.StatusCode, err.Error())
		}
		return bodyResp, fmt.Errorf("ошибка загрузки файла: статус %d, err: %s", resp.StatusCode, bodyErrResp.Description)
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return bodyResp, fmt.Errorf("Ошибка при чтении данных файла: %v", err)
		}

		if err := json.Unmarshal(bodyBytes, &bodyResp); err != nil {
			return bodyResp, fmt.Errorf("Ошибка json.Unmarshal: %v", err)
		}
		return bodyResp, nil
	}
}
