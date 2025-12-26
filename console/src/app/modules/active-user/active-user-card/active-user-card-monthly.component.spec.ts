import { ComponentFixture, TestBed } from "@angular/core/testing";
import { beforeEach } from "vitest";
import { ActiveUserCardMonthlyComponent } from "./active-user-card-monthly.component";

describe("ActiveUserCardMonthlyComponent", () => {
  let component: ComponentFixture<ActiveUserCardMonthlyComponent>;
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ActiveUserCardMonthlyComponent],
    }).compileComponents();

    component = TestBed.createComponent(ActiveUserCardMonthlyComponent);
  });

  test("should change to custom", async () => {
    component.componentInstance.start.set(new Date(0));
    component.componentInstance.end.set(new Date(0));
    component.componentInstance.state.set("custom");
    component.detectChanges();
  });
});
