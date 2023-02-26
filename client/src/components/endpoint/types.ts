
export type ConnectionStatus = 'connecting' | 'connected' | undefined

export interface Endpoint {
  date_added: string;
  enabled: boolean;
  id: string;
  name: string;
  rpc_url: string;
  type: 'http' | 'ws';
  user: string;
  symbol: string;
  block_explorer_url?: string,
}
