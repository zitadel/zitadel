import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SMTPProviderGoogleComponent } from './smtp-provider-google.component';

describe('SMTPProviderGoogleComponent', () => {
  let component: SMTPProviderGoogleComponent;
  let fixture: ComponentFixture<SMTPProviderGoogleComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SMTPProviderGoogleComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SMTPProviderGoogleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
