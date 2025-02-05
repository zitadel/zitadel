import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AddActionDialogComponent } from './add-action-dialog.component';

describe('AddKeyDialogComponent', () => {
  let component: AddActionDialogComponent;
  let fixture: ComponentFixture<AddActionDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [AddActionDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddActionDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
