import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ShowTokenDialogComponent } from './show-token-dialog.component';

describe('ShowKeyDialogComponent', () => {
  let component: ShowTokenDialogComponent;
  let fixture: ComponentFixture<ShowTokenDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ShowTokenDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShowTokenDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
