import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'redirect',
})
export class RedirectPipe implements PipeTransform {
    public transform(uri: string, isNative: boolean): boolean {
        console.log(uri, isNative);
        if (isNative) {
            if (uri.startsWith('http://localhost')) {
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
