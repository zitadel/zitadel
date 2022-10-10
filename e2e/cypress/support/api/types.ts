export interface API {
  authHeader: string;
  mgntBaseURL: string;
  adminBaseURL: string;
}

export type SearchResult = {
  entity: Entity | null;
  sequence: number;
  id: number;
};

// Entity is an object but not a function
export type Entity = { [k: string]: any } & ({ bind?: never } | { call?: never });
