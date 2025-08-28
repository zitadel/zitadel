import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { of } from 'rxjs';
import { AuthUserMfaComponent } from './auth-user-mfa.component';
import { ToastService } from 'src/app/services/toast.service';
import { MatDialog } from '@angular/material/dialog';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { SecondFactorType } from 'src/app/proto/generated/zitadel/policy_pb';
import { CardComponent } from 'src/app/modules/card/card.component';
import { RefreshTableComponent } from 'src/app/modules/refresh-table/refresh-table.component';
import { TranslateModule } from '@ngx-translate/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AuthFactor, AuthFactorState } from '@zitadel/proto/zitadel/user_pb';

describe('AuthUserMfaComponent', () => {
  // Create a test host component that extends the original component
  class TestHostComponent extends AuthUserMfaComponent {
    // Expose protected properties for testing
    public getOtpEmailDisabled$() {
      return this.otpEmailDisabled$;
    }

    public getOtpDisabled$() {
      return this.otpDisabled$;
    }

    public getOtpSmsDisabled$() {
      return this.otpSmsDisabled$;
    }
  }

  let component: TestHostComponent;
  let fixture: ComponentFixture<TestHostComponent>;
  let serviceStub: Partial<NewAuthService>;
  let toastStub: Partial<ToastService>;
  let dialogStub: Partial<MatDialog>;

  beforeEach(waitForAsync(() => {
    // Create stubs for required services
    serviceStub = {
      listMyMultiFactors: jasmine.createSpy('listMyMultiFactors').and.returnValue(
        Promise.resolve({
          result: [
            { type: { case: 'otp' }, state: AuthFactorState.READY, $typeName: 'zitadel.user.v1.AuthFactor' } as AuthFactor,
            {
              type: { case: 'otpSms' },
              state: AuthFactorState.READY,
              $typeName: 'zitadel.user.v1.AuthFactor',
            } as AuthFactor,
            {
              type: { case: 'otpEmail' },
              state: AuthFactorState.READY,
              $typeName: 'zitadel.user.v1.AuthFactor',
            } as AuthFactor,
          ],
        }),
      ),
      getMyLoginPolicy: jasmine.createSpy('getMyLoginPolicy').and.returnValue(
        Promise.resolve({
          policy: {
            secondFactorsList: [
              SecondFactorType.SECOND_FACTOR_TYPE_OTP,
              SecondFactorType.SECOND_FACTOR_TYPE_U2F,
              SecondFactorType.SECOND_FACTOR_TYPE_OTP_EMAIL,
              SecondFactorType.SECOND_FACTOR_TYPE_OTP_SMS,
            ],
          },
        }),
      ),
      removeMyMultiFactorOTP: jasmine.createSpy('removeMyMultiFactorOTP').and.returnValue(Promise.resolve()),
      removeMyMultiFactorU2F: jasmine.createSpy('removeMyMultiFactorU2F').and.returnValue(Promise.resolve()),
      removeMyAuthFactorOTPEmail: jasmine.createSpy('removeMyAuthFactorOTPEmail').and.returnValue(Promise.resolve()),
      removeMyAuthFactorOTPSMS: jasmine.createSpy('removeMyAuthFactorOTPSMS').and.returnValue(Promise.resolve()),
    };

    toastStub = {
      showInfo: jasmine.createSpy('showInfo'),
      showError: jasmine.createSpy('showError'),
    };

    dialogStub = {
      // Opened dialog returns a truthy value after closing
      open: jasmine.createSpy('open').and.returnValue({
        afterClosed: () => of(true),
      }),
    };

    TestBed.configureTestingModule({
      declarations: [TestHostComponent, CardComponent, RefreshTableComponent], // Use TestHostComponent instead
      imports: [MatIconModule, TranslateModule.forRoot(), MatTooltipModule, MatTableModule, BrowserAnimationsModule],
      providers: [
        { provide: NewAuthService, useValue: serviceStub },
        { provide: ToastService, useValue: toastStub },
        { provide: MatDialog, useValue: dialogStub },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TestHostComponent); // Use TestHostComponent
    component = fixture.componentInstance;
    // Optionally set the phoneVerified input if needed by your tests
    component.phoneVerified = true;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should call getMFAs and update dataSource and disable flags', async () => {
    // Call the method and wait for the Promise resolution
    await component.getMFAs();
    fixture.detectChanges();

    expect(serviceStub.listMyMultiFactors).toHaveBeenCalled();
    // Our stub returns 3 items
    expect(component.dataSource.data.length).toBe(3);

    // Use the public getter methods to access protected properties
    component.getOtpDisabled$().subscribe((value) => {
      expect(value).toBeTrue();
    });
    component.getOtpSmsDisabled$().subscribe((value) => {
      expect(value).toBeTrue();
    });
    component.getOtpEmailDisabled$().subscribe((value) => {
      expect(value).toBeTrue();
    });
  });

  it('should call deleteMFA and remove OTP factor', async () => {
    // OTP is set
    const factor = {
      type: { case: 'otp' },
      state: AuthFactorState.READY,
      $typeName: 'zitadel.user.v1.AuthFactor',
    } as AuthFactor;
    await component.deleteMFA(factor);

    // Verify that the service method for OTP removal was called
    expect(serviceStub.removeMyMultiFactorOTP).toHaveBeenCalled();
    expect(serviceStub.listMyMultiFactors).toHaveBeenCalled();
  });
});
