import { Config } from './config.js';

export default function url(path, options = {}) {
  let url = new URL(Config.host+path);

   if (options.searchParams !== undefined) {
    Object.entries(options.searchParams).forEach(([key, value]) => {
      url.searchParams.append(key, value);
    })
  }
  
  return url.toString();
}