package fastforex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type HTTPRepository struct {
	apiKey string
	client *resty.Client
}

func NewHTTPRepository(apiKey string) *HTTPRepository {
	return &HTTPRepository{apiKey: apiKey, client: resty.New()}
}

type fetchOneResult struct {
	Base    string             `json:"base"`
	Result  map[string]float64 `json:"result"`
	Updated string             `json:"updated"`
	Ms      int                `json:"ms"`
}

func (r *HTTPRepository) FetchOne(ctx context.Context, from, to string) (float64, error) {
	url := fmt.Sprintf("https://api.fastforex.io/fetch-one?from=%s&to=%s&api_key=%s", from, to, r.apiKey)
	rawResp, err := r.client.
		R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return 0, errors.Wrap(err, utils.GetFuncName(r.client.R().Get))
	}

	if rawResp.StatusCode() != 200 {
		return 0, errors.Wrap(
			errors.New(fmt.Sprintf("fast forex status code error: %d", rawResp.StatusCode())),
			utils.GetFuncName(r.client.R().Get))
	}

	var resp fetchOneResult
	if err = json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return 0, errors.Wrap(err, utils.GetFuncName(json.Unmarshal))
	}

	rate, ok := resp.Result[to]
	if !ok {
		return 0, errors.New(fmt.Sprintf("currency not supported: %s", to))
	}

	return rate, nil
}
