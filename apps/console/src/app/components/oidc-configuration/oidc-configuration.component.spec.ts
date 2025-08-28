import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OIDCConfigurationComponent } from './oidc-configuration.component';

describe('QuickstartComponent', () => {
  let component: OIDCConfigurationComponent;
  let fixture: ComponentFixture<OIDCConfigurationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [OIDCConfigurationComponent],
    });
    fixture = TestBed.createComponent(OIDCConfigurationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
