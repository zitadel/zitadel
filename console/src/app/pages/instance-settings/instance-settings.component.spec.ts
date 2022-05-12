import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InstanceSettingsComponent } from './instance-settings.component';

describe('InstanceSettingsComponent', () => {
  let component: InstanceSettingsComponent;
  let fixture: ComponentFixture<InstanceSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ InstanceSettingsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(InstanceSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
