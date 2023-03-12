import type Surreal from "surrealdb.js";
import { ENDPOINT_TABLE, type Endpoint } from "./types";

interface QueryParams {
  throwErr: boolean;
}

interface SurrealResponse<T> {
  status: string;
  time: string;
  result: T[];
}

interface QueryResponse<T> {
  loading: boolean;
  result: T;
  error: unknown | boolean;
}

export async function getEndpoint(db: Surreal, endpointId: string, params?: QueryParams) {
  let loading = true;
  let result: Endpoint; // always set if error not thrown
  let error: unknown | boolean = false;
  try {
    const response = await db.query<SurrealResponse<Endpoint>[]>(`SELECT * FROM ${ENDPOINT_TABLE} WHERE id=${endpointId};`);
    if (response[0].status !== "OK") {
      throw new Error(`failed to get endpoint; ${endpointId}`);
    }
    if (response[0].result.length === 0) {
      throw new Error(`endpoint not found; ${endpointId}`);
    }
    result = response[0].result[0];
  } catch (e) {
    error = e;
    if (params?.throwErr) throw error;
  } finally {
    loading = false;
  }
  return { loading, result: result!, error };
}

export async function getEndpoints(db: Surreal, params?: QueryParams) {
  let loading = true;
  let result: Endpoint[] = [];
  let error: unknown | boolean = false;
  try {
    const response = await db.query<SurrealResponse<Endpoint>[]>(
      "SELECT * FROM type::table($tb);",
      { tb: ENDPOINT_TABLE },
    );
    if (response[0].status !== "OK") {
      throw new Error(`failed to get endpoints`);
    }
    result = response[0].result;
  } catch (e) {
    error = e;
    if (params?.throwErr) throw error;
  } finally {
    loading = false;
  }
  return { loading, result, error };
}
