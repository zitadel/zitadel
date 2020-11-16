import { Component, ContentChildren, HostListener, QueryList } from '@angular/core';
import { Subject } from 'rxjs';

import { CNSL_ERROR, CnslErrorDirective } from './error.directive';

@Component({
    selector: 'cnsl-form-field',
    templateUrl: './form-field.component.html',
    styleUrls: ['./form-field.component.scss'],
})
export class FormFieldComponent {
    focused: boolean = false;
    readonly stateChanges: Subject<void> = new Subject<void>();

    @ContentChildren(CNSL_ERROR as any, { descendants: true }) _errorChildren!: QueryList<CnslErrorDirective>;

    @HostListener('blur', ['false'])
    _focusChanged(isFocused: boolean): void {
        console.log('blur1');
        if (isFocused !== this.focused && (!isFocused)) {
            this.focused = isFocused;
            this.stateChanges.next();
        }
    }

    constructor() { }
}
