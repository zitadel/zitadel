import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'roletransform',
})
export class RoleTransformPipe implements PipeTransform {
  public transform(value: string | string[]): string {
    if (typeof value === 'string') {
      return getNewString(value);
    } else if (typeof value === 'object' && value.length) {
      return (
        value
          .map((s) => getNewString(s))
          // .slice(0, -1)
          .join(', ')
      );
    } else {
      return '';
    }
  }
}

function getNewString(value: string): string {
  const splitted = value.toLowerCase().split('_');
  const uppercased = splitted.map((s) => `${s.substring(0, 1).toUpperCase()}${s.substring(1, s.length)}`);
  return uppercased.join(' ');
}
