import { coerceBooleanProperty } from '@angular/cdk/coercion';
import { Platform } from '@angular/cdk/platform';
import { AutofillMonitor } from '@angular/cdk/text-field';
import {
    AfterViewInit,
    Directive,
    DoCheck,
    ElementRef,
    HostListener,
    Inject,
    Input,
    OnChanges,
    OnDestroy,
    Optional,
    Self,
} from '@angular/core';
import { FormGroupDirective, NgControl, NgForm } from '@angular/forms';
import { CanUpdateErrorStateCtor, mixinErrorState } from '@angular/material/core';
import { MatFormFieldControl } from '@angular/material/form-field';
import { MAT_INPUT_VALUE_ACCESSOR } from '@angular/material/input';
import { Subject } from 'rxjs';

import { CNSL_FORM_FIELD, CnslFormFieldComponent } from '../form-field/form-field.component';
import { ErrorStateMatcher } from './error-options';

let nextUniqueId = 0;

class CnslInputBase {
    constructor(public _defaultErrorStateMatcher: ErrorStateMatcher,
        public _parentForm: NgForm,
        public _parentFormGroup: FormGroupDirective,
        /** @docs-private */
        public ngControl: NgControl) { }
}

const _CnslInputMixinBase: CanUpdateErrorStateCtor & typeof CnslInputBase =
    mixinErrorState(CnslInputBase);

