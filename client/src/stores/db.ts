import { writable } from 'svelte/store';
import Surreal from 'surrealdb.js';

interface Store { db: Surreal };

// docker run --rm -it --name surrealdb -p 8000:8000 surrealdb/surrealdb:latest start --log trace --user admin --pass admin memory
const SURREAL_DB_URL = 'http://127.0.0.1:8000/rpc';

export const dbStore = (() => {
  const { subscribe, update } = writable({ db: new Surreal() });

	return {
		subscribe,
    connect: async () => update(state => {
      state.db.connect(SURREAL_DB_URL);
      return state;
    }),
    disconnect: async () => update(state => {
      state.db.invalidate();
      return state;
    }),
	};
})();

