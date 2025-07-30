import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ActionKeysComponent } from './action-keys.component';

describe('ActionKeysComponent', () => {
  let component: ActionKeysComponent;
  let fixture: ComponentFixture<ActionKeysComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ActionKeysComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ActionKeysComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
