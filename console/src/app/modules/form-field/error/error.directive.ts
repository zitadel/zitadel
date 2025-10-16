import { Directive, InjectionToken, Input } from '@angular/core';

let nextUniqueId = 0;

export const CNSL_ERROR = new InjectionToken<CnslErrorDirective>('CnslError');

@Directive({
  selector: '[cnslError]',
  host: {
    class: 'cnsl-error',
    role: 'alert',
    '[attr.id]': 'id',
  },
  providers: [{ provide: CNSL_ERROR, useExisting: CnslErrorDirective }],
  standalone: false,
})
export class CnslErrorDirective {
  @Input() id: string = `cnsl-error-${nextUniqueId++}`;

  constructor() {}
}
