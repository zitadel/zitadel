import { MonthService } from "./month.service";
import { TestBed } from "@angular/core/testing";
import { ReplaySubject } from "rxjs";
import { TranslateService } from "@ngx-translate/core";

describe("MonthService", () => {
  const translateService = {
    getCurrentLang: () => "en-US",
    onLangChange: new ReplaySubject<{ lang: string }>(1),
  };

  beforeEach(() =>
    TestBed.configureTestingModule({
      providers: [
        MonthService,
        { provide: TranslateService, useValue: translateService },
      ],
    })
  );

  it("#getMonths should return monthNames", async () => {
    const monthService = TestBed.inject(MonthService);

    expect(monthService.monthNames()).toEqual([
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ]);

    translateService.onLangChange.next({ lang: "de-DE" });

    expect(monthService.monthNames()).toEqual([
      "Jan",
      "Feb",
      "Mär",
      "Apr",
      "Mai",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Okt",
      "Nov",
      "Dez",
    ]);
  });
});
