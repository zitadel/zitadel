import { Component, forwardRef, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ControlValueAccessor, FormControl, NG_VALUE_ACCESSOR } from '@angular/forms';
import { Observable, Subject, takeUntil } from 'rxjs';

@Component({
  selector: 'cnsl-redirect-uris',
  templateUrl: './redirect-uris.component.html',
  styleUrls: ['./redirect-uris.component.scss'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => RedirectUrisComponent),
      multi: true,
    },
  ],
})
export class RedirectUrisComponent implements ControlValueAccessor, OnInit, OnDestroy {
  @Input() title: string = '';
  @Input() devMode: boolean = false;
  @Input() isNative!: boolean;
  @Input() public getValues: Observable<void> = new Observable(); // adds formfieldinput to array on emission

  public redirectControl: FormControl = new FormControl<string>({ value: '', disabled: true });
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
      this.redirectControl.disable();
    } else {
      this.redirectControl.enable();
    }
  }

  public add(input: any): void {
    if (this.redirectControl.valid) {
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

  public remove(redirect: any): void {
    const index = this.value.indexOf(redirect);

    if (index >= 0) {
      this.value.splice(index, 1);
      this.onChange(this.value);
      this.onTouch(this.value);
    }
  }
}
