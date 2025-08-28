import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LDAPAttributesComponent } from './ldap-attributes.component';

describe('LDAPAttributesComponent', () => {
  let component: LDAPAttributesComponent;
  let fixture: ComponentFixture<LDAPAttributesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [LDAPAttributesComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(LDAPAttributesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
