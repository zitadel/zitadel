import { BooleanInput, coerceBooleanProperty } from '@angular/cdk/coercion';
import { getSupportedInputTypes, Platform } from '@angular/cdk/platform';
import { AutofillMonitor } from '@angular/cdk/text-field';
import {
    AfterViewInit,
    Directive,
    DoCheck,
    ElementRef,
    HostListener,
    Inject,
    Input,
    NgZone,
    OnChanges,
    OnDestroy,
    Optional,
    Self,
} from '@angular/core';
import { FormGroupDirective, NgControl, NgForm } from '@angular/forms';
import { CanUpdateErrorState, CanUpdateErrorStateCtor, ErrorStateMatcher, mixinErrorState } from '@angular/material/core';
import { MAT_FORM_FIELD, MatFormField, MatFormFieldControl } from '@angular/material/form-field';
import { getMatInputUnsupportedTypeError, MAT_INPUT_VALUE_ACCESSOR } from '@angular/material/input';
import { Subject } from 'rxjs';


// Invalid input type. Using one of these will throw an MatInputUnsupportedTypeError.
const MAT_INPUT_INVALID_TYPES = [
    'button',
    'checkbox',
    'file',
    'hidden',
    'image',
    'radio',
    'range',
    'reset',
    'submit',
];

let nextUniqueId = 0;

// Boilerplate for applying mixins to MatInput.
/** @docs-private */
class MatInputBase {
    constructor(public _defaultErrorStateMatcher: ErrorStateMatcher,
        public _parentForm: NgForm,
        public _parentFormGroup: FormGroupDirective,
        /** @docs-private */
        public ngControl: NgControl) { }
}
const _MatInputMixinBase: CanUpdateErrorStateCtor & typeof MatInputBase =
    mixinErrorState(MatInputBase);

