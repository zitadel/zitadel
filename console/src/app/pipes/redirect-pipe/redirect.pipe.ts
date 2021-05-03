import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'redirect',
})
export class RedirectPipe implements PipeTransform {
    public transform(uri: string, isNative: boolean): boolean {
        if (isNative) {
            if (uri.startsWith('http://localhost/') || uri.startsWith('http://localhost:') || uri.startsWith('http://127.0.0.1') || uri.startsWith('http://[::1]') || uri.startsWith('http://[0:0:0:0:0:0:0:1]')) {
                return true;
            }
            if (!uri.startsWith('https://') && !uri.startsWith('http://')) {
                return true;
            } else {
                return false;
            }
        } else {
            if (uri.startsWith('https://')) {
                return true;
            } else {
                return false;
            }
        }
    }
}
