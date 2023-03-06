
export type ConnectionStatus = 'connecting' | 'connected' | undefined

export interface Endpoint {
  enabled: boolean;
  name: string;
  rpc_url: string;
  symbol: string;
  block_explorer_url?: string,

  // added by db
  id: string;
  user: string;
  date_added: string;
  proxy_url: string,
}
