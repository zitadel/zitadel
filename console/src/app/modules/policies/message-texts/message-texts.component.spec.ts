import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MessageTextsComponent } from './message-texts.component';

describe('LoginPolicyComponent', () => {
  let component: MessageTextsComponent;
  let fixture: ComponentFixture<MessageTextsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [MessageTextsComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageTextsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
