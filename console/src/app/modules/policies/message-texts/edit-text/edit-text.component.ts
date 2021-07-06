import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'cnsl-edit-text',
  templateUrl: './edit-text.component.html',
  styleUrls: ['./edit-text.component.scss']
})
export class EditTextComponent implements OnInit, OnDestroy {
  @Input() label: string = 'hello';
  @Input() current$: Observable<string> = of('');
  @Input() default$: Observable<string> = of('');

  public value: string = '';
  public default: string = '';
  private destroy$: Subject<void> = new Subject();

  constructor() { }

  public ngOnInit(): void {
    this.current$.pipe(takeUntil(this.destroy$)).subscribe(value => {
      this.value = value;
    });
    this.default$.pipe(takeUntil(this.destroy$)).subscribe(value => {
      this.default = value;
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
