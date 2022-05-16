import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OIDCConfigurationComponent } from './oidc-configuration.component';

describe('OIDCConfigurationComponent', () => {
  let component: OIDCConfigurationComponent;
  let fixture: ComponentFixture<OIDCConfigurationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OIDCConfigurationComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(OIDCConfigurationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
