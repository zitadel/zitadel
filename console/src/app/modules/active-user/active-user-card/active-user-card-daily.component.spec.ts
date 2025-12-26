import { ComponentFixture, TestBed } from "@angular/core/testing";
import { beforeEach, expect } from "vitest";
import { ActiveUserCardDailyComponent } from "./active-user-card-daily.component";

describe("ActiveUserCardDailyComponent", () => {
  let component: ComponentFixture<ActiveUserCardDailyComponent>;
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ActiveUserCardDailyComponent],
    }).compileComponents();

    component = TestBed.createComponent(ActiveUserCardDailyComponent);
  });

  test("form changes should change outputs", async () => {
    component.componentInstance.start.set(new Date(0));
    component.componentInstance.end.set(new Date(0));
    component.componentInstance.state.set("custom");

    component.detectChanges();

    const start = new Date();
    start.setMonth(start.getMonth() - 11);

    const end = new Date();

    component.componentInstance.form().setValue({
      start,
      end,
    });

    component.detectChanges();

    expect(component.componentInstance.start()).toBe(start);
    expect(component.componentInstance.end()).toBe(end);
  });
});
