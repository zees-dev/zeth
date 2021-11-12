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
