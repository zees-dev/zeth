package bots

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getPrice(t *testing.T) {
	is := assert.New(t)

	// func getPrice(providerURL string, token1Addr, token2Addr string, input int, pairAddr, routerAddr string)

	providerURL := "https://api.avax.network/ext/bc/C/rpc"
	token1Addr := "0x130966628846bfd36ff31a822705796e8cb8c18d" // AVAX - MIM
	token2Addr := "0x1d60109178C48E4A937D8AB71699D8eBb6F7c5dE" // AVAX - MAG
	input := 1
	// factoryAddr := "0x9Ad6C38BE94206cA50bb0d90783181662f0Cfa10"
	pairAddr := "0x147B8eb97fD247D06C4006D269c90C1908Fb5D54"
	routerAddr := "0x60aE616a2155Ee3d9A68541Ba4544862310933d4"

	p, err := getPrice(providerURL, token1Addr, token2Addr, input, pairAddr, routerAddr)
	is.NoError(err)
	is.Equal(p, float64(0))
}
