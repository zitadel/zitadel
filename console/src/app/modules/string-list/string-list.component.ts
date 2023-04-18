import { Component, forwardRef, Input, OnDestroy, ViewChildren, ViewEncapsulation } from '@angular/core';
import { ControlValueAccessor, FormControl, NG_VALUE_ACCESSOR } from '@angular/forms';
import { Subject } from 'rxjs';
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

  @Input() public inputControl: FormControl = new FormControl<string>('', [requiredValidator]);

  private destroy$: Subject<void> = new Subject();
  @ViewChildren('stringInput') input!: any[];
  public val: string[] = [];

  constructor() {
    this.control.setValidators([minArrayLengthValidator]);
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

  setValueAtIndex(index: number, event: any) {
    const value = event?.target?.value ?? event.value;
    const toSet = value.trim();
    console.log(toSet);
    this.value[index] = toSet;
  }

  addArrayEntry() {
    this.value.push('');
  }

  removeEntryAtIndex(index: number) {
    this.value.splice(index, 1);
  }

  get value() {
    return this.val;
  }

  writeValue(value: string[]) {
    this.value = value;
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

  public add(input: any): void {
    if (this.control.valid) {
      const trimmed = input.value.trim();
      if (trimmed) {
        this.val ? this.val.push(input.value) : (this.val = [input.value]);
        this.onChange(this.val);
        this.onTouch(this.val);
      }
      if (input) {
        input.value = '';
      }
    }
  }

  public remove(str: string): void {
    const index = this.value.indexOf(str);

    if (index >= 0) {
      this.value.splice(index, 1);
      this.onChange(this.value);
      this.onTouch(this.value);
    }
  }
}
