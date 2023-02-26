
/**
 * Decode JWT token
 * @param token JWT token
 * @returns decoded JSW obj
 * source: https://stackoverflow.com/a/38552302/10813908
 */
export function parseJWT(token: string) {
  const base64Url = token.split('.')[1];
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));

  return JSON.parse(jsonPayload);
}

export function endpointType(rpcUrl: string) {
  return rpcUrl.startsWith('http') ? 'http' : 'ws';
}

/**
 * @param isoString Date ISO string like so: 2023-02-25T12:53:05.053Z
 * @returns date in format dd/mm/yyyy HH:MM -> example '25/02/2023 12:53'
 */
export function formatDate(isoString: string): string {
  const date = new Date(isoString);
  const day = date.getDate().toString().padStart(2, '0');
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const year = date.getFullYear().toString();
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  return `${day}/${month}/${year} ${hours}:${minutes}`;
}

export async function testRPCConnection(rpcUrl: string) {
  const type = endpointType(rpcUrl);
  if (type === 'http') {
    const res = await fetch(rpcUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "jsonrpc": "2.0",
        "method": "eth_blockNumber",
        "params": [],
        "id": 1
      }),
    });
    if (!res.ok) {
      throw new Error(`Failed to connect to ${rpcUrl}`);
    }
  } else {
    const socket = new WebSocket(rpcUrl);
    socket.onopen = () => {
      socket.send(JSON.stringify({
        "jsonrpc": "2.0",
        "method": "eth_blockNumber",
        "params": [],
        "id": 1
      }));
    }
    try {
      await new Promise((resolve, reject) => {
        socket.onmessage = (msg) => {
          const res = JSON.parse(msg.data);
          if (res.error) {
            reject(res.error);
            return;
          }
          resolve(res);
        }
        socket.onerror = (err) => reject(err);
      })
    } catch (err) {
      throw new Error(`Failed to connect to ${rpcUrl}`);
    } finally {
      socket!.close();
    }
  }
}