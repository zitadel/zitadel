import {Component, EventEmitter, forwardRef, Input, Output} from "@angular/core";
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from "@angular/forms";

@Component({
  selector: 'cnsl-datetime-local-input',
  template: `<input type="datetime-local" [value] = "_date"
             (change) = "onDateChange($any($event.target).value)" />`,
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => DateTimeLocalInputComponent),
      multi: true,
    }
  ]

})
export class DateTimeLocalInputComponent implements ControlValueAccessor{
  public _date: string = "";
  private onChange!: (value: Date) => void;
  private onTouched!: () => void;

  @Input() set date(d: Date) {
    this._date = this.toDateString(d);
  }
  @Output() dateChange: EventEmitter<Date>;
  constructor() {
    this.date = new Date();
    this.dateChange = new EventEmitter();
  }

  private toDateString(date: Date): string {
    return (date.getFullYear().toString() + '-'
        + ("0" + (date.getMonth() + 1)).slice(-2) + '-'
        + ("0" + (date.getDate())).slice(-2))
      + 'T' + date.toTimeString().slice(0,5);
  }

  private parseDateString(date:string): Date {
    date = date.replace('T','-');
    var parts = date.split('-');
    var timeParts = parts[3].split(':');
    return new Date(parseInt(parts[0]), parseInt(parts[1])-1, parseInt(parts[2]), parseInt(timeParts[0]), parseInt(timeParts[1]));
  }

  public onDateChange(value: any): void {
    if (value != this._date) {
      var parsedDate = this.parseDateString(value);
        this._date = value;
        this.dateChange.emit(parsedDate);
    }
  }

  registerOnChange(fn: (value: Date) => void): void {
    this.onChange = fn
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn
  }

  writeValue(obj: Date): void {
    this._date = this.toDateString(obj)
  }
}
