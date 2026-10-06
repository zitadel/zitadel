import { signal } from '@angular/core';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { of } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { LoginPolicyComponent } from './login-policy.component';
import { LoginPolicyModule } from './login-policy.module';

describe('LoginPolicyComponent', () => {
  let component: LoginPolicyComponent;
  let fixture: ComponentFixture<LoginPolicyComponent>;

  const lifetimeControls = [
    'passwordCheckLifetime',
    'externalLoginCheckLifetime',
    'mfaInitSkipLifetime',
    'secondFactorCheckLifetime',
    'multiFactorCheckLifetime',
  ];

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [LoginPolicyModule, NoopAnimationsModule, TranslateModule.forRoot()],
      providers: [
        provideRouter([]),
        {
          provide: GrpcAuthService,
          useValue: { isAllowed: () => of(true), hasRoles: () => true, hasRolesApi: () => of(true) },
        },
        { provide: NewAuthService, useValue: { listMyZitadelPermissionsQuery: () => ({ data: signal([]) }) } },
        { provide: ToastService, useValue: { showError: () => {}, showInfo: () => {} } },
        {
          provide: ManagementService,
          useValue: {
            getLoginPolicy: () =>
              Promise.resolve({
                policy: {
                  isDefault: false,
                  secondFactorsList: [],
                  multiFactorsList: [],
                  idpsList: [],
                  passwordCheckLifetime: { seconds: 0, nanos: 0 },
                  externalLoginCheckLifetime: { seconds: 0, nanos: 0 },
                  mfaInitSkipLifetime: { seconds: 0, nanos: 0 },
                  secondFactorCheckLifetime: { seconds: 0, nanos: 0 },
                  multiFactorCheckLifetime: { seconds: 0, nanos: 0 },
                },
              }),
          },
        },
      ],
    }).compileComponents();
  }));

  beforeEach(async () => {
    fixture = TestBed.createComponent(LoginPolicyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('accepts a policy whose lifetimes are all 0', () => {
    lifetimeControls.forEach((name) => {
      expect(component.lifetimeForm.get(name)?.value).toBe(0, `${name} should be loaded as 0`);
      expect(component.lifetimeForm.get(name)?.errors).toBeNull(`${name} should be valid at 0`);
    });
    expect(component.lifetimeForm.valid).toBe(true);
  });

  it('does not render the lifetime inputs with a minimum above 0', () => {
    lifetimeControls.forEach((name) => {
      const input: HTMLInputElement = fixture.nativeElement.querySelector(`input[name="${name}"]`);
      expect(input).withContext(name).not.toBeNull();
      expect(input.min).withContext(`${name} min attribute`).toBe('0');
    });
  });

  it('still rejects a negative lifetime', () => {
    lifetimeControls.forEach((name) => {
      const control = component.lifetimeForm.get(name)!;
      control.setValue(-1);
      expect(control.valid).withContext(name).toBe(false);
    });
  });
});
