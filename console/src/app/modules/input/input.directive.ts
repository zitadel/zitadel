import { coerceBooleanProperty } from '@angular/cdk/coercion';
import { Directive, HostListener, Input } from '@angular/core';
import { NgControl } from '@angular/forms';
import { Subject } from 'rxjs';

let nextUniqueId = 0;

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
})
export class InputDirective {
    /** Whether the element is readonly. */
    // @Input()
    // get readonly(): boolean { return this._readonly; }
    // set readonly(value: boolean) { this._readonly = coerceBooleanProperty(value); }
    // private _readonly = false;
    // readonly id!: string;
    protected _uid: string = `mat-input-${nextUniqueId++}`;

    readonly ngControl!: NgControl | null;
    protected _id!: string;
    protected _required: boolean = false;

    protected _disabled: boolean = false;

    focused: boolean = false;
    readonly stateChanges: Subject<void> = new Subject<void>();

    /**
    * Implemented as part of MatFormFieldControl.
    * @docs-private
    */
    @Input() placeholder!: string;

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


    constructor() {
        console.log('input');
    }

    // onContainerClick() {
    //     // Do not re-focus the input element if the element is already focused. Otherwise it can happen
    //     // that someone clicks on a time input and the cursor resets to the "hours" field while the
    //     // "minutes" field was actually clicked. See: https://github.com/angular/components/issues/12849
    //     if (!this.focused) {
    //       this.focus();
    //     }
    //   }
}
