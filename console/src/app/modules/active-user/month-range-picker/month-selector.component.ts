import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  inject,
  input,
  Output,
} from "@angular/core";
import { MonthService } from "@/modules/active-user/month-range-picker/month.service";

@Component({
  selector: "cnsl-month-selector",
  templateUrl: "./month-selector.component.html",
  providers: [MonthService],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
})
export class MonthSelectorComponent {
  protected readonly monthService = inject(MonthService);

  public readonly selectedMonth = input.required<number | undefined>();
  // need to use old EventEmitter because signals swallow duplicates
  @Output()
  public readonly selectedMonthChange = new EventEmitter<number>();

  public readonly notBefore = input<number>(-1);
  public readonly notAfter = input<number>(12);
}
