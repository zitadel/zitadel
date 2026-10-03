import { CommonModule, Location } from '@angular/common';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { of } from 'rxjs';

import { FrameworkChangeComponent } from 'src/app/components/framework-change/framework-change.component';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { NavigationService } from 'src/app/services/navigation.service';
import { ToastService } from 'src/app/services/toast.service';

import { IntegrateAppComponent } from './integrate.component';

function setup(queryParams: Record<string, string>) {
  TestBed.configureTestingModule({
    // FrameworkChangeComponent is imported for real, not stubbed: it is the
    // component that reads ?framework= on this page, so stubbing it would hide
    // exactly the behaviour these tests are about.
    imports: [CommonModule, TranslateModule.forRoot(), FrameworkChangeComponent],
    declarations: [IntegrateAppComponent],
    schemas: [NO_ERRORS_SCHEMA],
    providers: [
      {
        provide: ActivatedRoute,
        useValue: {
          snapshot: { paramMap: new Map([['projectid', 'project-1']]) },
          queryParams: of(queryParams),
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
}

describe('IntegrateAppComponent', () => {
  let component: IntegrateAppComponent;
  let fixture: ComponentFixture<IntegrateAppComponent>;

  function create() {
    fixture = TestBed.createComponent(IntegrateAppComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  /**
   * Reached without a framework in the URL — by deep link, reload after the
   * parameter is dropped, or the "other framework" path.
   *
   * Regression: the framework selector is the only thing that can set a
   * framework on this page. Rendering it behind `@if (framework())` deadlocked
   * the page — it could not mount until the framework it sets was already set,
   * so the user got a Create button that never became useful and no control to
   * choose a framework.
   */
  describe('with no framework in the URL', () => {
    beforeEach(waitForAsync(() => setup({})));
    beforeEach(create);

    it('should create', () => {
      expect(component).toBeTruthy();
    });

    it('still renders the framework selector, so the user can pick one', () => {
      expect(component.framework()).toBeUndefined();

      expect(fixture.nativeElement.querySelector('cnsl-framework-change'))
        .withContext('cnsl-framework-change must render with no framework set, or nothing can ever set one')
        .not.toBeNull();
    });
  });

  /**
   * The home page entry point: a quickstart card sends the user to app-create,
   * which forwards ?framework= here. Nothing else on the page picks the
   * framework up, so if this breaks the user reaches Create with an empty
   * app request.
   */
  describe('entered from the home page quickstart (?framework=next)', () => {
    beforeEach(waitForAsync(() => setup({ framework: 'next' })));
    beforeEach(create);

    it('reads the framework from the URL', () => {
      expect(component.framework()?.id).toBe('next');
    });

    it('builds the app request from that framework', () => {
      const request = component.oidcAppRequest.getValue().toObject();

      expect(request.projectId).toBe('project-1');
      expect(request.name).toBe('Next.js');
      expect(request.redirectUrisList.length).withContext('no redirect URIs were preconfigured').toBeGreaterThan(0);
    });
  });

  /**
   * Deep link / stale URL: the pickers no longer offer frameworks without an
   * OIDC configuration, but ?framework= can still name one (Java, Dart /
   * Flutter, or any SDK entry). The page must fall back to an empty request
   * rather than throw.
   */
  describe('with a framework that has no OIDC configuration (?framework=java)', () => {
    beforeEach(waitForAsync(() => setup({ framework: 'java' })));

    it('does not throw while building the app request', () => {
      expect(() => create()).not.toThrow();
    });
  });

  /**
   * Regression: OIDC_CONFIGURATIONS used to hold one shared AddOIDCAppRequest
   * per framework at module scope, and the page configured it in place, so
   * edits leaked into the next app created with the same framework. The
   * entries are factories now; this pins that down.
   */
  describe('creating two apps with the same framework', () => {
    beforeEach(waitForAsync(() => setup({ framework: 'next' })));
    beforeEach(create);

    it('does not carry redirect URIs over from the previous app', () => {
      const pristine = [...component.redirectUris];
      component.redirectUris = [...pristine, 'http://localhost:3000/leaked'];

      // A second visit to the page, as after navigating away and starting again.
      const second = TestBed.createComponent(IntegrateAppComponent);
      second.detectChanges();

      expect(second.componentInstance.redirectUris)
        .withContext('redirect URIs leaked between app creations')
        .toEqual(pristine);
    });
  });
});
