import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OrgSettingsComponent } from './org-settings.component';

describe('OrgSettingsComponent', () => {
  let component: OrgSettingsComponent;
  let fixture: ComponentFixture<OrgSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OrgSettingsComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(OrgSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
