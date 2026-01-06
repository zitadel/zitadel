import {
  provideTanStackQuery,
  QueryClient,
} from "@tanstack/angular-query-experimental";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActiveUserPageComponent } from "@/modules/active-user/active-user-page.component";
import { TranslateService } from "@ngx-translate/core";
import { beforeEach, expect } from "vitest";

describe("ActiveUserPageComponent", () => {
  // Arrange
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false, // ✅ faster failure tests
      },
    },
  });
  let component: ComponentFixture<ActiveUserPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ActiveUserPageComponent],
      providers: [
        provideTanStackQuery(queryClient),
        { provide: TranslateService, useValue: {} },
      ],
    }).compileComponents();

    component = TestBed.createComponent(ActiveUserPageComponent);
  });

  test("should fetch daily active users", async () => {
    vi.useFakeTimers();

    const start = new Date();
    start.setDate(start.getDate() - 10);
    const end = new Date();

    component.componentInstance.dailyStart.set(start);
    component.componentInstance.dailyEnd.set(end);

    component.detectChanges();
    await vi.advanceTimersByTimeAsync(0);
    await Promise.resolve();
    await component.whenStable();

    expect(component.componentInstance.dailyActiveUser.data()).toBeDefined();
    expect(component.componentInstance.dailyActiveUser.error()).toBeNull();
  });

  test("should fetch monthly active users", async () => {
    vi.useFakeTimers();

    const start = new Date();
    start.setMonth(start.getMonth() - 10);
    const end = new Date();

    component.componentInstance.monthlyStart.set(start);
    component.componentInstance.monthlyEnd.set(end);

    component.detectChanges();
    await vi.advanceTimersByTimeAsync(0);
    await Promise.resolve();
    await component.whenStable();

    expect(component.componentInstance.monthlyActiveUser.data()).toBeDefined();
    expect(component.componentInstance.monthlyActiveUser.error()).toBeNull();
  });
});
