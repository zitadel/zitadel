import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FilterEventsAggregateComponent } from './filter-events-aggregate.component';

describe('FilterEventsAggregateComponent', () => {
  let component: FilterEventsAggregateComponent;
  let fixture: ComponentFixture<FilterEventsAggregateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [FilterEventsAggregateComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FilterEventsAggregateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
