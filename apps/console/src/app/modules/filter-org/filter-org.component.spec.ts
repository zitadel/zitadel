import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FilterOrgComponent } from './filter-org.component';

describe('FilterOrgComponent', () => {
  let component: FilterOrgComponent;
  let fixture: ComponentFixture<FilterOrgComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [FilterOrgComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FilterOrgComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
