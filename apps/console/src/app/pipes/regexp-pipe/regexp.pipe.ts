import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'regexp',
})
export class RegexpPipe implements PipeTransform {
  public transform(value: string): RegExp {
    return new RegExp(value);
  }
}
