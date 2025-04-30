import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Observable, Subject, map, of, startWith, switchMap, tap } from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ConditionType } from '../actions-two-add-action-condition/actions-two-add-action-condition.component';

// export enum ExecutionType {
//   REQUEST = 'request',
//   RESPONSE = 'response',
//   EVENTS = 'event',
//   FUNCTIONS = 'function',
// }

@Component({
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  selector: 'cnsl-actions-two-add-action-type',
  templateUrl: './actions-two-add-action-type.component.html',
  styleUrls: ['./actions-two-add-action-type.component.scss'],
  imports: [TranslateModule, MatRadioModule, RouterModule, ReactiveFormsModule, FormsModule, CommonModule, MatButtonModule],
})
export class ActionsTwoAddActionTypeComponent {
  protected readonly typeForm: ReturnType<typeof this.buildActionTypeForm> = this.buildActionTypeForm();
  @Output() public readonly typeChanges$: Observable<ConditionType>;

  @Output() public readonly back = new EventEmitter<void>();
  @Output() public readonly continue = new EventEmitter<ConditionType>();
  @Input() public set initialValue(type: ConditionType) {
    this.typeForm.get('executionType')!.setValue(type);
  }

  constructor(private readonly fb: FormBuilder) {
    this.typeChanges$ = this.typeForm.get('executionType')!.valueChanges.pipe(
      startWith(this.typeForm.get('executionType')!.value), // Emit the initial value
    );
  }

  public buildActionTypeForm() {
    return this.fb.group({
      executionType: new FormControl<ConditionType>('request', {
        nonNullable: true,
      }),
    });
  }

  public submit() {
    this.continue.emit(this.typeForm.get('executionType')!.value);
  }
}
