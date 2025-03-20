import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Observable, map, of, startWith, switchMap, tap } from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';

export enum ExecutionType {
  REQUEST,
  RESPONSE,
  EVENTS,
  FUNCTIONS,
}

@Component({
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  selector: 'cnsl-actions-two-add-action-type',
  templateUrl: './actions-two-add-action-type.component.html',
  styleUrls: ['./actions-two-add-action-type.component.scss'],
  imports: [TranslateModule, MatRadioModule, RouterModule, ReactiveFormsModule, FormsModule, CommonModule, MatButtonModule],
})
export class ActionsTwoAddActionTypeComponent implements OnInit {
  public ExecutionType = ExecutionType;
  protected readonly typeForm: ReturnType<typeof this.buildActionTypeForm> = this.buildActionTypeForm();
  @Output() public continue: EventEmitter<void> = new EventEmitter();
  @Output() public typeChanges$: Observable<ExecutionType>;

  constructor(private readonly fb: FormBuilder) {
    // Initialize the Observable to emit form value changes
    this.typeChanges$ = this.typeForm.get('executionType')!.valueChanges.pipe(
      startWith(this.typeForm.get('executionType')!.value), // Emit the initial value
      tap((value) => console.log('ExecutionType changed:', value)), // Debugging/logging
    );
  }

  public ngOnInit(): void {}

  public buildActionTypeForm() {
    return this.fb.group({
      executionType: new FormControl<ExecutionType>(ExecutionType.REQUEST, {
        nonNullable: true,
      }),
    });
  }
}
