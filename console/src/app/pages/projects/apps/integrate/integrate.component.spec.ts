import { CommonModule } from '@angular/common';
import { Location } from '@angular/common';
import { NO_ERRORS_SCHEMA, Pipe, PipeTransform } from '@angular/core';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { of } from 'rxjs';

import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { NavigationService } from 'src/app/services/navigation.service';
import { ToastService } from 'src/app/services/toast.service';

import { IntegrateAppComponent } from './integrate.component';

@Pipe({ name: 'translate', standalone: false })
class TranslateStubPipe implements PipeTransform {
  public transform(value: string): string {
    return value;
  }
}

describe('IntegrateAppComponent', () => {
  let component: IntegrateAppComponent;
  let fixture: ComponentFixture<IntegrateAppComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [CommonModule],
      declarations: [IntegrateAppComponent, TranslateStubPipe],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: { paramMap: new Map([['projectid', 'project-1']]) },
            queryParams: of({ framework: 'next' }),
          },
        },
        { provide: Router, useValue: { navigate: jasmine.createSpy('navigate') } },
        { provide: ToastService, useValue: { showError: () => {}, showInfo: () => {} } },
        { provide: MatDialog, useValue: { open: () => ({ afterClosed: () => of(undefined) }) } },
        { provide: ManagementService, useValue: { ownedProjects: of([]), grantedProjects: of([]) } },
        { provide: Location, useValue: { back: () => {} } },
        { provide: BreadcrumbService, useValue: { setBreadcrumb: () => {} } },
        { provide: NavigationService, useValue: { isBackPossible: false } },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IntegrateAppComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  /**
   * Regression: the framework selector is the only thing that reads the
   * ?framework= query parameter and calls setFramework(). Rendering it behind
   * `@if (framework())` deadlocked the page — it could not mount until the
   * framework it sets was already set, so a user arriving from the quickstart
   * (app-create navigates here with queryParams: { framework }) got a Create
   * button disabled forever and no control to choose a framework.
   */
  it('renders the framework selector even when no framework is selected yet', () => {
    expect(component.framework()).toBeUndefined();

    const selector = fixture.nativeElement.querySelector('cnsl-framework-change');

    expect(selector)
      .withContext('cnsl-framework-change must render with no framework set, or nothing can ever set one')
      .not.toBeNull();
  });
});
