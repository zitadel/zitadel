import { Directive, Input } from '@angular/core';

let nextUniqueId = 0;

@Directive({
    selector: '[cnslError]',
    host: {
        'class': 'cnsl-error',
        'role': 'alert',
        '[attr.id]': 'id',
    },
})
export class ErrorDirective {
    @Input() id: string = `cnsl-error-${nextUniqueId++}`;

    constructor() { }
}
