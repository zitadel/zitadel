import {
    AfterContentInit,
    AfterViewInit,
    ChangeDetectorRef,
    Component,
    ContentChild,
    ContentChildren,
    ElementRef,
    HostListener,
    InjectionToken,
    QueryList,
    ViewChild,
} from '@angular/core';
import { NgControl } from '@angular/forms';
import { Subject } from 'rxjs';
import { startWith } from 'rxjs/operators';

import { cnslFormFieldAnimations } from './animations';
import { CNSL_ERROR, CnslErrorDirective } from './error.directive';
import { CnslFormFieldControlDirective } from './form-field-control.directive';
import { _CNSL_HINT, CnslHintDirective } from './hint.directive';

export const CNSL_FORM_FIELD = new InjectionToken<CnslFormFieldComponent>('CnslFormFieldComponent');


@Component({
    selector: 'cnsl-form-field',
    templateUrl: './form-field.component.html',
    styleUrls: ['./form-field.component.scss'],
    providers: [
        { provide: CNSL_FORM_FIELD, useExisting: CnslFormFieldComponent },
    ],
    host: {
        '[class.ng-untouched]': '_shouldForward("untouched")',
        '[class.ng-touched]': '_shouldForward("touched")',
        '[class.ng-pristine]': '_shouldForward("pristine")',
        '[class.ng-dirty]': '_shouldForward("dirty")',
        '[class.ng-valid]': '_shouldForward("valid")',
        '[class.ng-invalid]': '_shouldForward("invalid")',
        '[class.ng-pending]': '_shouldForward("pending")',
        // '[class.cnsl-form-field-invalid]': '_control.errorState',
    },
    animations: [cnslFormFieldAnimations.transitionMessages],
})
export class CnslFormFieldComponent implements AfterContentInit, AfterViewInit {
    focused: boolean = false;
    readonly stateChanges: Subject<void> = new Subject<void>();

    @ViewChild('inputContainer') _inputContainerRef!: ElementRef;
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

    _subscriptAnimationState: string = '';

    @ContentChildren(CNSL_ERROR as any, { descendants: true }) _errorChildren!: QueryList<CnslErrorDirective>;
    @ContentChildren(_CNSL_HINT, { descendants: true }) _hintChildren!: QueryList<CnslHintDirective>;

    @HostListener('blur', ['false'])
    _focusChanged(isFocused: boolean): void {
        console.log('blur1');
        if (isFocused !== this.focused && (!isFocused)) {
            this.focused = isFocused;
            this.stateChanges.next();
        }
    }

    constructor(public _elementRef: ElementRef, private _changeDetectorRef: ChangeDetectorRef) {
    }

    public ngAfterViewInit(): void {
        // Avoid animations on load.
        this._subscriptAnimationState = 'enter';
        this._changeDetectorRef.detectChanges();
    }

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

    /** Determines whether a class from the NgControl should be forwarded to the host element. */
    _shouldForward(prop: keyof NgControl): boolean {
        const ngControl = this._control ? this._control.ngControl : null;
        return ngControl && ngControl[prop];
    }

    /** Determines whether to display hints or errors. */
    _getDisplayedMessages(): 'error' | 'hint' {
        return (this._errorChildren && this._errorChildren.length > 0) ? 'error' : 'hint';
    }
}
