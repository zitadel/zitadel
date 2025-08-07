import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ActionTableComponent } from './action-table.component';

describe('ActionTableComponent', () => {
  let component: ActionTableComponent;
  let fixture: ComponentFixture<ActionTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ActionTableComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ActionTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
