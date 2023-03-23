import { Component, forwardRef, Input, OnDestroy, OnInit, ViewChild, ViewEncapsulation } from '@angular/core';
import { ControlValueAccessor, FormControl, NG_VALUE_ACCESSOR } from '@angular/forms';
import { Observable, Subject, takeUntil } from 'rxjs';

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
export class StringListComponent implements ControlValueAccessor, OnInit, OnDestroy {
  @Input() title: string = '';
  @Input() public getValues: Observable<void> = new Observable(); // adds formfieldinput to array on emission

  @Input() public control: FormControl = new FormControl<string>({ value: '', disabled: true });
  private destroy$: Subject<void> = new Subject();
  @ViewChild('redInput') input!: any;

  ngOnInit(): void {
    this.getValues.pipe(takeUntil(this.destroy$)).subscribe(() => {
      this.add(this.input.nativeElement);
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onChange: any = () => {};
  onTouch: any = () => {};

  private val: string[] = [];

  set value(val: string[]) {
    console.log('setvalue', val);
    if (val !== undefined && this.val !== val) {
      this.val = val;
      this.onChange(val);
      this.onTouch(val);
    }
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
    console.log(input.value);
    if (this.control.valid) {
      if (input.value !== '' && input.value !== ' ' && input.value !== '/') {
        this.val.push(input.value);
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
