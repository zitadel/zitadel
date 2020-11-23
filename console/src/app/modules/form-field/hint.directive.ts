import { Directive, InjectionToken, Input } from '@angular/core';

let nextUniqueId = 0;

/**
 * Injection token that can be used to reference instances of `MatHint`. It serves as
 * alternative token to the actual `MatHint` class which could cause unnecessary
 * retention of the class and its directive metadata.
 *
 * *Note*: This is not part of the public API as the MDC-based form-field will not
 * need a lightweight token for `MatHint` and we want to reduce breaking changes.
 */
export const _CNSL_HINT = new InjectionToken<CnslHintDirective>('CnslHintDirective');

/** Hint text to be shown underneath the form field control. */
@Directive({
    selector: 'cnsl-hint',
    host: {
        'class': 'cnsl-hint',
        '[class.cnsl-form-field-hint-end]': 'align === "end"',
        '[attr.id]': 'id',
        // Remove align attribute to prevent it from interfering with layout.
        '[attr.align]': 'null',
    },
    providers: [{ provide: _CNSL_HINT, useExisting: CnslHintDirective }],
})
export class CnslHintDirective {
    /** Whether to align the hint label at the start or end of the line. */
    @Input() align: 'start' | 'end' = 'start';

    /** Unique ID for the hint. Used for the aria-describedby on the form field control. */
    @Input() id: string = `mat-hint-${nextUniqueId++}`;
}
