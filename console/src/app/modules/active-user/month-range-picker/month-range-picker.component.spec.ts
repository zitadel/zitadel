import { beforeEach, describe, expect } from "vitest";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { TranslateService } from "@ngx-translate/core";
import { NEVER } from "rxjs";
import {
  MonthRangePickerComponent,
  MonthRangePickerTriggerDirective,
} from "./month-range-picker.component";
import { Component, model, viewChild } from "@angular/core";
import { By } from "@angular/platform-browser";
import { MonthSelectorComponent } from "@/modules/active-user/month-range-picker/month-selector.component";

@Component({
  template: `
    <button cnslMonthRangePickerTrigger>
      Open
      <cnsl-month-range-picker [(start)]="start" [(end)]="end" />
    </button>
  `,
  imports: [MonthRangePickerComponent, MonthRangePickerTriggerDirective],
})
class TestComponent {
  public readonly start = model.required<Date>();
  public readonly end = model.required<Date>();
  public readonly picker = viewChild(MonthRangePickerComponent);
}

describe("MonthRangePickerComponent", () => {
  let component: ComponentFixture<TestComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TestComponent],
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

    component = TestBed.createComponent(TestComponent);
  });

  it("should work", async () => {
    const start = new Date();
    start.setMonth(start.getMonth() - 8);

    const end = new Date();
    end.setMonth(end.getMonth() - 2);

    component.componentRef.setInput("start", start);
    component.componentRef.setInput("end", end);
    component.detectChanges();

    const picker = component.debugElement.query(
      By.css("cnsl-month-range-picker")
    );
    const pickerInstance = toBeInstanceOf(
      picker.componentInstance,
      MonthRangePickerComponent
    );
    pickerInstance.toggle();
    component.detectChanges();

    const startMonthSelector = picker.query(By.css("cnsl-month-selector"));
    const startMonthSelectorInstance = toBeInstanceOf(
      startMonthSelector.componentInstance,
      MonthSelectorComponent
    );
    startMonthSelectorInstance.selectedMonthChange.emit(3);
    component.detectChanges();

    const endMonthSelector = picker.query(By.css("cnsl-month-selector"));
    const endMonthSelectorInstance = toBeInstanceOf(
      endMonthSelector.componentInstance,
      MonthSelectorComponent
    );
    endMonthSelectorInstance.selectedMonthChange.emit(8);
    component.detectChanges();

    expect(component.componentInstance.start().toLocaleDateString()).toBe(
      "4/1/2025"
    );
    expect(component.componentInstance.end().toLocaleDateString()).toBe(
      "9/30/2025"
    );
  });
});

function toBeInstanceOf<T extends object>(
  obj: unknown,
  constructor: new (...params: unknown[]) => T
): T {
  expect(obj).toBeInstanceOf(constructor);
  return obj as T;
}
