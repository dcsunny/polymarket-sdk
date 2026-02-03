// order_builder.go 模块
package polymarket

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	order_utils_builder "github.com/polymarket/go-order-utils/pkg/builder"
	order_utils_model "github.com/polymarket/go-order-utils/pkg/model"
)

// OrderBuilder 封装 go-order-utils ExchangeOrderBuilder。
type OrderBuilder struct {
	builder    order_utils_builder.ExchangeOrderBuilder
	privateKey *ecdsa.PrivateKey
	address    common.Address
	funder     common.Address
	sigType    int
}

// NewOrderBuilder 创建新的订单构建器。
func NewOrderBuilder(client *CLOBClient, privateKey, address string, sigType int, funder string) *OrderBuilder {
	pk, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		panic(fmt.Sprintf("failed to parse private key: %v", err))
	}

	var funderAddr common.Address
	if funder == "" {
		funderAddr = common.HexToAddress(address)
	} else {
		funderAddr = common.HexToAddress(funder)
	}

	return &OrderBuilder{
		builder:    order_utils_builder.NewExchangeOrderBuilderImpl(big.NewInt(client.ChainID()), nil),
		privateKey: pk,
		address:    common.HexToAddress(address),
		funder:     funderAddr,
		sigType:    sigType,
	}
}

// BuildAndSignOrder 构建并签名限价订单。
func (ob *OrderBuilder) BuildAndSignOrder(args *OrderArgs, negRisk bool) (*order_utils_model.SignedOrder, error) {
	var side order_utils_model.Side
	if args.Side == SideBuy {
		side = order_utils_model.BUY
	} else {
		side = order_utils_model.SELL
	}

	nonceStr := args.Nonce
	if nonceStr == "" {
		nonceStr = "0"
	}
	if args.Taker == "" {
		args.Taker = "0x0000000000000000000000000000000000000000"
	}
	if args.FeeRateBps == "" {
		args.FeeRateBps = "0"
	}
	if args.Expiration == "" {
		args.Expiration = "0"
	}

	orderData := &order_utils_model.OrderData{
		Maker:         ob.funder.Hex(),
		Signer:        ob.address.Hex(),
		Taker:         args.Taker,
		TokenId:       args.TokenID,
		MakerAmount:   args.MakerAmount,
		TakerAmount:   args.TakerAmount,
		Side:          side,
		FeeRateBps:    args.FeeRateBps,
		Nonce:         nonceStr,
		Expiration:    args.Expiration,
		SignatureType: order_utils_model.SignatureType(ob.sigType),
	}

	contract := order_utils_model.CTFExchange
	if negRisk {
		contract = order_utils_model.NegRiskCTFExchange
	}

	return ob.builder.BuildSignedOrder(ob.privateKey, orderData, contract)
}
