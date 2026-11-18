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
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';

import { AppCreateComponent } from './app-create.component';

/** Mirrors what cnsl-search-project-autocomplete emits for an owned project. */
function selectOwnedProject(component: AppCreateComponent, name: string, id = 'project-1') {
  component.selectProject({ project: { id, name } as any, type: ProjectType.PROJECTTYPE_OWNED, name });
  component.projectName = name;
}

function setup(queryParams: Record<string, string>) {
  const router = { navigate: jasmine.createSpy('navigate') };

  TestBed.configureTestingModule({
    // FrameworkChangeComponent is imported for real, not stubbed: it is the
    // component that reads ?framework= on this page, so stubbing it would hide
    // exactly the behaviour these tests are about.
    imports: [CommonModule, TranslateModule.forRoot(), FrameworkChangeComponent],
    declarations: [AppCreateComponent],
    schemas: [NO_ERRORS_SCHEMA],
    providers: [
      { provide: MatDialog, useValue: { open: () => ({ afterClosed: () => of(undefined) }) } },
      { provide: ActivatedRoute, useValue: { queryParams: of(queryParams), snapshot: { paramMap: new Map() } } },
      { provide: Router, useValue: router },
      { provide: ManagementService, useValue: { addProject: () => Promise.resolve({ id: 'new-project' }) } },
      { provide: BreadcrumbService, useValue: { setBreadcrumb: () => {} } },
      { provide: NavigationService, useValue: { isBackPossible: false } },
      { provide: Location, useValue: { back: () => {} } },
    ],
  }).compileComponents();

  return router;
}

describe('AppCreateComponent', () => {
  let component: AppCreateComponent;
  let fixture: ComponentFixture<AppCreateComponent>;

  describe('entered from the projects tab (no ?framework=)', () => {
    beforeEach(waitForAsync(() => {
      setup({});
    }));

    beforeEach(() => {
      fixture = TestBed.createComponent(AppCreateComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();
    });

    it('should create', () => {
      expect(component).toBeTruthy();
    });

    it('renders a framework picker the user can operate', () => {
      expect(fixture.nativeElement.querySelector('cnsl-framework-autocomplete')).not.toBeNull();
    });

    /**
     * frameworks.json also lists SDK/library entries. Offering one here sends
     * the user to an integrate page that has no app configuration to build
     * from, which used to throw.
     */
    it('only offers frameworks that can actually be integrated', () => {
      const notIntegrable = component.frameworks.filter((f) => !f.id || !OIDC_CONFIGURATIONS[f.id]);

      expect(notIntegrable.map((f) => f.title))
        .withContext('framework picker offers frameworks with no OIDC configuration')
        .toEqual([]);
      expect(component.frameworks.length).toBeGreaterThan(0);
    });
  });

  /**
   * The home page quickstart cards link to
   * /projects/app-create?framework=<id>, so this is the entry point every user
   * who starts from the home page lands on.
   */
  describe('entered from the home page quickstart (?framework=angular)', () => {
    let router: { navigate: jasmine.Spy };

    beforeEach(waitForAsync(() => {
      router = setup({ framework: 'angular' });
    }));

    beforeEach(() => {
      fixture = TestBed.createComponent(AppCreateComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();
    });

    it('preselects the framework named in the query parameter', () => {
      expect(component.framework()?.id).toBe('angular');
    });

    it('renders a framework picker the user can operate', () => {
      const picker =
        fixture.nativeElement.querySelector('cnsl-framework-autocomplete') ??
        fixture.nativeElement.querySelector('cnsl-framework-change');

      expect(picker).withContext('no control on the page can change the framework').not.toBeNull();
    });

    it('enables the continue button once a project is picked', () => {
      selectOwnedProject(component, 'My Project');
      fixture.detectChanges();

      const button: HTMLButtonElement = fixture.nativeElement.querySelector('.continue-button');
      expect(button).not.toBeNull();
      expect(button.disabled).withContext('continue button stayed disabled').toBeFalse();
    });

    it('continues to the integrate page for the selected project', () => {
      selectOwnedProject(component, 'My Project');
      fixture.detectChanges();

      component.goToAppIntegratePage();

      expect(router.navigate).toHaveBeenCalledWith(['/projects', 'project-1', 'apps', 'integrate'], {
        queryParams: { framework: 'angular' },
      });
    });
  });
});
