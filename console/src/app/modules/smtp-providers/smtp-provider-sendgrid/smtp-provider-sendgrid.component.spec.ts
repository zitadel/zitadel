import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SMTPProviderSendgridComponent } from './smtp-provider-sendgrid.component';

describe('SMTPProviderSendgridComponent', () => {
  let component: SMTPProviderSendgridComponent;
  let fixture: ComponentFixture<SMTPProviderSendgridComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SMTPProviderSendgridComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SMTPProviderSendgridComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
