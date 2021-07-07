import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'cnsl-edit-text',
  templateUrl: './edit-text.component.html',
  styleUrls: ['./edit-text.component.scss']
})
export class EditTextComponent implements OnInit, OnDestroy {
  @Input() label: string = 'hello';
  @Input() current$!: Observable<{ [key: string]: string; }>;
  @Input() default$!: Observable<{ [key: string]: string; }>;
  @Output() changedValues: EventEmitter<{ [key: string]: string; }> = new EventEmitter();
  public currentMap: { [key: string]: string; } = {};
  private destroy$: Subject<void> = new Subject();
  public form!: FormGroup;
  constructor() { }

  public ngOnInit(): void {
    this.current$.pipe(takeUntil(this.destroy$)).subscribe(value => {
      console.log('current', value);
      this.currentMap = value;
      this.form = new FormGroup({});
      Object.keys(value).map(key => {
        const control = new FormControl(value[key]);
        this.form.addControl(key, control);
      });

      this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe(values => this.changedValues.emit(values));
    });
  }

  public ngOnDestroy(): void {
    console.log('destroy');
    this.destroy$.next();
    this.destroy$.complete();
  }
}
