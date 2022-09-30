import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KeyboardShortcutsComponent } from './keyboard-shortcuts.component';

describe('KeyboardShortcutsComponent', () => {
  let component: KeyboardShortcutsComponent;
  let fixture: ComponentFixture<KeyboardShortcutsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KeyboardShortcutsComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KeyboardShortcutsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
