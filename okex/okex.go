package okex

import (
	"net/url"

	"strings"

	"strconv"

	"github.com/leek-box/sheep/consts"
	"github.com/leek-box/sheep/proto"
	"github.com/pkg/errors"
)

type OKEX struct {
	accessKey string
	secretKey string
}

func (o *OKEX) GetExchangeType() string {
	return consts.ExchangeTypeOKEX
}

func (o *OKEX) GetAccountBalance() ([]proto.AccountBalance, error) {
	path := "userinfo.do"
	var ret BalanceReturn
	err := o.apiKeyPost(url.Values{}, path, &ret)
	if err != nil {
		return nil, err
	}

	if !ret.Result {
		return nil, errors.New("result error")
	}

	var res []proto.AccountBalance
	for k, v := range ret.Info.Funds.Free {
		if v != "0" {
			var item proto.AccountBalance
			item.Currency = k
			item.Balance = v
			item.Type = proto.AccountBalanceTypeTrade

			res = append(res, item)
		}
	}

	for k, v := range ret.Info.Funds.Freezed {
		if v != "0" {
			var item proto.AccountBalance
			item.Currency = k
			item.Balance = v
			item.Type = proto.AccountBalanceTypeFrozen

			res = append(res, item)
		}
	}

	return res, nil
}

//访问频率 20次/2秒
func (o *OKEX) OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error) {
	path := "trade.do"

	values := url.Values{}
	values.Set("symbol", strings.ToLower(params.BaseCurrencyID)+"_"+strings.ToLower(params.QuoteCurrencyID))
	values.Set("type", TransOrderType(params.Type))
	values.Set("price", strconv.FormatFloat(params.Price, 'f', -1, 64))
	values.Set("amount", strconv.FormatFloat(params.Amount, 'f', -1, 64))

	var okRet OrderPlaceReturn
	err := o.apiKeyPost(values, path, &okRet)
	if err != nil {
		return nil, err
	}
	if okRet.ErrorCode != 0 {
		return nil, codeError(okRet.ErrorCode)
	}
	var ret proto.OrderPlaceReturn

	ret.OrderID = strconv.FormatInt(okRet.OrderID, 10)

	return &ret, nil
}

func (o *OKEX) OrderCancel(params *proto.OrderCancelParams) error {
	path := "cancel_order.do"
	values := url.Values{}
	values.Set("order_id", params.OrderID)
	values.Set("symbol", strings.ToLower(params.BaseCurrencyID)+"_"+strings.ToLower(params.QuoteCurrencyID))

	var okRet CancelOrderReturn
	err := o.apiKeyPost(values, path, &okRet)
	if err != nil {
		return err
	}
	if okRet.ErrorCode != 0 {
		return codeError(okRet.ErrorCode)
	}
	if !okRet.Result {
		return errors.New("撤单失败")
	}

	return nil

}

func NewOKEX(apiKey, secretKey string) (*OKEX, error) {
	return &OKEX{
		accessKey: apiKey,
		secretKey: secretKey,
	}, nil
}