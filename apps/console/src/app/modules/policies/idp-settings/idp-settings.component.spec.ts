import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IdpSettingsComponent } from './idp-settings.component';

describe('IdpSettingsComponent', () => {
  let component: IdpSettingsComponent;
  let fixture: ComponentFixture<IdpSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [IdpSettingsComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IdpSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
