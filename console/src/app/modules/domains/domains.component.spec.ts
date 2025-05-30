import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DomainsComponent } from './domains.component';

describe('OrgDomainsComponent', () => {
  let component: DomainsComponent;
  let fixture: ComponentFixture<DomainsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [DomainsComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(DomainsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
