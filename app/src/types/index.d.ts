/**
 * Can be made globally available by placing this
 * inside `global.d.ts` and removing `export` keyword
 */

// API response interfaces
export interface SettingsResponse {
  nodeSettings: {
    supportedNodes: {
      version: string
    }[]
  }
}

export interface SyncStatus {
  currentBlock: string
  highestBlock: string
  knownStates: string
  pulledStates: string
  startingBlock: string
}

export interface RPCModules {
  [key: string | undefined]: string
  admin?: string
  clique?: string
  debug?: string
  eth?: string
  miner?: string
  net?: string
  personal?: string
  rpc?: string
  txpool?: string
  web3?: string
}

export interface NodeResponse {
  id: string
  dateAdded: string
  name: string
  isDev: boolean
  enabled: boolean
  explorerUrl: string
  rpc: {
    http: string
    ws: string
  }
}

export interface NodesListResponse {
  nodes: NodeResponse[]
}

export type NetworkName = 'ETHEREUM' | 'BINANCE_SMART_CHAIN' | 'POLYGON'

export interface Chain {
  key: NetworkName
  name: string
  apiURL: URL
  explorerURL: URL
  gasURL: URL
}

export interface Locals {
  userid: string
}

export interface TokenData {
  id: string
  symbol: string
  name: string
  address: string
  logoURI: string
}

export interface TokenResponse {
  tokens: {
    [key: string]: TokenData
  }
}

export interface AMM {
  id: string
  name: string
  chainID: number
  url: string
  routerAddress: string
  factoryAddress: string
} 
