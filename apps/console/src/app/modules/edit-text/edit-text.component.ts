import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

import { InfoSectionType } from '../info-section/info-section.component';

@Component({
  selector: 'cnsl-edit-text',
  templateUrl: './edit-text.component.html',
  styleUrls: ['./edit-text.component.scss'],
})
export class EditTextComponent implements OnInit, OnDestroy {
  @Input() label: string = '';
  @Input() current$!: Observable<{ [key: string]: string | boolean }>;
  @Input() default$!: Observable<{ [key: string]: string | boolean }>;
  @Input() currentlyDragged: string = '';
  @Output() changedValues: EventEmitter<{ [key: string]: string }> = new EventEmitter();
  public currentMap: { [key: string]: string | boolean } = {}; // boolean because of isDefault
  private destroy$: Subject<void> = new Subject();
  public form!: UntypedFormGroup;
  public warnText: { [key: string]: string | boolean | undefined } = {};

  @Input() public chips: any[] = [];
  @Input() public disabled: boolean = true;

  public copied: string = '';
  public InfoSectionType: any = InfoSectionType;

  public ngOnInit(): void {
    this.current$.pipe(takeUntil(this.destroy$)).subscribe((value) => {
      this.currentMap = value;
      this.form = new UntypedFormGroup({});
      Object.keys(value).map((key) => {
        if (key !== 'isDefault') {
          const control = new UntypedFormControl({ value: value[key], disabled: this.disabled });
          this.form.addControl(key, control);
        }
      });

      this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe((values) => this.changedValues.emit(values));
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public setWarnText(key: string, text: string | boolean | undefined): void {
    this.warnText[key] = text;
  }

  public addChip(key: string, value: string): void {
    const c = this.form.get(key)?.value;
    this.form.get(key)?.setValue(`${c} ${value}`);
  }
}
