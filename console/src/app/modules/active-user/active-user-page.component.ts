import {
  ChangeDetectionStrategy,
  Component,
  computed,
  effect,
  ElementRef,
  inject,
  Signal,
  signal,
} from "@angular/core";
import { ActiveUserService, PrecisionType } from "./active-user.service";
import { injectQuery } from "@tanstack/angular-query-experimental";
import { ToastService } from "src/app/services/toast.service";
import {
  ActiveUserGrpcMockProviderService,
  ActiveUserGrpcProviderService,
} from "./active-user-grpc-provider.service";
import { CardModule } from "@/modules/card/card.module";
import { ActiveUserCardDailyComponent } from "@/modules/active-user/active-user-card/active-user-card-daily.component";
import { ActiveUserCardMonthlyComponent } from "@/modules/active-user/active-user-card/active-user-card-monthly.component";
import { GrpcService } from "@/services/grpc.service";
import { TimestampToDatePipe } from "@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe";
import { AgChartOptions } from "@/modules/active-user/active-user-card/active-user-card.component";
import { ThemeService } from "@/services/theme.service";
import { toSignal } from "@angular/core/rxjs-interop";
import { map } from "rxjs/operators";
import { delay } from "rxjs";

@Component({
  templateUrl: "./active-user-page.component.html",
  imports: [
    CardModule,
    ActiveUserCardDailyComponent,
    ActiveUserCardMonthlyComponent,
  ],
  providers: [
    ActiveUserService,
    {
      provide: ActiveUserGrpcProviderService,
      useClass: ActiveUserGrpcMockProviderService,
    },
    {
      provide: GrpcService,
      useValue: {},
    },
  ],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActiveUserPageComponent {
  private readonly toastService = inject(ToastService);
  private readonly activeUserService = inject(ActiveUserService);

  public readonly dailyActiveUser: ReturnType<typeof this.getActiveUser>;
  public readonly monthlyActiveUser: ReturnType<typeof this.getActiveUser>;

  public readonly dailyStart = signal<Date>(new Date(0));
  public readonly dailyEnd = signal<Date>(new Date(0));
  public readonly monthlyStart = signal<Date>(new Date(0));
  public readonly monthlyEnd = signal<Date>(new Date(0));

  constructor() {
    this.dailyActiveUser = this.getActiveUser(
      this.dailyStart,
      this.dailyEnd,
      "dailyPrecision"
    );
    this.monthlyActiveUser = this.getActiveUser(
      this.monthlyStart,
      this.monthlyEnd,
      "monthlyPrecision"
    );

    effect(() => {
      const dailyError = this.dailyActiveUser.error();
      const monthlyError = this.monthlyActiveUser.error();

      if (dailyError) {
        this.toastService.showError(dailyError);
      }
      if (monthlyError) {
        this.toastService.showError(monthlyError);
      }
    });
  }

  private getActiveUser(
    startSignal: Signal<Date>,
    endSignal: Signal<Date>,
    precision: PrecisionType
  ) {
    const request = computed(() => {
      const start = startSignal();
      const end = endSignal();
      if (!start?.getTime() || !end.getTime()) {
        return undefined;
      }

      return {
        precision,
        startingDateInclusive: start,
        endingDateInclusive: end,
      } as const;
    });

    const pipe = new TimestampToDatePipe();
    const foregroundColorSignal = this.getColor("--mat-app-text-color");
    const tooltipBackgroundColorSignal = this.getColor(
      "--mat-dialog-container-color"
    );

    return injectQuery(() => {
      const foregroundColor = foregroundColorSignal();
      const tooltipBackgroundColor = tooltipBackgroundColorSignal();
      return {
        ...this.activeUserService.getActiveUser(request()),
        select: (entries): AgChartOptions => ({
          background: {
            visible: false,
          },
          theme: {
            params: {
              foregroundColor,
              tooltipBackgroundColor,
            },
            palette: {
              fills: precision === "dailyPrecision" ? ["#A194FF"] : ["#347F4A"],
            },
          },
          data: entries.map((entry) => ({
            date: pipe.transform(entry.date),
            value: Number(entry.value),
          })),
          series: [
            {
              type: "line",
              xKey: "date",
              yKey: "value",
              yName: "Users",
              interpolation: { type: "smooth" },
              label: {
                enabled: true,
              },
            },
          ],
        }),
        placeholderData: (prev) => prev,
        staleTime: 5 * 60 * 1000, // 5 minutes
      };
    });
  }

  private getColor(variable: string) {
    const element = inject(ElementRef<HTMLElement>);
    const themeService = inject(ThemeService);

    const compute = () =>
      window.getComputedStyle(element.nativeElement).getPropertyValue(variable);

    const foreGroundColor$ = themeService.isDarkTheme.pipe(
      delay(0),
      map(compute)
    );

    return toSignal(foreGroundColor$, { initialValue: compute() });
  }
}
