import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'cnsl-edit-text',
  templateUrl: './edit-text.component.html',
  styleUrls: ['./edit-text.component.scss'],
})
export class EditTextComponent implements OnInit, OnDestroy {
  @Input() label: string = '';
  @Input() current$!: Observable<{ [key: string]: any | string; }>;
  @Input() default$!: Observable<{ [key: string]: any | string; }>;
  @Input() currentlyDragged: string = '';
  @Output() changedValues: EventEmitter<{ [key: string]: string; }> = new EventEmitter();
  public currentMap: { [key: string]: string; } = {};
  private destroy$: Subject<void> = new Subject();
  public form!: FormGroup;
  public warnText: { [key: string]: string | undefined; } = {};

  @Input() public chips: any[] = [];
  @Input() public disabled: boolean = true;

  public copied: string = '';

  public ngOnInit(): void {
    this.current$.pipe(takeUntil(this.destroy$)).subscribe(value => {
      this.currentMap = value;
      this.form = new FormGroup({});
      Object.keys(value).map(key => {
        const control = new FormControl({ value: value[key], disabled: this.disabled });
        this.form.addControl(key, control);
      });

      this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe(values => this.changedValues.emit(values));
    });

  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public setWarnText(key: string, text: string | undefined): void {
    this.warnText[key] = text;
  }

  public addChip(key: string, value: string): void {
    const c = this.form.get(key)?.value;
    this.form.get(key)?.setValue(`${c} ${value}`);
  }
}
