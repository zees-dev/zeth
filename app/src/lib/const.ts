// URLS
export const settingsURL = '/api/v1/settings'

export const nodesURL = '/api/v1/nodes'
export const rpcNodeURL = `${nodesURL}/rpc`

export const httpNodeRPCURL = (id: string) => `http://${window.location.host}${rpcNodeURL}/${id}`
export const wsNodeRPCURL = (id: string) => `ws://${window.location.host}${rpcNodeURL}/${id}`

// Supported Node types
export enum NodeType {
  TypeGethNodeInProcess = 1,
  TypeGethNode,
  TypeRemoteNode,
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
  1: 'Mainnet',
  3: 'Ropsten',
  4: 'Rinkeby',
  5: 'Goerli',
  2018: 'Dev',
}

// BSC testnets
// source: https://docs.binance.org/smart-chain/wallet/metamask.html
