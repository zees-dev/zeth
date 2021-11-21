import { writable } from 'svelte/store'
import { Node } from '../../lib/Models/Node'

export const nodeStore = (() => {
  const { subscribe, set } = writable<Node>()
  return { subscribe, set }
})()
