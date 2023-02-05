import { writable } from 'svelte/store';
import { replace } from 'svelte-spa-router';
import { web3ProviderStore } from './web3provider';
import { dbStore } from './db';
import { parseJWT } from '../lib/utils';

const USER_TOKEN_KEY = 'token';

export const loginStore = (() => {
  const token = localStorage.getItem(USER_TOKEN_KEY);
  const userId = !!token ? parseJWT(token!)?.ID : '';
  const { subscribe, set } = writable({
    loggedIn: !!token,
    token: token,
    userId,
  });

	return {
		subscribe,
		// login: () => update(state => ({ })),
		login: (token: string) => {
      localStorage.setItem(USER_TOKEN_KEY, token);
      set({
        loggedIn: true,
        token,
        userId: parseJWT(token!)?.ID,
      });
    },
    logout: () => {
      localStorage.removeItem(USER_TOKEN_KEY);
      set({
        loggedIn: false,
        token: '',
        userId: '',
      });
      web3ProviderStore.disconnect();
      dbStore.disconnect();
      replace('/'); // redirect to home
    },
	};
})();
