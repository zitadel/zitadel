import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'redirect',
  standalone: false,
})
export class RedirectPipe implements PipeTransform {
  public transform(uri: string, isNative: boolean): boolean {
    let parsedURI = URL.parse(uri);
    if (parsedURI === null) {
      return false;
    }
    if (isNative) {
      if (parsedURI.protocol === 'http:' || parsedURI.protocol === 'https:') {
        let hostname = parsedURI.hostname;
        if (
          hostname === 'localhost' ||
          hostname === '127.0.0.1' ||
          hostname === '[::1]'
        ) {
          return true;
        } else {
          return false;
        }
      } else {
        return true;
      }
    } else {
      return parsedURI.protocol === 'https:'
    }
  }
}
