import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OrgDomainsComponent } from './org-domains.component';

describe('OrgDomainsComponent', () => {
  let component: OrgDomainsComponent;
  let fixture: ComponentFixture<OrgDomainsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OrgDomainsComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(OrgDomainsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
