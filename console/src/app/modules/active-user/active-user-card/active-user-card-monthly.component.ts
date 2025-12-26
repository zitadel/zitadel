import {
  ChangeDetectionStrategy,
  Component,
  effect,
  input,
  model,
  signal,
} from "@angular/core";
import {
  ActiveUserCardComponent,
  AgChartOptions,
} from "@/modules/active-user/active-user-card/active-user-card.component";
import { MatFormFieldModule } from "@angular/material/form-field";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import {
  MatButtonToggle,
  MatButtonToggleGroup,
} from "@angular/material/button-toggle";
import { MatDatepickerModule } from "@angular/material/datepicker";
import {
  MonthRangePickerComponent,
  MonthRangePickerTriggerDirective,
} from "@/modules/active-user/month-range-picker/month-range-picker.component";

@Component({
  selector: "cnsl-active-user-card-monthly",
  templateUrl: "./active-user-card-monthly.component.html",
  standalone: true,
  imports: [
    MatFormFieldModule,
    MatDatepickerModule,
    ReactiveFormsModule,
    MatButtonToggleGroup,
    MatButtonToggle,
    ActiveUserCardComponent,
    MonthRangePickerComponent,
    MonthRangePickerTriggerDirective,
    FormsModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActiveUserCardMonthlyComponent {
  public readonly state = signal<3 | 6 | 12 | "custom">(3);

  public readonly start = model.required<Date>();
  public readonly end = model.required<Date>();

  public readonly chart = input<AgChartOptions>();

  constructor() {
    effect(() => {
      const state = this.state();
      if (state == "custom") {
        return;
      }

      const start = new Date();
      start.setMonth(start.getMonth() - state);
      this.start.set(start);

      this.end.set(new Date());
    });
  }
}
