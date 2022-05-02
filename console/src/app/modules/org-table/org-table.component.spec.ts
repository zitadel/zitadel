import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OrgTableComponent } from './org-table.component';

describe('OrgsComponent', () => {
  let component: OrgTableComponent;
  let fixture: ComponentFixture<OrgTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OrgTableComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(OrgTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
