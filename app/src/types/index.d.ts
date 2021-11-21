/**
 * Can be made globally available by placing this
 * inside `global.d.ts` and removing `export` keyword
 */

// API response interfaces
export interface SettingsResponse {
  nodeSettings: {
    supportedNodes: {
      nodeType: number
      version: string
    }[]
    defaultNodeID: string
  }
}

export interface SyncStatus {
  currentBlock: string
  highestBlock: string
  knownStates: string
  pulledStates: string
  startingBlock: string
}


export interface Node {
  id: string
  dateAdded: string
  name: string
  nodeType: number
  isDev: boolean
  enabled: boolean
  rpc: {
    http: string
    ws: string
  }
}

export interface NodesResponse {
  nodes: Node[]
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
