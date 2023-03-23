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
  Input,
  OnDestroy,
  QueryList,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { NgControl } from '@angular/forms';
import { MatLegacyFormFieldControl as MatFormFieldControl } from '@angular/material/legacy-form-field';
import { combineLatest, map, Observable, of, startWith, Subject, takeUntil } from 'rxjs';

import { CnslErrorDirective, CNSL_ERROR } from '../error/error.directive';
import { cnslFormFieldAnimations } from './animations';

export const CNSL_FORM_FIELD = new InjectionToken<CnslFormFieldComponent>('CnslFormFieldComponent');

class CnslFormFieldBase {
  constructor(public _elementRef: ElementRef) {}
}

interface Help {
  type: 'hints' | 'errors';
  validationErrors?: Array<ValidationError>;
}
interface ValidationError {
  i18nKey: string;
  params: any;
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
  @Input() public disableValidationErrors = false;

  get _control(): MatFormFieldControl<any> {
    return this._explicitFormFieldControl || this._controlNonStatic || this._controlStatic;
  }
  set _control(value: MatFormFieldControl<any>) {
    this._explicitFormFieldControl = value;
  }

  private _explicitFormFieldControl!: MatFormFieldControl<any>;
  readonly stateChanges: Subject<void> = new Subject<void>();
  public help$?: Observable<Help>;

  _subscriptAnimationState: string = '';

  @ContentChildren(CNSL_ERROR as any, { descendants: true }) _errorChildren!: QueryList<CnslErrorDirective>;

  // TODO: Remove?
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
  ) {
    super(_elementRef);
  }

  public ngAfterViewInit(): void {
    this._changeDetectorRef.detectChanges();
  }

  public ngOnDestroy(): void {
    this._destroyed.next();
    this._destroyed.complete();
  }

  public ngAfterContentInit(): void {
    this._validateControlChild();
    this.mapHelp();

    const control = this._control;
    control.stateChanges.pipe(startWith(null), takeUntil(this._destroyed)).subscribe(() => {
      this._syncDescribedByIds();
      this._changeDetectorRef.markForCheck();
    });

    // Run change detection if the value changes.
    // TODO: Is that not redundant (see immediately above)?
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

  private mapHelp(): void {
    const validationErrors$: Observable<Array<ValidationError>> = this.disableValidationErrors
      ? of([])
      : this._control.stateChanges?.pipe(
          map(() => this.currentValidationErrors()),
          startWith([]),
        ) || of([]);

    const childrenErrors$: Observable<boolean> = this._errorChildren.changes.pipe(
      map(() => {
        return this._errorChildren.length > 0;
      }),
      startWith(false),
    );

    this.help$ = combineLatest([validationErrors$, childrenErrors$]).pipe(
      map((combined) => {
        return combined[0].length >= 1
          ? <Help>{
              type: 'errors',
              validationErrors: combined[0],
            }
          : combined[1]
          ? <Help>{
              type: 'errors',
            }
          : <Help>{
              type: 'hints',
              validationErrors: undefined,
            };
      }),
    );
    this._changeDetectorRef.markForCheck();
  }

  private currentValidationErrors(): Array<ValidationError> {
    return (
      Object.entries(this._control.ngControl?.control?.errors || [])
        ?.filter(this.filterErrorsProperties)
        .map(this.mapToValidationError)
        .filter(this.distinctFilter) || []
    );
  }

  private filterErrorsProperties(kv: [string, any]): boolean {
    return typeof kv[1] == 'object' && (kv[1] as { valid: boolean }).valid === false;
  }

  private mapToValidationError(kv: [string, any]): ValidationError {
    return {
      i18nKey: 'ERRORS.INVALID_FORMAT',
      ...(kv[1] as ValidationError | any),
    };
  }

  private distinctFilter(_: ValidationError, index: number, arr: Array<ValidationError>): boolean {
    return arr.findIndex((item) => item.i18nKey) === index;
  }

  /** Determines whether a class from the NgControl should be forwarded to the host element. */
  _shouldForward(prop: keyof NgControl): boolean {
    const ngControl: any = this._control ? this._control.ngControl : null;
    return ngControl && ngControl[prop];
  }
}
