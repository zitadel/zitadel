import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddTokenDialogComponent } from './add-token-dialog.component';

describe('AddTokenDialogComponent', () => {
  let component: AddTokenDialogComponent;
  let fixture: ComponentFixture<AddTokenDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AddTokenDialogComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AddTokenDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
