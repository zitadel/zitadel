import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { Router, ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { MatDialog } from '@angular/material/dialog';
import { of } from 'rxjs';

import { AppQuickCreateComponent } from './app-quick-create.component';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { NavigationService } from 'src/app/services/navigation.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ThemeService } from 'src/app/services/theme.service';

describe('AppQuickCreateComponent', () => {
  let component: AppQuickCreateComponent;
  let fixture: ComponentFixture<AppQuickCreateComponent>;

  beforeEach(waitForAsync(() => {
    const mockRouter = {
      navigate: jasmine.createSpy('navigate'),
    };

    const mockActivatedRoute = {
      queryParams: of({}),
    };

    const mockManagementService = jasmine.createSpyObj('ManagementService', ['addProject', 'addOIDCApp']);
    const mockToastService = jasmine.createSpyObj('ToastService', ['showInfo', 'showError']);
    const mockBreadcrumbService = jasmine.createSpyObj('BreadcrumbService', ['setBreadcrumb']);
    const mockLocation = jasmine.createSpyObj('Location', ['back']);
    const mockNavigationService = {
      isBackPossible: false,
    };
    const mockAuthService = jasmine.createSpyObj('GrpcAuthService', ['getActiveOrg']);
    const mockThemeService = {
      isDarkTheme: of(false),
    };
    const mockDialog = jasmine.createSpyObj('MatDialog', ['open']);

    TestBed.configureTestingModule({
      declarations: [AppQuickCreateComponent],
      providers: [
        { provide: Router, useValue: mockRouter },
        { provide: ActivatedRoute, useValue: mockActivatedRoute },
        { provide: ManagementService, useValue: mockManagementService },
        { provide: ToastService, useValue: mockToastService },
        { provide: BreadcrumbService, useValue: mockBreadcrumbService },
        { provide: Location, useValue: mockLocation },
        { provide: NavigationService, useValue: mockNavigationService },
        { provide: GrpcAuthService, useValue: mockAuthService },
        { provide: ThemeService, useValue: mockThemeService },
        { provide: MatDialog, useValue: mockDialog },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppQuickCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
