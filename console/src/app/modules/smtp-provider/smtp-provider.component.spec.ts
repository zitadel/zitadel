import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SMTPProviderComponent } from './smtp-provider.component';

describe('SMTPProviderSendgridComponent', () => {
  let component: SMTPProviderComponent;
  let fixture: ComponentFixture<SMTPProviderComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SMTPProviderComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SMTPProviderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
