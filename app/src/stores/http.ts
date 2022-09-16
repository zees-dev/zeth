import { writable, Writable } from 'svelte/store'

// TODO: use simple getters and setters to GET/PUT data from/to the store
// https://youtu.be/1Df-9EKvZr0?t=6261

export const fetchStates = {
  IDLE: 'idle',
  LOADING: 'loading',
  ERROR: 'error',
  SUCCESS: 'success'
}

interface APIResponse<T> {
  status: keyof typeof fetchStates | string
  data: {
    error?: unknown | Error
    response?: T
  }
}

/**
 * fetchData is a store that holds the state of the fetching of data
 * @source https://svelte.dev/repl/a63591efba17499da9165341dcc9fb13?version=3.37.0
 * @param url
 * @returns
 */
export function fetchData<T>(url: string): [Writable<APIResponse<T>>, () => Promise<void>] {
  const store = writable({ status: fetchStates.IDLE, data: {} })
  async function get() {
    store.set({ status: fetchStates.LOADING, data: {} })
    try {
      const response = await fetch(url)
      store.set({
        status: fetchStates.SUCCESS,
        data: {
          response: (await response.json()) as T
        }
      })
    } catch (e) {
      store.set({ status: fetchStates.ERROR, data: { error: e } })
    }
  }
  // eslint-disable-next-line @typescript-eslint/no-floating-promises
  get()
  return [store, get]
}



/**
 * subscribedFetch returns a store with HTTP access functions for get, post, patch, delete
 * anytime an HTTP request is made, the store is updated and all subscribers are notified.
 * @source https://gist.github.com/joshnuss/e4c4a4965f12b1d6012393a4ccdb7462
 * @param initial
 * @returns
 */
// export function subscribedFetch(initial) {
//   // create the underlying store
//   const store = writable(initial)

//   // define a request function that will do `fetch` and update store when request finishes
//   store.request = async (method, url, params=null) => {
//     // before we fetch, clear out previous errors and set `loading` to `true`
//     store.update(data => {
//       delete data.errors
//       data.loading = true

//       return data
//     })

//     // define headers and body
//     const headers = {
//       "Content-type": "application/json"
//     }
//     const body = params ? JSON.stringify(params) : undefined

//     // execute fetch
//     const response = await fetch(url, { method, body, headers })
//     // pull out json body
//     const json = await response.json()

//     // if response is 2xx
//     if (response.ok) {
//       // update the store, which will cause subscribers to be notified
//       store.set(json)
//     } else {
//       // response failed, set `errors` and clear `loading` flag
//       store.update(data => {
//         data.loading = false
//         data.errors = json.errors
//         return data
//       })
//     }
//   }

//   // convenience wrappers for get, post, patch, and delete
//   store.get = (url) => store.request('GET', url)
//   store.post = (url, params) => store.request('POST', url, params)
//   store.patch = (url, params) => store.request('PATCH', url, params)
//   store.delete = (url, params) => store.request('DELETE', url, params)

//   // return the customized store
//   return store
// }