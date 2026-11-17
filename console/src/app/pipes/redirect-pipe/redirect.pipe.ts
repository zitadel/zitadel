import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'redirect',
  standalone: false,
})
export class RedirectPipe implements PipeTransform {
  public transform(uri: string, isNative: boolean): boolean {
    const parsedURI = URL.parse(uri);
    if (parsedURI === null) {
      return false;
    }
    // https is valid for all application types
    if (parsedURI.protocol === 'https:') {
      return true;
    }
    if (!isNative) {
      return false;
    }
    // native applications can also use custom schemes
    if (parsedURI.protocol !== 'http:') {
      return true;
    }
    // native http redirects are restricted to loopback addresses
    const hostname = parsedURI.hostname;
    return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '[::1]';
  }
}
