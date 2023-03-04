import { writable } from 'svelte/store';
import { ethers } from "ethers";
import Web3Modal from "web3modal";

/**
 * WEB3_CONNECT_CACHED_PROVIDER is the key used to store the provider connection status in localStorage
 * It is used by web3modal to determine if a provider is already connected.
 * note: metamask will set this storage item automatically to `"injected"` upon connection
 */
const METAMASK_CONNECTION_KEY = 'WEB3_CONNECT_CACHED_PROVIDER';

export const web3ProviderStore = (() => {
  interface Store { provider?: ethers.providers.Web3Provider };
  const initialState: Store = { provider: undefined }; 
  const { subscribe, set } = writable(initialState);

  const isConnected = localStorage.getItem(METAMASK_CONNECTION_KEY);
  const providerOptions = {
    /* See Provider Options Section */
  };
  const web3Modal = new Web3Modal({
    network: "mainnet", // optional
    cacheProvider: true, // optional
    providerOptions // required
  });

  async function connect() {
    let instance;
    try {
      instance = await web3Modal.connect();
    } catch (e) {
      console.warn("could not get a wallet connection; disconnecting...", e);
      disconnect();
      return;
    }

    const provider = new ethers.providers.Web3Provider(instance);
    // const signer = provider.getSigner();

    const accounts = await provider.listAccounts();
    const network = await provider.getNetwork();
    console.log(accounts);
    console.log(network.chainId);

    set({ provider });
  }

  async function disconnect() {
    web3Modal.clearCachedProvider();
    localStorage.removeItem(METAMASK_CONNECTION_KEY);
    set(initialState);
  }

  if (isConnected) {
    connect();
  };

	return {
		subscribe,
		// login: () => update(state => ({ })),
		connect,
    disconnect,
	};
})();
