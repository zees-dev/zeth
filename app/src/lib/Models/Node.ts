import type { SyncStatus, NodeResponse, RPCModules } from '../../types'
import { ethers } from 'ethers'

export class Node {
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

  httpProvider?: ethers.providers.JsonRpcProvider
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
    const { id, dateAdded, name, nodeType, isDev, enabled, rpc } = nodeResponse
    this.id = id
    this.dateAdded = dateAdded
    this.name = name
    this.nodeType = nodeType
    this.isDev = isDev
    this.enabled = enabled
    this.rpc = rpc

    this.connected = false
    this.block = 0
    this.version = ''
    this.syncing = false
    this.peers = 0
    this.modules = {}
    this.mining = false
    this.isDefault = false
  }

  setHTTPProvider(httpProvider: ethers.providers.JsonRpcProvider) {
    this.httpProvider = httpProvider
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
      this.httpProvider.send('eth_mining', [])
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

}
