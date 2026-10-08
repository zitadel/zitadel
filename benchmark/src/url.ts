// @ts-ignore Import module
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

import { Config } from './config';

export type options = {
  searchParams?: { [name: string]: string };
};

export default function url(path: string, options: options = {}) {
  let url = new URL(Config.host + path);

  if (options.searchParams) {
    Object.entries(options.searchParams).forEach(([key, value]) => {
      url.searchParams.append(key, value);
    });
  }

  return url.toString();
}
