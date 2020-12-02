import {
    AfterContentInit,
    AfterViewInit,
    ChangeDetectionStrategy,
    ChangeDetectorRef,
    Component,
    ContentChild,
    ContentChildren,
    ElementRef,
    HostListener,
    Inject,
    InjectionToken,
    OnDestroy,
    QueryList,
    ViewChild,
    ViewEncapsulation,
} from '@angular/core';
import { NgControl } from '@angular/forms';
import { MatFormFieldControl } from '@angular/material/form-field';
import { Subject } from 'rxjs';
import { startWith, takeUntil } from 'rxjs/operators';

import { cnslFormFieldAnimations } from './animations';
import { CNSL_ERROR, CnslErrorDirective } from './error.directive';
import { _CNSL_HINT, CnslHintDirective } from './hint.directive';

export const CNSL_FORM_FIELD = new InjectionToken<CnslFormFieldComponent>('CnslFormFieldComponent');

class CnslFormFieldBase {
    constructor(public _elementRef: ElementRef) { }
}

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
        '[class.cnsl-form-field-disabled]': '_control.disabled',
        '[class.cnsl-form-field-autofilled]': '_control.autofilled',
        '[class.cnsl-focused]': '_control.focused',
        '[class.cnsl-form-field-invalid]': '_control.errorState',
    },
    encapsulation: ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush,
    animations: [cnslFormFieldAnimations.transitionMessages],
})
export class CnslFormFieldComponent extends CnslFormFieldBase implements OnDestroy, AfterContentInit, AfterViewInit {
    focused: boolean = false;
    private _destroyed: Subject<void> = new Subject<void>();

    @ViewChild('connectionContainer', { static: true }) _connectionContainerRef!: ElementRef;
    @ViewChild('inputContainer') _inputContainerRef!: ElementRef;
    @ContentChild(MatFormFieldControl) _controlNonStatic!: MatFormFieldControl<any>;
    @ContentChild(MatFormFieldControl, { static: true }) _controlStatic!: MatFormFieldControl<any>;
    get _control(): MatFormFieldControl<any> {
        // TODO(crisbeto): we need this workaround in order to support both Ivy and ViewEngine.
        //  We should clean this up once Ivy is the default renderer.
        return this._explicitFormFieldControl || this._controlNonStatic || this._controlStatic;
    }
    set _control(value: MatFormFieldControl<any>) {
        this._explicitFormFieldControl = value;
    }
    private _explicitFormFieldControl!: MatFormFieldControl<any>;
    readonly stateChanges: Subject<void> = new Subject<void>();

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

    constructor(public _elementRef: ElementRef, private _changeDetectorRef: ChangeDetectorRef,
        @Inject(ElementRef)
        // Use `ElementRef` here so Angular has something to inject.
        _labelOptions: any) {
        super(_elementRef);

    }

    public ngAfterViewInit(): void {
        // Avoid animations on load.
        this._subscriptAnimationState = 'enter';
        this._changeDetectorRef.detectChanges();
    }

    public ngOnDestroy(): void {
        this._destroyed.next();
        this._destroyed.complete();
    }

    public ngAfterContentInit(): void {
        this._validateControlChild();

        const control = this._control;
        control.stateChanges.pipe(startWith(null)).subscribe(() => {
            this._syncDescribedByIds();
            this._changeDetectorRef.markForCheck();
        });

        // Run change detection if the value changes.
        if (control.ngControl && control.ngControl.valueChanges) {
            control.ngControl.valueChanges
                .pipe(takeUntil(this._destroyed))
                .subscribe(() => this._changeDetectorRef.markForCheck());
        }

        // Update the aria-described by when the number of errors changes.
        this._errorChildren.changes.pipe(startWith(null)).subscribe(() => {
            this._syncDescribedByIds();
            this._changeDetectorRef.markForCheck();
        });
    }

    /** Throws an error if the form field's control is missing. */
    protected _validateControlChild(): void {
        if (!this._control) {
            throw Error('cnsl-form-field must contain a MatFormFieldControl.');
        }
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
