import { Location } from '@angular/common';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';
import { UpdateGenericOAuthProviderRequest as MgmtUpdateGenericOAuthProviderRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { AutoLinkingOption, Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';

import { TranslateModule } from '@ngx-translate/core';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import ProvidersModule from '../providers.module';
import { ProviderNextService } from '../provider-next/provider-next.service';
import { ProviderOAuthComponent } from './provider-oauth.component';

describe('ProviderOAuthComponent', () => {
  let component: ProviderOAuthComponent;
  let fixture: ComponentFixture<ProviderOAuthComponent>;
  let capturedRequest: MgmtUpdateGenericOAuthProviderRequest | undefined;

  const storedOptions: Options.AsObject = {
    isAutoCreation: true,
    isAutoUpdate: true,
    isCreationAllowed: false,
    isLinkingAllowed: false,
    autoLinking: AutoLinkingOption.AUTO_LINKING_OPTION_USERNAME,
  };

  const provider: Provider.AsObject = {
    id: 'id1',
    name: 'Test Provider',
    config: {
      oauth: {
        clientId: 'client',
        authorizationEndpoint: 'https://example.com/auth',
        tokenEndpoint: 'https://example.com/token',
        userEndpoint: 'https://example.com/user',
        idAttribute: 'sub',
        scopesList: ['openid'],
        usePkce: false,
      },
      options: storedOptions,
    },
  } as unknown as Provider.AsObject;

  const managementStub = {
    getProviderByID: () => Promise.resolve({ idp: provider }),
    updateGenericOAuthProvider: (req: MgmtUpdateGenericOAuthProviderRequest) => {
      capturedRequest = req;
      return Promise.resolve({});
    },
  };

  const routeStub = {
    data: of({ serviceType: PolicyComponentServiceType.MGMT }),
    paramMap: of(),
    snapshot: { paramMap: { get: (name: string) => (name === 'id' ? 'id1' : null) } },
  };

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [ProvidersModule, TranslateModule.forRoot()],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        { provide: GrpcAuthService, useValue: { isAllowed: () => of(true) } },
        { provide: ManagementService, useValue: managementStub },
        { provide: ActivatedRoute, useValue: routeStub },
        { provide: ToastService, useValue: { showError: () => {} } },
        { provide: BreadcrumbService, useValue: { setBreadcrumb: () => {} } },
        {
          provide: ProviderNextService,
          useValue: {
            service: () => of(managementStub),
            id: () => of('id1'),
            exists: () => of(true),
            autofillLink: () => of(''),
            activateLink: () => of(''),
            expandWhatNow: () => of(false),
            callbackUrls: () => of([]),
          },
        },
        { provide: Location, useValue: { back: () => {} } },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    capturedRequest = undefined;
    fixture = TestBed.createComponent(ProviderOAuthComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('saves the stored provider options when the optional section is never expanded', async () => {
    await fixture.whenStable();
    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(component.provider).toBeDefined();
    expect(component.showOptional).toBe(false);

    component.submitForm();

    expect(capturedRequest).toBeDefined();
    const sent = capturedRequest?.getProviderOptions();
    expect(sent?.getIsAutoCreation()).toBe(true);
    expect(sent?.getIsAutoUpdate()).toBe(true);
    expect(sent?.getIsCreationAllowed()).toBe(false);
    expect(sent?.getIsLinkingAllowed()).toBe(false);
    expect(sent?.getAutoLinking()).toBe(AutoLinkingOption.AUTO_LINKING_OPTION_USERNAME);
  });
});
