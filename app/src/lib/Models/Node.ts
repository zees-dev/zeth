import { ethers } from 'ethers'
import type { SyncStatus, NodeResponse, RPCModules } from '../../types'
import { EthNetworks, httpNodeRPCURL, wsNodeRPCURL } from '../../lib/const'

export class Node {
  id: string
  dateAdded: string
  name: string
  isDev: boolean
  enabled: boolean
  rpc: {
    http: string
    ws: string
  }
  explorerUrl: string

  httpProvider?: ethers.providers.JsonRpcProvider
  wsProvider?: ethers.providers.WebSocketProvider

  connected: boolean
  network?: ethers.providers.Network
  block: number
  version: string
  syncing: SyncStatus | boolean
  peers: number
  modules: RPCModules
  mining: boolean
  isDefault: boolean

  constructor(nodeResponse: NodeResponse) {
    // fields set from API request
    const { id, dateAdded, name, isDev, enabled, rpc, explorerUrl } = nodeResponse
    this.id = id
    this.dateAdded = dateAdded
    this.name = name
    this.isDev = isDev
    this.enabled = enabled
    this.rpc = rpc
    this.explorerUrl = explorerUrl

    this.connected = false
    this.block = 0
    this.version = ''
    this.syncing = false
    this.peers = 0
    this.modules = {}
    this.mining = false
    this.isDefault = false

    if (this.rpc.http) {
      this.httpProvider = new ethers.providers.JsonRpcProvider(httpNodeRPCURL(id))
    }
    if (this.rpc.ws) {
      this.wsProvider = new ethers.providers.WebSocketProvider(wsNodeRPCURL(id))
    }
  }

  setHTTPProvider(url: string) {
    this.httpProvider = new ethers.providers.JsonRpcProvider(url)
  }

  setWSProvider(url: string) {
    this.wsProvider = new ethers.providers.WebSocketProvider(url)
  }

  async getRPCData() {
    if (!this.httpProvider) {
      throw new Error("JsonRpcProvider required")
    }

    const [networkP, blockP, versionP, syncingP, peersP, modulesP, miningP] = await Promise.allSettled([
      this.httpProvider.getNetwork(),
      this.httpProvider.getBlockNumber(),
      this.httpProvider.send('web3_clientVersion', []),
      this.httpProvider.send('eth_syncing', []),
      this.httpProvider.send('net_peerCount', []),
      this.httpProvider.send('rpc_modules', []),
      this.httpProvider.send('eth_mining', []),
    ])

    if (networkP.status === 'fulfilled') {
      this.network = networkP.value
      this.connected = true
    }
    if (blockP.status === 'fulfilled') {
      this.block = blockP.value
    }
    if (versionP.status === 'fulfilled') {
      this.version = versionP.value as string
    }
    if (syncingP.status === 'fulfilled') {
      this.syncing = syncingP.value as SyncStatus
    }
    if (peersP.status === 'fulfilled') {
      this.peers = parseInt(peersP.value as string, 16)
    }
    if (modulesP.status === 'fulfilled') {
      this.modules = modulesP.value as RPCModules
    }
    if (miningP.status === 'fulfilled') {
      this.mining = miningP.value as boolean
    }

    return this
  }

  async getCoinbase() {
    if (!this.httpProvider) {
      throw new Error("JsonRpcProvider required")
    }

    try {
      const coinbase = (await this.httpProvider.send('eth_coinbase', []) as string)
      return coinbase
    } catch (err) {
      console.error(err)
      // TODO: notify
    }
  }

}


export function getNetworkName(chainId: number) {
  // TODO: account for other non-eth chains
  return EthNetworks[chainId]
}

export function getNodeSyncStatus(node: Node) {
  if (!node.connected) {
    return undefined
  }
  if (node.syncing) {
    return 'syncing'
  }
  if (node.syncing === undefined) {
    return undefined
  }
  return 'synced'
}

/**
 * getVersion returns the geth version from the version string.
 * It gets the string between first '/' and second '/' characters.
 * @param gethVersion example: "Geth/v1.10.9-omnibus-e03773e6/linux-amd64/go1.17.2"
 * @returns example: "v1.10.9-omnibus-e03773e6"
 */
export function getVersion(gethVersion: string) {
  return gethVersion.split('/')[1] ?? gethVersion
}

export function dateWithoutTZ(date: Date) {
  return date.toString().substring(0, date.toString().indexOf('GMT') - 1)
}

export function getSortedModules(modules: RPCModules) {
  return Object.keys(modules)
    .sort()
    .map((key: string) => ({ module: key, version: modules[key] }))
}