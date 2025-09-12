import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FilterUserGrantsComponent } from './filter-user-grants.component';

describe('FilterUserComponent', () => {
  let component: FilterUserGrantsComponent;
  let fixture: ComponentFixture<FilterUserGrantsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [FilterUserGrantsComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FilterUserGrantsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
