import { wrap } from "svelte-spa-router/wrap";

export default {
  '/': wrap({ asyncComponent: () => import('../components/endpoints-list/EndpointsListLayout.svelte') }),
  '/endpoints': wrap({ asyncComponent: () => import('../components/endpoints-list/EndpointsListLayout.svelte') }),
  '/endpoints/:endpointId': wrap({ asyncComponent: () => import('../components/endpoint/EndpointLayout.svelte') }),
  '/endpoints/:endpointId/*': wrap({ asyncComponent: () => import('../components/endpoint/EndpointLayout.svelte') }), // TODO <- there may be a better way to do this
  '/accounts': wrap({ asyncComponent: () => import('../components/Accounts.svelte') }),
  '/contracts': wrap({ asyncComponent: () => import('../components/Contracts.svelte') }),
  '/settings': wrap({ asyncComponent: () => import('../components/Settings.svelte') }),
  // Catch-all route last
  '*': wrap({ asyncComponent: () => import('../components/NotFound.svelte') }),

  // example of fallback loading component (when fetching routed component)
  // '/hello/:first/:last?': wrap({
  //   // Note that this is a function that returns the import
  //   asyncComponent: () => import('./routes/Name.svelte'),
  //   // Show the loading component while the component is being downloaded
  //   loadingComponent: Loading,
  //   // Pass values for the `params` prop of the loading component
  //   loadingParams: {
  //       message: 'Loading the Name route…'
  //   }
  // }),
}