/** Directive that allows a native input to work inside a `MatFormField`. */
@Directive({
    selector: `input[cnslInput], textarea[cnslInput], select[cnslNativeControl]`,
    exportAs: 'cnslInput',
    host: {
        /**
         * @breaking-change 8.0.0 remove .mat-form-field-autofill-control in favor of AutofillMonitor.
         */
        // 'class': 'cnsl-input-element cnsl-form-field-autofill-control',
        // '[class.mat-input-server]': '_isServer',
        // Native input properties that are overwritten by Angular inputs need to be synced with
        // the native input element. Otherwise property bindings for those don't work.
        '[attr.id]': 'id',
        // At the time of writing, we have a lot of customer tests that look up the input based on its
        // placeholder. Since we sometimes omit the placeholder attribute from the DOM to prevent screen
        // readers from reading it twice, we have to keep it somewhere in the DOM for the lookup.
        '[attr.data-placeholder]': 'placeholder',
        '[disabled]': 'disabled',
        '[required]': 'required',
        '[attr.readonly]': 'readonly && !_isNativeSelect || null',
        '[attr.aria-invalid]': 'errorState',
        '[attr.aria-required]': 'required.toString()',
    },
    providers: [{ provide: MatFormFieldControl, useExisting: InputDirective }],
})
export class InputDirective extends _MatInputMixinBase implements MatFormFieldControl<any>, OnChanges,
    OnDestroy, AfterViewInit, DoCheck, CanUpdateErrorState {
    protected _uid: string = `cnsl-input-${nextUniqueId++}`;
    protected _previousNativeValue: any;
    private _inputValueAccessor: { value: any; };
    private _previousPlaceholder!: string | null;

    /** Whether the component is being rendered on the server. */
    readonly _isServer: boolean;

    /** Whether the component is a native html select. */
    readonly _isNativeSelect: boolean;

    /** Whether the component is a textarea. */
    readonly _isTextarea: boolean;

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    focused: boolean = false;

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    readonly stateChanges: Subject<void> = new Subject<void>();

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    controlType: string = 'mat-input';

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    autofilled: boolean = false;

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
    protected _disabled: boolean = false;

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    @Input()
    get id(): string { return this._id; }
    set id(value: string) { this._id = value || this._uid; }
    protected _id!: string;

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    @Input() placeholder!: string;

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    @Input()
    get required(): boolean { return this._required; }
    set required(value: boolean) { this._required = coerceBooleanProperty(value); }
    protected _required: boolean = false;

    /** Input type of the element. */
    @Input()
    get type(): string { return this._type; }
    set type(value: string) {
        this._type = value || 'text';
        this._validateType();

        // When using Angular inputs, developers are no longer able to set the properties on the native
        // input element. To ensure that bindings for `type` work, we need to sync the setter
        // with the native property. Textarea elements don't support the type property or attribute.
        if (!this._isTextarea && getSupportedInputTypes().has(this._type)) {
            (this._elementRef.nativeElement as HTMLInputElement).type = this._type;
        }
    }
    protected _type: string = 'text';

    /** An object used to control when error messages are shown. */
    @Input() errorStateMatcher!: ErrorStateMatcher;

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    // tslint:disable-next-line:no-input-rename
    @Input('aria-describedby') userAriaDescribedBy!: string;

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    @Input()
    get value(): string { return this._inputValueAccessor.value; }
    set value(value: string) {
        if (value !== this.value) {
            this._inputValueAccessor.value = value;
            this.stateChanges.next();
        }
    }

    /** Whether the element is readonly. */
    @Input()
    get readonly(): boolean { return this._readonly; }
    set readonly(value: boolean) { this._readonly = coerceBooleanProperty(value); }
    private _readonly: boolean = false;

    protected _neverEmptyInputTypes: string[] = [
        'date',
        'datetime',
        'datetime-local',
        'month',
        'time',
        'week',
    ].filter(t => getSupportedInputTypes().has(t));

    constructor(
        protected _elementRef: ElementRef<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>,
        protected _platform: Platform,
        /** @docs-private */
        @Optional() @Self() public ngControl: NgControl,
        @Optional() _parentForm: NgForm,
        @Optional() _parentFormGroup: FormGroupDirective,
        _defaultErrorStateMatcher: ErrorStateMatcher,
        @Optional() @Self() @Inject(MAT_INPUT_VALUE_ACCESSOR) inputValueAccessor: any,
        private _autofillMonitor: AutofillMonitor,
        ngZone: NgZone,
        @Optional() @Inject(MAT_FORM_FIELD) private _formField?: MatFormField) {

        super(_defaultErrorStateMatcher, _parentForm, _parentFormGroup, ngControl);

        const element = this._elementRef.nativeElement;
        const nodeName = element.nodeName.toLowerCase();

        // If no input value accessor was explicitly specified, use the element as the input value
        // accessor.
        this._inputValueAccessor = inputValueAccessor || element;

        this._previousNativeValue = this.value;

        // Force setter to be called in case id was not specified.
        this.id = this.id;

        // On some versions of iOS the caret gets stuck in the wrong place when holding down the delete
        // key. In order to get around this we need to "jiggle" the caret loose. Since this bug only
        // exists on iOS, we only bother to install the listener on iOS.
        if (_platform.IOS) {
            ngZone.runOutsideAngular(() => {
                _elementRef.nativeElement.addEventListener('keyup', (event: Event) => {
                    const el: HTMLInputElement = event.target as HTMLInputElement;
                    if (!el.value && !el.selectionStart && !el.selectionEnd) {
                        // Note: Just setting `0, 0` doesn't fix the issue. Setting
                        // `1, 1` fixes it for the first time that you type text and
                        // then hold delete. Toggling to `1, 1` and then back to
                        // `0, 0` seems to completely fix it.
                        el.setSelectionRange(1, 1);
                        el.setSelectionRange(0, 0);
                    }
                });
            });
        }

        this._isServer = !this._platform.isBrowser;
        this._isNativeSelect = nodeName === 'select';
        this._isTextarea = nodeName === 'textarea';

        if (this._isNativeSelect) {
            this.controlType = (element as HTMLSelectElement).multiple ? 'mat-native-select-multiple' :
                'mat-native-select';
        }
    }

    ngAfterViewInit(): void {
        if (this._platform.isBrowser) {
            this._autofillMonitor.monitor(this._elementRef.nativeElement).subscribe(event => {
                this.autofilled = event.isAutofilled;
                this.stateChanges.next();
            });
        }
    }

    ngOnChanges(): void {
        this.stateChanges.next();
    }

    ngOnDestroy(): void {
        this.stateChanges.complete();

        if (this._platform.isBrowser) {
            this._autofillMonitor.stopMonitoring(this._elementRef.nativeElement);
        }
    }

    ngDoCheck(): void {
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
        this._dirtyCheckPlaceholder();
    }

    /** Focuses the input. */
    focus(options?: FocusOptions): void {
        this._elementRef.nativeElement.focus(options);
    }

    // We have to use a `HostListener` here in order to support both Ivy and ViewEngine.
    // In Ivy the `host` bindings will be merged when this class is extended, whereas in
    // ViewEngine they're overwritten.
    /** Callback for the cases where the focused state of the input changes. */
    // tslint:disable:no-host-decorator-in-concrete
    @HostListener('focus', ['true'])
    @HostListener('blur', ['false'])
    // tslint:enable:no-host-decorator-in-concrete
    _focusChanged(isFocused: boolean): void {
        if (isFocused !== this.focused && (!this.readonly || !isFocused)) {
            this.focused = isFocused;
            this.stateChanges.next();
        }
    }

    // We have to use a `HostListener` here in order to support both Ivy and ViewEngine.
    // In Ivy the `host` bindings will be merged when this class is extended, whereas in
    // ViewEngine they're overwritten.
    // tslint:disable-next-line:no-host-decorator-in-concrete
    @HostListener('input')
    _onInput(): void {
        // This is a noop function and is used to let Angular know whenever the value changes.
        // Angular will run a new change detection each time the `input` event has been dispatched.
        // It's necessary that Angular recognizes the value change, because when floatingLabel
        // is set to false and Angular forms aren't used, the placeholder won't recognize the
        // value changes and will not disappear.
        // Listening to the input event wouldn't be necessary when the input is using the
        // FormsModule or ReactiveFormsModule, because Angular forms also listens to input events.
    }

    /** Does some manual dirty checking on the native input `placeholder` attribute. */
    private _dirtyCheckPlaceholder(): void {
        // If we're hiding the native placeholder, it should also be cleared from the DOM, otherwise
        // screen readers will read it out twice: once from the label and once from the attribute.
        const placeholder = this._formField?._hideControlPlaceholder?.() ? null : this.placeholder;
        if (placeholder !== this._previousPlaceholder) {
            const element = this._elementRef.nativeElement;
            this._previousPlaceholder = placeholder;
            placeholder ?
                element.setAttribute('placeholder', placeholder) : element.removeAttribute('placeholder');
        }
    }

    /** Does some manual dirty checking on the native input `value` property. */
    protected _dirtyCheckNativeValue(): void {
        const newValue = this._elementRef.nativeElement.value;

        if (this._previousNativeValue !== newValue) {
            this._previousNativeValue = newValue;
            this.stateChanges.next();
        }
    }

    /** Make sure the input is a supported type. */
    protected _validateType(): void {
        if (MAT_INPUT_INVALID_TYPES.indexOf(this._type) > -1) {
            throw getMatInputUnsupportedTypeError(this._type);
        }
    }

    /** Checks whether the input type is one of the types that are never empty. */
    protected _isNeverEmpty(): boolean {
        return this._neverEmptyInputTypes.indexOf(this._type) > -1;
    }

    /** Checks whether the input is invalid based on the native validation. */
    protected _isBadInput(): boolean {
        // The `validity` property won't be present on platform-server.
        const validity = (this._elementRef.nativeElement as HTMLInputElement).validity;
        return validity && validity.badInput;
    }

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    get empty(): boolean {
        return !this._isNeverEmpty() && !this._elementRef.nativeElement.value && !this._isBadInput() &&
            !this.autofilled;
    }

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    get shouldLabelFloat(): boolean {
        if (this._isNativeSelect) {
            // For a single-selection `<select>`, the label should float when the selected option has
            // a non-empty display value. For a `<select multiple>`, the label *always* floats to avoid
            // overlapping the label with the options.
            const selectElement = this._elementRef.nativeElement as HTMLSelectElement;
            const firstOption: HTMLOptionElement | undefined = selectElement.options[0];

            // On most browsers the `selectedIndex` will always be 0, however on IE and Edge it'll be
            // -1 if the `value` is set to something, that isn't in the list of options, at a later point.
            return this.focused || selectElement.multiple || !this.empty ||
                !!(selectElement.selectedIndex > -1 && firstOption && firstOption.label);
        } else {
            return this.focused || !this.empty;
        }
    }

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    setDescribedByIds(ids: string[]): void {
        if (ids.length) {
            this._elementRef.nativeElement.setAttribute('aria-describedby', ids.join(' '));
        } else {
            this._elementRef.nativeElement.removeAttribute('aria-describedby');
        }
    }

    /**
     * Implemented as part of MatFormFieldControl.
     * @docs-private
     */
    onContainerClick(): void {
        // Do not re-focus the input element if the element is already focused. Otherwise it can happen
        // that someone clicks on a time input and the cursor resets to the "hours" field while the
        // "minutes" field was actually clicked. See: https://github.com/angular/components/issues/12849
        if (!this.focused) {
            this.focus();
        }
    }
    // tslint:disable
    static ngAcceptInputType_disabled: BooleanInput;
    static ngAcceptInputType_readonly: BooleanInput;
    static ngAcceptInputType_required: BooleanInput;

    // Accept `any` to avoid conflicts with other directives on `<input>` that may
    // accept different types.
    static ngAcceptInputType_value: any;
    // tslint:enable
}
