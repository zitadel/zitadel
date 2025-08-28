import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProviderSamlSpComponent } from './provider-saml-sp.component';

describe('ProviderSamlSpComponent', () => {
  let component: ProviderSamlSpComponent;
  let fixture: ComponentFixture<ProviderSamlSpComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderSamlSpComponent],
    });
    fixture = TestBed.createComponent(ProviderSamlSpComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
