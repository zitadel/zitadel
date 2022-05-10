import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SettingsGridComponent } from './settings-grid.component';

describe('SettingsGridComponent', () => {
  let component: SettingsGridComponent;
  let fixture: ComponentFixture<SettingsGridComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SettingsGridComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SettingsGridComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
