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
import { MatLegacyFormFieldControl as MatFormFieldControl } from '@angular/material/legacy-form-field';
import { Observable, of, Subject } from 'rxjs';
import { distinctUntilChanged, map, mergeMap, startWith, takeUntil } from 'rxjs/operators';

import { cnslFormFieldAnimations } from './animations';
import { CnslErrorDirective, CNSL_ERROR } from '../error/error.directive';
import { KeyValue, KeyValuePipe } from '@angular/common';

export const CNSL_FORM_FIELD = new InjectionToken<CnslFormFieldComponent>('CnslFormFieldComponent');

class CnslFormFieldBase {
  constructor(public _elementRef: ElementRef) {}
}

@Component({
  selector: 'cnsl-form-field',
  templateUrl: './form-field.component.html',
  styleUrls: ['./form-field.component.scss'],
  providers: [{ provide: CNSL_FORM_FIELD, useExisting: CnslFormFieldComponent }],
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
    return this._explicitFormFieldControl || this._controlNonStatic || this._controlStatic;
  }
  set _control(value: MatFormFieldControl<any>) {
    this._explicitFormFieldControl = value;
  }
  private _explicitFormFieldControl!: MatFormFieldControl<any>;
  readonly stateChanges: Subject<void> = new Subject<void>();
  public errori18nKeys$?: Observable<Array<string>>;

  _subscriptAnimationState: string = '';

  @ContentChildren(CNSL_ERROR as any, { descendants: true }) _errorChildren!: QueryList<CnslErrorDirective>;

  @HostListener('blur', ['false'])
  _focusChanged(isFocused: boolean): void {
    if (isFocused !== this.focused && !isFocused) {
      this.focused = isFocused;
      this.stateChanges.next();
    }
  }

  constructor(
    public _elementRef: ElementRef,
    private _changeDetectorRef: ChangeDetectorRef,
    @Inject(ElementRef)
    _labelOptions: // Use `ElementRef` here so Angular has something to inject.
    any,
    private kvPipe: KeyValuePipe,
  ) {
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
    this.defineI18nErrors()

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

      if (this._control.userAriaDescribedBy && typeof this._control.userAriaDescribedBy === 'string') {
        ids.push(...this._control.userAriaDescribedBy.split(' '));
      }

      if (this._errorChildren) {
        ids.push(...this._errorChildren.map((error) => error.id));
      }

      this._control.setDescribedByIds(ids);
    }
  }

  private defineI18nErrors(): void {
    let ctrl = this._control.ngControl?.control
    this.errori18nKeys$ = ctrl?.valueChanges?.pipe(
      mergeMap(() => ctrl?.statusChanges || of([])),
      map(() => this.currentErrors()),
      distinctUntilChanged(),
    ) || of([]);
  }

  private currentErrors(): Array<string> {
    return (
      this.kvPipe
        .transform(this._control.ngControl?.control?.errors)
        ?.filter(this.filterErrorsProperties)
        .map(this.mapErrorToI18nKey)
        .filter(this.distinctFilter) || []
    );
  }

  private filterErrorsProperties(kv: KeyValue<unknown, unknown>): boolean {
    return typeof kv.value == "object" && (kv.value as {valid: boolean}).valid === false
  }

  private mapErrorToI18nKey(kv: KeyValue<unknown, unknown>): string {
    return (kv.value as { i18nKey: string }).i18nKey || 'ERRORS.INVALID_FORMAT';
  }

  private distinctFilter(item: string, index: number, arr: Array<string>): boolean {
    return arr.indexOf(item) === index;
  }


  /** Determines whether a class from the NgControl should be forwarded to the host element. */
  _shouldForward(prop: keyof NgControl): boolean {
    const ngControl: any = this._control ? this._control.ngControl : null;
    return ngControl && ngControl[prop];
  }
}
