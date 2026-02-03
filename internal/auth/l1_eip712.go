// l1_eip712.go 模块
package auth

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	clobAuthDomain  = "ClobAuthDomain"
	clobAuthVersion = "1"
	clobAuthMessage = "This message attests that I control the given wallet"
)

// ClobAuthSignature 签名 CLOB 认证类型数据并返回十六进制签名。
func ClobAuthSignature(privateKeyHex, address, timestamp string, nonce int, chainID int64) (string, error) {
	if privateKeyHex == "" || address == "" {
		return "", errors.New("missing private key or address")
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return "", err
	}

	hash, err := clobAuthTypedDataHash(address, timestamp, nonce, chainID)
	if err != nil {
		return "", err
	}

	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", err
	}
	if sig[64] < 27 {
		sig[64] += 27
	}
	return "0x" + hex.EncodeToString(sig), nil
}

func clobAuthTypedDataHash(address, timestamp string, nonce int, chainID int64) ([]byte, error) {
	typeHash := crypto.Keccak256([]byte("ClobAuth(address address,string timestamp,uint256 nonce,string message)"))

	addr := common.HexToAddress(address)
	addressBytes := common.LeftPadBytes(addr.Bytes(), 32)

	timestampHash := crypto.Keccak256([]byte(timestamp))
	nonceBytes := common.LeftPadBytes(big.NewInt(int64(nonce)).Bytes(), 32)
	messageHash := crypto.Keccak256([]byte(clobAuthMessage))

	data := make([]byte, 0, 32*5)
	data = append(data, typeHash...)
	data = append(data, addressBytes...)
	data = append(data, timestampHash...)
	data = append(data, nonceBytes...)
	data = append(data, messageHash...)
	structHash := crypto.Keccak256(data)

	domainTypeHash := crypto.Keccak256([]byte("EIP712Domain(string name,string version,uint256 chainId)"))
	nameHash := crypto.Keccak256([]byte(clobAuthDomain))
	versionHash := crypto.Keccak256([]byte(clobAuthVersion))
	chainIDBytes := common.LeftPadBytes(big.NewInt(chainID).Bytes(), 32)

	domainData := make([]byte, 0, 32*4)
	domainData = append(domainData, domainTypeHash...)
	domainData = append(domainData, nameHash...)
	domainData = append(domainData, versionHash...)
	domainData = append(domainData, chainIDBytes...)
	domainSeparator := crypto.Keccak256(domainData)

	finalData := make([]byte, 0, 2+32+32)
	finalData = append(finalData, 0x19, 0x01)
	finalData = append(finalData, domainSeparator...)
	finalData = append(finalData, structHash...)

	return crypto.Keccak256(finalData), nil
}