@Directive({
    selector: '[cnslInput]',
    host: {
        /**
         * @breaking-change 8.0.0 remove .mat-form-field-autofill-control in favor of AutofillMonitor.
         */
        // 'class': 'mat-input-element mat-form-field-autofill-control',
        // '[class.mat-input-server]': '_isServer',
        // Native input properties that are overwritten by Angular inputs need to be synced with
        // the native input element. Otherwise property bindings for those don't work.
        '[attr.data-placeholder]': 'placeholder',
        '[attr.id]': 'id',
        '[disabled]': 'disabled',
        '[required]': 'required',
        '[attr.aria-invalid]': 'errorState',
        '[attr.aria-required]': 'required.toString()',
    },
    providers: [{ provide: MatFormFieldControl, useExisting: InputDirective }],
})
export class InputDirective extends _CnslInputMixinBase implements MatFormFieldControl<any>,
    OnChanges, AfterViewInit, OnDestroy, DoCheck {
    /** Whether the element is readonly. */
    // @Input()
    // get readonly(): boolean { return this._readonly; }
    // set readonly(value: boolean) { this._readonly = coerceBooleanProperty(value); }
    // private _readonly = false;
    // readonly id!: string;

    /* tslint:disable */
    static ngAcceptInputType_value: any;
    /* tslint:enable */

    protected _previousNativeValue: any;
    private _inputValueAccessor: { value: any; };
    autofilled: boolean = false;
    shouldLabelFloat!: boolean;
    // controlType?: string | undefined;
    userAriaDescribedBy?: string | undefined;

    protected _uid: string = `cnsl-input-${nextUniqueId++}`;

    protected _id!: string;
    protected _required: boolean = false;
    protected _type: string = 'text';

    protected _disabled: boolean = false;

    focused: boolean = false;
    readonly stateChanges: Subject<void> = new Subject<void>();

    controlType: string = 'cnslInput';

    /**
    * Implemented as part of MatFormFieldControl.
    * @docs-private
    */
    @Input() placeholder!: string;

    @Input() errorStateMatcher!: ErrorStateMatcher;

    @HostListener('blur', ['false'])
    _focusChanged(isFocused: boolean): void {
        console.log('blur');
        if (isFocused !== this.focused && (!isFocused)) {
            this.focused = isFocused;
            this.stateChanges.next();
        }
    }

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    @Input()
    get disabled(): boolean {
        if (this.ngControl && this.ngControl.disabled !== null) {
            return this.ngControl.disabled;
        }
        return this._disabled;
    }
    set disabled(value: boolean) {
        this._disabled = coerceBooleanProperty(value);

        // Browsers may not fire the blur event if the input is disabled too quickly.
        // Reset from here to ensure that the element doesn't become stuck.
        if (this.focused) {
            this.focused = false;
            this.stateChanges.next();
        }
    }

    @HostListener('input') _onInput(): any {
        // This is a noop function and is used to let Angular know whenever the value changes.
        // Angular will run a new change detection each time the `input` event has been dispatched.
        // It's necessary that Angular recognizes the value change, because when floatingLabel
        // is set to false and Angular forms aren't used, the placeholder won't recognize the
        // value changes and will not disappear.
        // Listening to the input event wouldn't be necessary when the input is using the
        // FormsModule or ReactiveFormsModule, because Angular forms also listens to input events.
    }

    /**
   * Implemented as part of MatFormFieldControl.
   * @docs-private
   */
    @Input()
    get id(): string { return this._id; }
    set id(value: string) { this._id = value || this._uid; }

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    @Input()
    get required(): boolean { return this._required; }
    set required(value: boolean) { this._required = coerceBooleanProperty(value); }


    @Input()
    get type(): string { return this._type; }
    set type(value: string) {
        this._type = value || 'text';

        // When using Angular inputs, developers are no longer able to set the properties on the native
        // input element. To ensure that bindings for `type` work, we need to sync the setter
        // with the native property. Textarea elements don't support the type property or attribute.
        (this._elementRef.nativeElement as HTMLInputElement).type = this._type;
    }

    @Input()
    get value(): string { return this._inputValueAccessor.value; }
    set value(value: string) {
        if (value !== this.value) {
            this._inputValueAccessor.value = value;
            this.stateChanges.next();
        }
    }

    constructor(
        protected _elementRef: ElementRef<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>,
        protected _platform: Platform,
        @Optional() @Self() public ngControl: NgControl,
        @Optional() _parentForm: NgForm,
        @Optional() _parentFormGroup: FormGroupDirective,
        _defaultErrorStateMatcher: ErrorStateMatcher,
        @Optional() @Self() @Inject(MAT_INPUT_VALUE_ACCESSOR) inputValueAccessor: any,
        private _autofillMonitor: AutofillMonitor,
        @Optional() @Inject(CNSL_FORM_FIELD) private _formField?: CnslFormFieldComponent,
    ) {
        super(_defaultErrorStateMatcher, _parentForm, _parentFormGroup, ngControl);

        const element = this._elementRef.nativeElement;
        this._inputValueAccessor = inputValueAccessor || element;
        this._previousNativeValue = this.value;
        this.id = this.id;
    }

    public ngAfterViewInit(): void {
        if (this._platform.isBrowser) {
            this._autofillMonitor.monitor(this._elementRef.nativeElement).subscribe(event => {
                this.autofilled = event.isAutofilled;
                this.stateChanges.next();
            });
        }
    }

    public ngOnChanges(): void {
        this.stateChanges.next();
    }

    public ngOnDestroy(): void {
        this.stateChanges.complete();

        if (this._platform.isBrowser) {
            this._autofillMonitor.stopMonitoring(this._elementRef.nativeElement);
        }
    }

    public ngDoCheck(): void {
        if (this.ngControl) {
            // We need to re-evaluate this on every change detection cycle, because there are some
            // error triggers that we can't subscribe to (e.g. parent form submissions). This means
            // that whatever logic is in here has to be super lean or we risk destroying the performance.
            this.updateErrorState();
        }

        // We need to dirty-check the native element's value, because there are some cases where
        // we won't be notified when it changes (e.g. the consumer isn't using forms or they're
        // updating the value using `emitEvent: false`).
        this._dirtyCheckNativeValue();

        // We need to dirty-check and set the placeholder attribute ourselves, because whether it's
        // present or not depends on a query which is prone to "changed after checked" errors.
    }

    focus(options?: FocusOptions): void {
        this._elementRef.nativeElement.focus(options);
    }

    /** Does some manual dirty checking on the native input `value` property. */
    protected _dirtyCheckNativeValue(): void {
        const newValue = this._elementRef.nativeElement.value;

        if (this._previousNativeValue !== newValue) {
            this._previousNativeValue = newValue;
            this.stateChanges.next();
        }
    }

    /** Checks whether the input is invalid based on the native validation. */
    protected _isBadInput(): boolean {
        // The `validity` property won't be present on platform-server.
        const validity = (this._elementRef.nativeElement as HTMLInputElement).validity;
        return validity && validity.badInput;
    }

    get empty(): boolean {
        return !this._elementRef.nativeElement.value && !this._isBadInput() &&
            !this.autofilled;
    }

    public setDescribedByIds(ids: string[]): void {
        if (ids.length) {
            this._elementRef.nativeElement.setAttribute('aria-describedby', ids.join(' '));
        } else {
            this._elementRef.nativeElement.removeAttribute('aria-describedby');
        }
    }

    public onContainerClick(): void {
        // Do not re-focus the input element if the element is already focused. Otherwise it can happen
        // that someone clicks on a time input and the cursor resets to the "hours" field while the
        // "minutes" field was actually clicked. See: https://github.com/angular/components/issues/12849
        if (!this.focused) {
            this.focus();
        }
    }

}
