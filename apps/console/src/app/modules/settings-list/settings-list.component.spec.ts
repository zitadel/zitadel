import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SettingsListComponent } from './settings-list.component';

describe('SettingsListComponent', () => {
  let component: SettingsListComponent;
  let fixture: ComponentFixture<SettingsListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SettingsListComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SettingsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
