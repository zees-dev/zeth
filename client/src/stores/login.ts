import { writable } from 'svelte/store';
import { replace } from 'svelte-spa-router';
import { web3ProviderStore } from './web3provider';
import { dbStore } from './db';

const USER_TOKEN_KEY = 'token';

export const loginStore = (() => {
  const token = localStorage.getItem(USER_TOKEN_KEY);
  const { subscribe, set } = writable({
    loggedIn: !!token,
    token: token,
  });

	return {
		subscribe,
		// login: () => update(state => ({ })),
		login: (token: string) => {
      localStorage.setItem(USER_TOKEN_KEY, token);
      set({
        loggedIn: true,
        token,
      });
    },
    logout: () => {
      localStorage.removeItem(USER_TOKEN_KEY);
      set({
        loggedIn: false,
        token: '',
      });
      web3ProviderStore.disconnect();
      dbStore.disconnect();
      replace('/'); // redirect to home
    },
	};
})();
