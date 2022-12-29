export interface API {
  token: string;
  systemToken: string,
  mgmtBaseURL: string;
  adminBaseURL: string;
  systemBaseURL: string;
}

export type SearchResult = {
  entity: Entity | null;
  sequence: number;
  id: number;
};

// Entity is an object but not a function
export type Entity = { [k: string]: any } & ({ bind?: never } | { call?: never });
