import { Component, Input, OnDestroy, OnInit } from '@angular/core';
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

  public currentMap: { [key: string]: string; } = {};
  private destroy$: Subject<void> = new Subject();

  constructor() { }

  public ngOnInit(): void {
    this.current$.pipe(takeUntil(this.destroy$)).subscribe(value => {
      console.log('current', value);
      this.currentMap = value;
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
