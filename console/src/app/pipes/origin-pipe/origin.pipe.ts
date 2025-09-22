import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'origin',
  standalone: false,
})
export class OriginPipe implements PipeTransform {
  public transform(value: string): boolean {
    return new RegExp(/^((https?:\/\/).*?([\w\d-]*\.[\w\d]+))($|\/.*$)/gm).test(value);
  }
}
