import { writable } from 'svelte/store';
import { ethers } from "ethers";
import Web3Modal from "web3modal";

const METAMASK_CONNECTION_KEY = 'metamask';

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
    const instance = await web3Modal.connect();
    const provider = new ethers.providers.Web3Provider(instance);
    // const signer = provider.getSigner();

    const accounts = await provider.listAccounts();
    const network = await provider.getNetwork();
    console.log(accounts);
    console.log(network.chainId);

    localStorage.setItem(METAMASK_CONNECTION_KEY, 'true');

    set({ provider });
  }

  if (isConnected) {
    connect();
  };

	return {
		subscribe,
		// login: () => update(state => ({ })),
		connect,
    disconnect: async () => {
      web3Modal.clearCachedProvider();
      localStorage.removeItem(METAMASK_CONNECTION_KEY);
      set(initialState);
    },
	};
})();
