import { ChangeDetectionStrategy, Component, computed, effect, inject } from '@angular/core';
import { ActiveUserService, averageActiveUserEntries } from './active-user.service';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormBuilder, FormControl, ReactiveFormsModule } from '@angular/forms';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { provideNativeDateAdapter } from '@angular/material/core';
import { ToastService } from 'src/app/services/toast.service';
import { ActiveUserGrpcMockProviderService, ActiveUserGrpcProviderService } from './active-user-grpc-provider.service';

@Component({
  templateUrl: './active-user-page.component.html',
  imports: [MatFormFieldModule, MatDatepickerModule, ReactiveFormsModule],
  providers: [
    provideNativeDateAdapter(),
    ActiveUserService,
    { provide: ActiveUserGrpcProviderService, useClass: ActiveUserGrpcMockProviderService },
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActiveUserPageComponent {
  private readonly toastService = inject(ToastService);
  private readonly activeUserService = inject(ActiveUserService);
  private readonly formBuilder = inject(FormBuilder);

  protected readonly form: ReturnType<typeof this.buildForm>;
  protected readonly activeUser: ReturnType<typeof this.getActiveUser>;

  constructor() {
    this.form = this.buildForm();
    this.activeUser = this.getActiveUser(this.form);

    effect(() => {
      const error = this.activeUser.error();
      if (error) {
        this.toastService.showError(error);
      }
    });
  }

  private buildForm() {
    return this.formBuilder.group({
      start: new FormControl<Date | undefined>(undefined),
      end: new FormControl<Date | undefined>(undefined),
    });
  }

  private getActiveUser(form: typeof this.form) {
    const formValuesSignal = toSignal(form.valueChanges, { initialValue: form.value });
    const activeUserRequest = computed(() => {
      const { start, end } = formValuesSignal();
      if (!start || !end) {
        return undefined;
      }
      return {
        precision: 'dailyPrecision',
        startingDateInclusive: start,
        endingDateInclusive: end,
      } as const;
    });

    return injectQuery(() => ({
      ...this.activeUserService.getActiveUser(activeUserRequest()),
      select: (entries) => averageActiveUserEntries(entries),
    }));
  }
}
