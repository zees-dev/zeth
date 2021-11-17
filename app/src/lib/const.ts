// URLS
export const settingsURL = '/api/v1/settings'

export const nodesURL = '/api/v1/nodes'
export const rpcNodeURL = `${nodesURL}/rpc`

export const httpNodeRPCURL = (id: string) => `http://${window.location.host}${rpcNodeURL}/${id}`
export const wsNodeRPCURL = (id: string) => `ws://${window.location.host}${rpcNodeURL}/${id}`

// Supported Node types
export enum NodeType {
  GethNodeInProcess = 1,
  GethNode,
  RemoteNode,
}

export enum NetworkType {
  Mainnet = 1,
  BinanceSmartChain = 56,
  Polygon = 137,
}


// ETH testnets
// source: https://besu.hyperledger.org/en/stable/Concepts/NetworkID-And-ChainID/
// mainnet	ETH	1	1	Production
// ropsten	ETH	3	3	Test
// rinkeby	ETH	4	4	Test
// goerli	ETH	5	5	Test
// dev	ETH	2018	2018	Development
export const EthNetworks: { [id: number]: string } = {
  // Eth networks
  1: 'Mainnet',
  3: 'Ropsten',
  4: 'Rinkeby',
  5: 'Goerli',
  2018: 'Dev',

  // Binance Smart Chain
  // source: https://docs.binance.org/smart-chain/wallet/metamask.html
  56: 'Binance Smart Chain',
  97: 'BSC - testnet',

  // Polygon
  // source: https://docs.polygon.io/docs/polygon-node-setup
  137: 'Polygon',

  // Avalanche
  // source: https://docs.avalanche.io/docs/avalanche-node-setup
  43114: 'Avalanche',

  // Fantom
  250: 'Fantom',

  // Harmony
  1666600000: 'Harmony',
}


