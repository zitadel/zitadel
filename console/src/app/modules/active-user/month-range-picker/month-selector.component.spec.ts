import { beforeEach, describe, expect } from "vitest";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MonthSelectorComponent } from "@/modules/active-user/month-range-picker/month-selector.component";
import { By } from "@angular/platform-browser";
import { TranslateService } from "@ngx-translate/core";
import { firstValueFrom, NEVER } from "rxjs";

describe("MonthSelectorComponent", () => {
  let component: ComponentFixture<MonthSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MonthSelectorComponent],
      providers: [
        {
          provide: TranslateService,
          useValue: {
            getCurrentLang: () => "en-US",
            onLangChange: NEVER,
          },
        },
      ],
    }).compileComponents();

    component = TestBed.createComponent(MonthSelectorComponent);
  });

  it("should work", async () => {
    component.componentRef.setInput("selectedMonth", 3);
    // todo: write not before and not after tests
    component.componentRef.setInput("notBefore", 9);
    component.componentRef.setInput("notAfter", 9);
    component.detectChanges();

    const el = component.debugElement;
    const buttons = el.queryAll(By.css("button"));

    expect(buttons).toHaveLength(12);
    expect(
      (buttons[0].nativeElement as HTMLButtonElement).textContent
    ).toContain("Jan");
    expect(
      (buttons[1].nativeElement as HTMLButtonElement).textContent
    ).toContain("Feb");
    expect(
      (buttons[2].nativeElement as HTMLButtonElement).textContent
    ).toContain("Mar");
    expect(
      (buttons[3].nativeElement as HTMLButtonElement).textContent
    ).toContain("Apr");
    expect(buttons[3].classes).toHaveProperty(
      "bg-(--mat-datepicker-calendar-date-selected-state-background-color)"
    );
    expect(
      (buttons[4].nativeElement as HTMLButtonElement).textContent
    ).toContain("May");
    expect(
      (buttons[5].nativeElement as HTMLButtonElement).textContent
    ).toContain("Jun");
    expect(
      (buttons[6].nativeElement as HTMLButtonElement).textContent
    ).toContain("Jul");
    expect(
      (buttons[7].nativeElement as HTMLButtonElement).textContent
    ).toContain("Aug");
    expect(
      (buttons[8].nativeElement as HTMLButtonElement).textContent
    ).toContain("Sep");
    expect(
      (buttons[9].nativeElement as HTMLButtonElement).textContent
    ).toContain("Oct");
    expect(
      (buttons[10].nativeElement as HTMLButtonElement).textContent
    ).toContain("Nov");
    expect(
      (buttons[11].nativeElement as HTMLButtonElement).textContent
    ).toContain("Dec");

    component.componentRef.setInput("selectedMonth", 4);
    component.detectChanges();

    expect(buttons[3].classes).not.toHaveProperty(
      "bg-(--mat-datepicker-calendar-date-selected-state-background-color)"
    );
    expect(buttons[4].classes).toHaveProperty(
      "bg-(--mat-datepicker-calendar-date-selected-state-background-color)"
    );

    const selectedMonth = firstValueFrom(
      component.componentInstance.selectedMonthChange
    );

    (buttons[9].nativeElement as HTMLButtonElement).click();

    await expect(selectedMonth).resolves.toBe(9);
  });
});
