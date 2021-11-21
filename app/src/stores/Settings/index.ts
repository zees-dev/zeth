import { writable, derived } from 'svelte/store'
import type { SettingsResponse } from '../../types'
import { settingsURL, rpcNodeURL } from '../../lib/const'
import { ethers } from 'ethers'

export const settingsStore = (() => {
	const { subscribe, set } = writable<SettingsResponse>()

	const getSettings = async () => {
		// TODO error handling
		const raw = await fetch(settingsURL)
		const response = (await raw.json()) as SettingsResponse
		set(response)
	}

	// eslint-disable-next-line @typescript-eslint/no-floating-promises
	getSettings()

	return { subscribe }
})()

//
export const rpcURL = (() => {
	const { subscribe, set } = writable<string>()
	return { subscribe, set }
})()


export const ethersProvider = (() => {
	const { subscribe, update } = writable<ethers.providers.WebSocketProvider>(new ethers.providers.WebSocketProvider(''))

	// initially set to defaultProvider
	const unsubscribeFromDefaultRPC = settingsStore.subscribe((settings) => {
		const defaultNodeID = settings?.nodeSettings.defaultNodeID
		if (defaultNodeID) {
			const rpcURL = `ws://${window.location.host}${rpcNodeURL}/${defaultNodeID}`
			setRPCURL(rpcURL)
		}
	})

	function setRPCURL(rpcURL: string) {
		// unsubscribe on first call to setRPCURL
		unsubscribeFromDefaultRPC()
		update((provider) => {
			provider.removeAllListeners()
			// eslint-disable-next-line @typescript-eslint/no-floating-promises
			provider.destroy()
			return new ethers.providers.WebSocketProvider(rpcURL)
		})
	}

	return { subscribe, setRPCURL }
})()
