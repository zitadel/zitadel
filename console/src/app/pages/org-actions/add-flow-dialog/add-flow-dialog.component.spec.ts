import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AddFlowDialogComponent } from './add-flow-dialog.component';

describe('AddKeyDialogComponent', () => {
  let component: AddFlowDialogComponent;
  let fixture: ComponentFixture<AddFlowDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [AddFlowDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddFlowDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
