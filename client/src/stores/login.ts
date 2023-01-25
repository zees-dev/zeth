import { writable } from 'svelte/store';
import { replace } from 'svelte-spa-router';

export const loginStore = (() => {
  const token = localStorage.getItem('token');
  const { subscribe, set } = writable({
    loggedIn: !!token,
    token: token,
  });

	return {
		subscribe,
		// login: () => update(state => ({ })),
		login: (token: string) => {
      localStorage.setItem('token', token);
      set({
        loggedIn: true,
        token,
      });
    },
    logout: () => {
      localStorage.removeItem('token');
      set({
        loggedIn: false,
        token: '',
      });
      replace('/'); // redirect to home
    },
	};
})();
