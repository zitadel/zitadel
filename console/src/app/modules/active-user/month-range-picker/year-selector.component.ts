import {
  ChangeDetectionStrategy,
  Component,
  computed,
  EventEmitter,
  input,
  Output,
} from "@angular/core";

@Component({
  selector: "cnsl-year-selector",
  templateUrl: "./year-selector.component.html",
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class YearSelectorComponent {
  public readonly range = input<number>(8);
  public readonly notBefore = input<number>(-1);

  public readonly selectedYear = input.required<number>();
  // need to use old EventEmitter because signals swallow duplicates
  @Output()
  public readonly selectedYearChange = new EventEmitter<number>();

  protected readonly years = computed(() => {
    const currentYear = new Date().getFullYear();
    return Array.from(
      { length: this.range() },
      (_, i) => currentYear - i
    ).reverse();
  });
}
