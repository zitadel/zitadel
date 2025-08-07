import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderLDAPComponent } from './provider-ldap.component';

describe('ProviderLDAPComponent', () => {
  let component: ProviderLDAPComponent;
  let fixture: ComponentFixture<ProviderLDAPComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderLDAPComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderLDAPComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
