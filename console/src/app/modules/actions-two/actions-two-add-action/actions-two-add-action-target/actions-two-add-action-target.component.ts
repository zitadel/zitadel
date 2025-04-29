import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  computed,
  effect,
  EventEmitter,
  Input,
  Output,
  signal,
  Signal,
} from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ReplaySubject, switchMap } from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { InputModule } from 'src/app/modules/input/input.module';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MessageInitShape } from '@bufbuild/protobuf';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import { ExecutionTargetTypeSchema } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { MatSelectModule } from '@angular/material/select';
import { ActionConditionPipeModule } from 'src/app/pipes/action-condition-pipe/action-condition-pipe.module';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { startWith } from 'rxjs/operators';
import { TypeSafeCellDefModule } from 'src/app/directives/type-safe-cell-def/type-safe-cell-def.module';
import { CdkDrag, CdkDragDrop, CdkDropList, moveItemInArray } from '@angular/cdk/drag-drop';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { minArrayLengthValidator } from '../../../form-field/validators/validators';
import { ProjectRoleChipModule } from '../../../project-role-chip/project-role-chip.module';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TableActionsModule } from '../../../table-actions/table-actions.module';

@Component({
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  selector: 'cnsl-actions-two-add-action-target',
  templateUrl: './actions-two-add-action-target.component.html',
  styleUrls: ['./actions-two-add-action-target.component.scss'],
  imports: [
    TranslateModule,
    MatRadioModule,
    RouterModule,
    ReactiveFormsModule,
    InputModule,
    MatAutocompleteModule,
    FormsModule,
    ActionConditionPipeModule,
    CommonModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatTableModule,
    TypeSafeCellDefModule,
    CdkDrag,
    CdkDropList,
    ProjectRoleChipModule,
    MatTooltipModule,
    TableActionsModule,
  ],
})
export class ActionsTwoAddActionTargetComponent {
  @Input() public hideBackButton = false;
  @Input()
  public set preselectedTargetIds(preselectedTargetIds: string[]) {
    this.preselectedTargetIds$.next(preselectedTargetIds);
  }

  @Output() public readonly back = new EventEmitter<void>();
  @Output() public readonly continue = new EventEmitter<MessageInitShape<typeof ExecutionTargetTypeSchema>[]>();

  private readonly preselectedTargetIds$ = new ReplaySubject<string[]>(1);

  protected readonly form: ReturnType<typeof this.buildForm>;
  protected readonly targets: ReturnType<typeof this.listTargets>;
  private readonly selectedTargetIds: Signal<string[]>;
  protected readonly selectableTargets: Signal<Target[]>;
  protected readonly dataSource: MatTableDataSource<Target>;

  constructor(
    private readonly fb: FormBuilder,
    private readonly actionService: ActionService,
    private readonly toast: ToastService,
  ) {
    this.form = this.buildForm();
    this.targets = this.listTargets();

    this.selectedTargetIds = this.getSelectedTargetIds(this.form);
    this.selectableTargets = this.getSelectableTargets(this.targets, this.selectedTargetIds);
    this.dataSource = this.getDataSource(this.targets, this.selectedTargetIds);
  }

  private buildForm() {
    const preselectedTargetIds = toSignal(this.preselectedTargetIds$, { initialValue: [] as string[] });

    return computed(() => {
      return this.fb.group({
        autocomplete: new FormControl('', { nonNullable: true }),
        selectedTargetIds: new FormControl(preselectedTargetIds(), {
          nonNullable: true,
          validators: [minArrayLengthValidator(1)],
        }),
      });
    });
  }

  private listTargets() {
    const targetsSignal = signal({ state: 'loading' as 'loading' | 'loaded', targets: new Map<string, Target>() });

    this.actionService
      .listTargets({})
      .then(({ result }) => {
        const targets = result.reduce((acc, target) => {
          acc.set(target.id, target);
          return acc;
        }, new Map<string, Target>());

        targetsSignal.set({ state: 'loaded', targets });
      })
      .catch((error) => {
        this.toast.showError(error);
      });

    return computed(targetsSignal);
  }

  private getSelectedTargetIds(form: typeof this.form) {
    const selectedTargetIds$ = toObservable(form).pipe(
      startWith(form()),
      switchMap((form) => {
        const { selectedTargetIds } = form.controls;
        return selectedTargetIds.valueChanges.pipe(startWith(selectedTargetIds.value));
      }),
    );
    return toSignal(selectedTargetIds$, { requireSync: true });
  }

  private getSelectableTargets(targets: typeof this.targets, selectedTargetIds: Signal<string[]>) {
    return computed(() => {
      const targetsCopy = new Map(targets().targets);
      for (const selectedTargetId of selectedTargetIds()) {
        targetsCopy.delete(selectedTargetId);
      }
      return Array.from(targetsCopy.values());
    });
  }

  private getDataSource(targetsSignal: typeof this.targets, selectedTargetIdsSignal: Signal<string[]>) {
    const selectedTargets = computed(() => {
      // get this out of the loop so angular can track this dependency
      // even if targets is empty
      const { targets, state } = targetsSignal();
      const selectedTargetIds = selectedTargetIdsSignal();

      if (state === 'loading') {
        return [];
      }

      return selectedTargetIds.map((id) => {
        const target = targets.get(id);
        if (!target) {
          throw new Error(`Target with id ${id} not found`);
        }
        return target;
      });
    });

    const dataSource = new MatTableDataSource<Target>(selectedTargets());
    effect(() => {
      dataSource.data = selectedTargets();
    });

    return dataSource;
  }

  protected addTarget(target: Target) {
    const { selectedTargetIds } = this.form().controls;
    selectedTargetIds.setValue([target.id, ...selectedTargetIds.value]);
    this.form().controls.autocomplete.setValue('');
  }

  protected removeTarget(index: number) {
    const { selectedTargetIds } = this.form().controls;
    const data = [...selectedTargetIds.value];
    data.splice(index, 1);
    selectedTargetIds.setValue(data);
  }

  protected drop(event: CdkDragDrop<undefined>) {
    const { selectedTargetIds } = this.form().controls;

    const data = [...selectedTargetIds.value];
    moveItemInArray(data, event.previousIndex, event.currentIndex);
    selectedTargetIds.setValue(data);
  }

  protected handleEnter(event: Event) {
    const selectableTargets = this.selectableTargets();
    if (selectableTargets.length !== 1) {
      return;
    }

    event.preventDefault();
    this.addTarget(selectableTargets[0]);
  }

  protected submit() {
    const selectedTargets = this.selectedTargetIds().map((value) => ({
      type: {
        case: 'target' as const,
        value,
      },
    }));

    this.continue.emit(selectedTargets);
  }

  protected trackTarget(_: number, target: Target) {
    return target.id;
  }
}
