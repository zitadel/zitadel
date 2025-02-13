import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FilterGroupGrantsComponent } from './filter-group-grants.component';

describe('FilterGroupComponent', () => {
  let component: FilterGroupGrantsComponent;
  let fixture: ComponentFixture<FilterGroupGrantsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [FilterGroupGrantsComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FilterGroupGrantsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
