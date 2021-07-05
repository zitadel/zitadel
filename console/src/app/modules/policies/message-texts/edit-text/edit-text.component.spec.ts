import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditTextComponent } from './edit-text.component';

describe('EditTextComponent', () => {
  let component: EditTextComponent;
  let fixture: ComponentFixture<EditTextComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EditTextComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(EditTextComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
