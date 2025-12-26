import { Component, ChangeDetectionStrategy, input } from "@angular/core";
import { CardModule } from "@/modules/card/card.module";
import { AgCharts } from "ag-charts-angular";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";

export type AgChartOptions = AgCharts["options"];

@Component({
  selector: "cnsl-active-user-card",
  templateUrl: "./active-user-card.component.html",
  imports: [CardModule, AgCharts, MatProgressSpinnerModule],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActiveUserCardComponent {
  public readonly title = input.required<string>();
  public readonly chart = input<AgChartOptions>();
}
