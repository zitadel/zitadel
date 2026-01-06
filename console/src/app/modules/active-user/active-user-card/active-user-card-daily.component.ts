import {
  ChangeDetectionStrategy,
  Component,
  computed,
  effect,
  inject,
  input,
  model,
  signal,
  untracked,
} from "@angular/core";
import {
  ActiveUserCardComponent,
  AgChartOptions,
} from "@/modules/active-user/active-user-card/active-user-card.component";
import { MatFormFieldModule } from "@angular/material/form-field";
import {
  FormBuilder,
  FormControl,
  FormsModule,
  ReactiveFormsModule,
} from "@angular/forms";
import {
  MatButtonToggle,
  MatButtonToggleGroup,
} from "@angular/material/button-toggle";
import { MatDatepickerModule } from "@angular/material/datepicker";
import { takeUntilDestroyed, toObservable } from "@angular/core/rxjs-interop";
import { switchMap } from "rxjs";
import { provideNativeDateAdapter } from "@angular/material/core";

@Component({
  selector: "cnsl-active-user-card-daily",
  templateUrl: "./active-user-card-daily.component.html",
  standalone: true,
  imports: [
    MatFormFieldModule,
    MatDatepickerModule,
    ReactiveFormsModule,
    MatButtonToggleGroup,
    MatButtonToggle,
    ActiveUserCardComponent,
    FormsModule,
  ],
  providers: [provideNativeDateAdapter()],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActiveUserCardDailyComponent {
  public readonly form = this.buildForm();
  public readonly state = signal<7 | 30 | "custom">(7);

  public readonly start = model.required<Date>();
  public readonly end = model.required<Date>();

  public readonly chart = input<AgChartOptions>();

  constructor() {
    effect(() => {
      const state = this.state();
      if (state === "custom") {
        return;
      }

      const start = new Date();
      start.setDate(start.getDate() - state);
      this.start.set(start);

      this.end.set(new Date());
    });
  }

  private buildForm() {
    const fb = inject(FormBuilder);

    // we need to create the form inside computed because the models
    // can only be read once the inputs are set
    // to make sure we don't create new form on each change we use untracked
    const form = computed(() =>
      untracked(() =>
        fb.group({
          start: new FormControl<Date>(this.start()),
          end: new FormControl<Date>(this.end()),
        })
      )
    );

    // sync form values when models change
    effect(() => {
      const start = this.start();
      const end = this.end();

      form().setValue({
        start,
        end,
      });
    });

    // update models when form values change
    toObservable(form)
      .pipe(
        switchMap((form) => form.valueChanges),
        takeUntilDestroyed()
      )
      .subscribe(({ start, end }) => {
        if (start && end && start <= end) {
          this.start.set(start);
          this.end.set(end);
        }
      });

    return form;
  }
}
