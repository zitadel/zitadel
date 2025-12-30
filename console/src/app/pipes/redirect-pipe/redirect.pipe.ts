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
    if (!isNative) {
      return parsedURI.protocol === 'https:';
    }
    if (parsedURI.protocol !== 'http:' && parsedURI.protocol !== 'https:') {
      return true;
    }
    const hostname = parsedURI.hostname;
    return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '[::1]';
  }
}
