package amm

import uuid "github.com/satori/go.uuid"

type AMM struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	ChainID        uint      `json:"chainID"`
	RouterAddress  string    `json:"routerAddress"`
	FactoryAddress string    `json:"factoryAddress"`
}

var SupportedAMMs = map[int][]AMM{
	43114: { // avalanche mainnet
		{
			Name:           "Trader Joe",
			URL:            "https://traderjoexyz.com/",
			ChainID:        43114, // avalanche mainnet
			FactoryAddress: "0x9Ad6C38BE94206cA50bb0d90783181662f0Cfa10",
			RouterAddress:  "0x60aE616a2155Ee3d9A68541Ba4544862310933d4",
		},
		{
			Name:           "Sushi Swap",
			URL:            "https://app.sushi.com/",
			ChainID:        43114, // avalanche mainnet
			FactoryAddress: "0xc35DADB65012eC5796536bD9864eD8773aBc74C4",
			RouterAddress:  "0x1b02dA8Cb0d097eB8D57A175b88c7D8b47997506",
		},
		{
			Name:           "Pangolin",
			URL:            "https://app.pangolin.exchange/",
			ChainID:        43114, // avalanche mainnet
			FactoryAddress: "0xefa94DE7a4656D787667C749f7E1223D71E9FD88",
			RouterAddress:  "0xE54Ca86531e17Ef3616d22Ca28b0D458b6C89106",
		},
	},
	250: { // fantom mainnet
		{
			Name:           "Spooky Swap",
			URL:            "https://spookyswap.finance/",
			ChainID:        250, // fantom mainnet
			FactoryAddress: "0x152eE697f2E276fA89E96742e9bB9aB1F2E61bE3",
			RouterAddress:  "0xF491e7B69E4244ad4002BC14e878a34207E38c29",
		},
		{
			Name:           "Spirit Swap",
			URL:            "https://app.spiritswap.finance/",
			ChainID:        250, // fantom mainnet
			FactoryAddress: "0xef45d134b73241eda7703fa787148d9c9f4950b0",
			RouterAddress:  "0x16327e3fbdaca3bcf7e38f5af2599d2ddc33ae52",
		},
	},
	1313161554: { // aurora mainnet
		{
			Name:           "Trisolaris",
			URL:            "https://www.trisolaris.io/",
			ChainID:        1313161554, // aurora mainnet
			FactoryAddress: "0xc66F594268041dB60507F00703b152492fb176E7",
			RouterAddress:  "0x2CB45Edb4517d5947aFdE3BEAbF95A582506858B",
		},
		{
			Name:           "Wanna Swap",
			URL:            "https://wannaswap.finance/",
			ChainID:        1313161554, // aurora mainnet
			FactoryAddress: "0x7928D4FeA7b2c90C732c10aFF59cf403f0C38246",
			RouterAddress:  "0xa3a1eF5Ae6561572023363862e238aFA84C72ef5",
		},
		{
			Name:           "Rose Finance",
			URL:            "https://app.rose.fi/",
			ChainID:        1313161554, // aurora mainnet
			FactoryAddress: "",         // TODO
			RouterAddress:  "",         // TODO
		},
		{
			Name:           "Near Pad",
			URL:            "https://dex.nearpad.io/",
			ChainID:        1313161554, // aurora mainnet
			FactoryAddress: "",         // TODO
			RouterAddress:  "",         // TODO
		},
	},
}
