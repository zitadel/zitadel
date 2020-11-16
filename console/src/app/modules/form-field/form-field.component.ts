import {
    AfterContentInit,
    ChangeDetectorRef,
    Component,
    ContentChild,
    ContentChildren,
    HostListener,
    InjectionToken,
    QueryList,
} from '@angular/core';
import { Subject } from 'rxjs';
import { startWith } from 'rxjs/operators';

import { CNSL_ERROR, CnslErrorDirective } from './error.directive';
import { CnslFormFieldControlDirective } from './form-field-control.directive';

export const CNSL_FORM_FIELD = new InjectionToken<FormFieldComponent>('CnslFormField');

@Component({
    selector: 'cnsl-form-field',
    templateUrl: './form-field.component.html',
    styleUrls: ['./form-field.component.scss'],
    providers: [
        { provide: CNSL_FORM_FIELD, useExisting: FormFieldComponent },
    ],
})
export class FormFieldComponent implements AfterContentInit {
    focused: boolean = false;
    readonly stateChanges: Subject<void> = new Subject<void>();

    @ContentChild(CnslFormFieldControlDirective) _controlNonStatic!: CnslFormFieldControlDirective<any>;
    @ContentChild(CnslFormFieldControlDirective, { static: true }) _controlStatic!: CnslFormFieldControlDirective<any>;
    get _control(): CnslFormFieldControlDirective<any> {
        // TODO(crisbeto): we need this workaround in order to support both Ivy and ViewEngine.
        //  We should clean this up once Ivy is the default renderer.
        return this._explicitFormFieldControl || this._controlNonStatic || this._controlStatic;
    }
    set _control(value: CnslFormFieldControlDirective<any>) {
        this._explicitFormFieldControl = value;
    }
    private _explicitFormFieldControl!: CnslFormFieldControlDirective<any>;


    @ContentChildren(CNSL_ERROR as any, { descendants: true }) _errorChildren!: QueryList<CnslErrorDirective>;

    @HostListener('blur', ['false'])
    _focusChanged(isFocused: boolean): void {
        console.log('blur1');
        if (isFocused !== this.focused && (!isFocused)) {
            this.focused = isFocused;
            this.stateChanges.next();
        }
    }

    constructor(private _changeDetectorRef: ChangeDetectorRef) { }

    public ngAfterContentInit(): void {
        // Update the aria-described by when the number of errors changes.
        this._errorChildren.changes.pipe(startWith(null)).subscribe(() => {
            this._syncDescribedByIds();
            this._changeDetectorRef.markForCheck();
        });
    }

    private _syncDescribedByIds(): void {
        if (this._control) {
            const ids: string[] = [];

            // TODO(wagnermaciel): Remove the type check when we find the root cause of this bug.
            if (this._control.userAriaDescribedBy &&
                typeof this._control.userAriaDescribedBy === 'string') {
                ids.push(...this._control.userAriaDescribedBy.split(' '));
            }

            if (this._errorChildren) {
                ids.push(...this._errorChildren.map(error => error.id));
            }

            this._control.setDescribedByIds(ids);
        }
    }
}
