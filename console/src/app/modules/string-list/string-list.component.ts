import { Component, forwardRef, Input, OnDestroy, ViewChildren, ViewEncapsulation } from '@angular/core';
import { ControlValueAccessor, FormArray, FormControl, NG_VALUE_ACCESSOR } from '@angular/forms';
import { distinctUntilChanged, Subject, takeUntil } from 'rxjs';
import { minArrayLengthValidator, requiredValidator } from '../form-field/validators/validators';

@Component({
  selector: 'cnsl-string-list',
  templateUrl: './string-list.component.html',
  styleUrls: ['./string-list.component.scss'],
  encapsulation: ViewEncapsulation.None,
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => StringListComponent),
      multi: true,
    },
  ],
})
export class StringListComponent implements ControlValueAccessor, OnDestroy {
  @Input() title: string = '';
  @Input() required: boolean = false;

  @Input() public control: FormControl = new FormControl<string[]>({ value: [], disabled: true });

  private destroy$: Subject<void> = new Subject();
  @ViewChildren('stringInput') input!: any[];
  public val: string[] = [];

  public formArray: FormArray = new FormArray([new FormControl('', [requiredValidator])]);

  constructor() {
    this.control.setValidators([minArrayLengthValidator(1)]);
    this.formArray.valueChanges.pipe(takeUntil(this.destroy$), distinctUntilChanged()).subscribe((value) => {
      this.value = value;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onChange: any = () => {};
  onTouch: any = () => {};

  set value(val: string[]) {
    if (val !== undefined && this.val !== val) {
      this.val = val;
      this.onChange(val);
      this.onTouch(val);
    }
  }

  addArrayEntry() {
    this.formArray.push(new FormControl('', [requiredValidator]));
  }

  removeEntryAtIndex(index: number) {
    this.formArray.removeAt(index);
  }

  clearEntryAtIndex(index: number) {
    this.formArray.controls[index].setValue('');
  }
  get value() {
    return this.val;
  }

  writeValue(value: string[]) {
    this.value = value;
    value.map((v, i) => this.formArray.setControl(i, new FormControl(v, [requiredValidator])));
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched(fn: any) {
    this.onTouch = fn;
  }

  public setDisabledState(isDisabled: boolean): void {
    if (isDisabled) {
      this.control.disable();
    } else {
      this.control.enable();
    }
  }
}
