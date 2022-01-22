package bots

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zees-dev/zeth/pkg/solidity/ethutil"
	periphery "github.com/zees-dev/zeth/pkg/solidity/uniswap/v2-periphery"
)

// AVAX_TOKENS.MAG, AVAX_TOKENS.MIM, magBalanceUser, MAG_MIM_PAIR_ADDRESS, TRADERJOE_ROUTER
func getPrice(providerURL string, token1Addr, token2Addr string, input int, pairAddr, routerAddr string) (float64, error) {
	client, err := ethclient.Dial(providerURL)
	if err != nil {
		log.Fatal(err)
	}

	// erc20Token1, err := standard.NewERC20(common.HexToAddress(token1Addr), client)
	// if err != nil {
	// 	return 0, err
	// }

	// traderjoe factory
	// factory, err := core.NewUniswapV2Factory(common.HexToAddress("0x9Ad6C38BE94206cA50bb0d90783181662f0Cfa10"), client)
	// if err != nil {
	// 	return 0, err
	// }

	// pairAddress, err := factory.GetPair(&bind.CallOpts{}, common.HexToAddress(token1Addr), common.HexToAddress(token2Addr))
	// if err != nil {
	// 	return 0, err
	// }

	// pair, err := core.NewUniswapV2Pair(pairAddress, client)
	// if err != nil {
	// 	return 0, err
	// }

	// reserves, err := pair.GetReserves(&bind.CallOpts{})
	// if err != nil {
	// 	return 0, err
	// }

	router, err := periphery.NewUniswapV2Router02(common.HexToAddress(routerAddr), client)
	if err != nil {
		return 0, err
	}

	outTokens, err := router.GetAmountsOut(
		&bind.CallOpts{},
		ethutil.ToWei(100, 18),
		[]common.Address{common.HexToAddress(token1Addr), common.HexToAddress(token2Addr)},
	)
	if err != nil {
		return 0, err
	}

	fmt.Println(ethutil.ToDecimal(outTokens[0], 18), ethutil.ToDecimal(outTokens[1], 9))

	// inTokens, err := router.GetAmountsIn(
	// 	&bind.CallOpts{},
	// 	ethutil.ToWei(10, 9),
	// 	[]common.Address{common.HexToAddress(token1Addr), common.HexToAddress(token2Addr)},
	// )
	// if err != nil {
	// 	return 0, err
	// }

	// fmt.Println(ethutil.ToDecimal(inTokens[0], 18), ethutil.ToDecimal(inTokens[1], 9))

	return 0, nil
}
