import { CdkOverlayOrigin, OverlayModule } from "@angular/cdk/overlay";
import {
  ChangeDetectionStrategy,
  Component,
  contentChild,
  Directive,
  HostListener,
  inject,
  model,
  signal,
} from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MonthSelectorComponent } from "@/modules/active-user/month-range-picker/month-selector.component";
import { YearSelectorComponent } from "@/modules/active-user/month-range-picker/year-selector.component";
import { MatAutocompleteModule } from "@angular/material/autocomplete";
import { InputModule } from "@/modules/input/input.module";
import { heroChevronDown } from "@ng-icons/heroicons/outline";
import { NgIcon, provideIcons } from "@ng-icons/core";

type State =
  | (
      | { page: "start" }
      | { page: "end"; endMonth?: number; endYear: number }
    ) & {
      selector: "month" | "year";
      startMonth: number;
      startYear: number;
    };

@Component({
  selector: "cnsl-month-range-picker",
  templateUrl: "./month-range-picker.component.html",
  styleUrls: ["./month-range-picker.component.scss"],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  providers: [provideIcons({ heroChevronDown })],
  imports: [
    OverlayModule,
    MatAutocompleteModule,
    InputModule,
    ReactiveFormsModule,
    MonthSelectorComponent,
    YearSelectorComponent,
    NgIcon,
  ],
})
export class MonthRangePickerComponent {
  protected readonly overlayOrigin = inject(CdkOverlayOrigin, { host: true });

  public readonly start = model.required<Date>();
  public readonly end = model.required<Date>();
  protected readonly state = signal<State | undefined>(undefined);

  protected readonly currentDate = new Date();

  protected modifyStart(
    state: Extract<State, { page: "start" }>,
    startMonth: number,
    startYear: number
  ) {
    return {
      ...state,
      startMonth,
      startYear,
    };
  }

  protected modifyEnd(state: Extract<State, { page: "end" }>, endYear: number) {
    return {
      ...state,
      endYear,
    };
  }

  protected apply(state: Extract<State, { page: "end" }>, endMonth: number) {
    // emit the newly selected date
    this.start.set(new Date(Date.UTC(state.startYear, state.startMonth, 1)));
    this.end.set(new Date(Date.UTC(state.endYear, endMonth + 1, 0)));
    // close the picker
    this.toggle();
  }

  public toggle() {
    const start = this.start();

    this.state.update((state) => {
      if (state) {
        return undefined;
      }

      return {
        page: "start",
        selector: "month",
        startMonth: start.getMonth(),
        startYear: start.getFullYear(),
      };
    });
  }

  protected toggleSelector(state: State) {
    this.state.set({
      ...state,
      selector: state.selector === "month" ? "year" : "month",
    });
  }

  protected selectEnd(state: Extract<State, { page: "start" }>) {
    // we either use the previously selected year or if the start year is now after it, we use the start year
    const endYear = Math.max(this.end().getUTCFullYear(), state.startYear);

    const endMonth =
      // we can only keep the previously selected month if the end date is after the start date
      this.end().getUTCMonth() > state.startMonth || endYear > state.startYear
        ? this.end().getUTCMonth()
        : undefined;

    this.state.set({
      ...state,
      page: "end" as const,
      endYear,
      endMonth,
    });
  }
}

@Directive({
  selector: "[cnslMonthRangePickerTrigger]",
  hostDirectives: [CdkOverlayOrigin],
  standalone: true,
})
export class MonthRangePickerTriggerDirective {
  private readonly component = contentChild(MonthRangePickerComponent);

  @HostListener("click")
  public toggle() {
    this.component()?.toggle();
  }
}
