import { beforeEach, describe, expect } from "vitest";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { firstValueFrom } from "rxjs";
import { YearSelectorComponent } from "@/modules/active-user/month-range-picker/year-selector.component";

describe("YearSelectorComponent", () => {
  let component: ComponentFixture<YearSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [YearSelectorComponent],
    }).compileComponents();

    component = TestBed.createComponent(YearSelectorComponent);
  });

  it("should work", async () => {
    component.componentRef.setInput("selectedYear", 2024);
    component.detectChanges();

    const el = component.debugElement;
    const buttons = el.queryAll(By.css("button"));

    expect(buttons).toHaveLength(8);
    expect(
      (buttons[0].nativeElement as HTMLButtonElement).textContent
    ).toContain("2018");
    expect(
      (buttons[1].nativeElement as HTMLButtonElement).textContent
    ).toContain("2019");
    expect(
      (buttons[2].nativeElement as HTMLButtonElement).textContent
    ).toContain("2020");
    expect(
      (buttons[3].nativeElement as HTMLButtonElement).textContent
    ).toContain("2021");
    expect(
      (buttons[4].nativeElement as HTMLButtonElement).textContent
    ).toContain("2022");
    expect(
      (buttons[5].nativeElement as HTMLButtonElement).textContent
    ).toContain("2023");
    expect(
      (buttons[6].nativeElement as HTMLButtonElement).textContent
    ).toContain("2024");
    expect(buttons[6].classes).toHaveProperty(
      "bg-(--mat-datepicker-calendar-date-selected-state-background-color)"
    );
    expect(
      (buttons[7].nativeElement as HTMLButtonElement).textContent
    ).toContain("2025");

    const selectedYear = firstValueFrom(
      component.componentInstance.selectedYearChange
    );

    (buttons[2].nativeElement as HTMLButtonElement).click();

    await expect(selectedYear).resolves.toBe(2020);
  });
});
