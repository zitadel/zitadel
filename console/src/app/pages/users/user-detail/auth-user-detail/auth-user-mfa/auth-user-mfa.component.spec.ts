import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { of } from 'rxjs';
import { AuthUserMfaComponent } from './auth-user-mfa.component';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatDialog } from '@angular/material/dialog';
import { AuthFactor, AuthFactorState } from 'src/app/proto/generated/zitadel/user_pb';
import { SecondFactorType } from 'src/app/proto/generated/zitadel/policy_pb';
import { CardComponent } from '../../../../../modules/card/card.component';
import { RefreshTableComponent } from '../../../../../modules/refresh-table/refresh-table.component';
import { TranslateModule } from '@ngx-translate/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('AuthUserMfaComponent', () => {
  let component: AuthUserMfaComponent;
  let fixture: ComponentFixture<AuthUserMfaComponent>;
  let serviceStub: Partial<GrpcAuthService>;
  let toastStub: Partial<ToastService>;
  let dialogStub: Partial<MatDialog>;

  beforeEach(waitForAsync(() => {
    // Create stubs for required services
    serviceStub = {
      listMyMultiFactors: jasmine.createSpy('listMyMultiFactors').and.returnValue(Promise.resolve({
        resultList: [
          { otp: true, state: AuthFactorState.AUTH_FACTOR_STATE_READY } as AuthFactor.AsObject,
          { otpSms: true, state: AuthFactorState.AUTH_FACTOR_STATE_READY } as AuthFactor.AsObject,
          { otpEmail: true, state: AuthFactorState.AUTH_FACTOR_STATE_READY } as AuthFactor.AsObject,
          { state: AuthFactorState.AUTH_FACTOR_STATE_NOT_READY } as AuthFactor.AsObject
        ]
      })),
      getMyLoginPolicy: jasmine.createSpy('getMyLoginPolicy').and.returnValue(Promise.resolve({
        policy: {
          secondFactorsList: [
            SecondFactorType.SECOND_FACTOR_TYPE_OTP,
            SecondFactorType.SECOND_FACTOR_TYPE_U2F,
            SecondFactorType.SECOND_FACTOR_TYPE_OTP_EMAIL,
            SecondFactorType.SECOND_FACTOR_TYPE_OTP_SMS
          ]
        }
      })),
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
        afterClosed: () => of(true)
      })
    };

    TestBed.configureTestingModule({
      declarations: [AuthUserMfaComponent, CardComponent, RefreshTableComponent],
      imports: [MatIconModule, TranslateModule.forRoot(), MatTooltipModule, MatTableModule, BrowserAnimationsModule],
      providers: [
        { provide: GrpcAuthService, useValue: serviceStub },
        { provide: ToastService, useValue: toastStub },
        { provide: MatDialog, useValue: dialogStub },
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AuthUserMfaComponent);
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
    // Our stub returned 4 items
    expect(component.dataSource.data.length).toBe(4);

    // Pipes were updated
    component.otpDisabled$.subscribe(value => {
      expect(value).toBeTrue();
    });
    component.otpSmsDisabled$.subscribe(value => {
      expect(value).toBeTrue();
    });
    component.otpEmailDisabled$.subscribe(value => {
      expect(value).toBeTrue();
    });
  });

  it('should call deleteMFA and remove OTP factor', async () => {
    // OTP is set
    const factor = { otp: true, state: AuthFactorState.AUTH_FACTOR_STATE_READY } as AuthFactor.AsObject;
    await component.deleteMFA(factor);

    // Verify that the service method for OTP removal was called
    expect(serviceStub.removeMyMultiFactorOTP).toHaveBeenCalled();
    expect(serviceStub.listMyMultiFactors).toHaveBeenCalled();
  });
});