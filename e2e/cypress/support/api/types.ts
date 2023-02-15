export interface Token {
  token: string;
}

export interface API extends Token {
  mgmtBaseURL: string;
  adminBaseURL: string;
  authBaseURL: string;
  assetsBaseURL: string;
  oidcBaseURL: string;
  oauthBaseURL: string;
  samlBaseURL: string;
}

export interface SystemAPI extends Token {
  baseURL: string;
}

export type SearchResult = {
  entity: Entity | null;
  sequence: number;
  id: string;
};

// Entity is an object but not a function
export type Entity = { [k: string]: any } & ({ bind?: never } | { call?: never });
