import { animate, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, ValidatorFn, Validators } from '@angular/forms';
import { MatLegacySlideToggleChange as MatSlideToggleChange } from '@angular/material/legacy-slide-toggle';
import { Router } from '@angular/router';
import { take } from 'rxjs/operators';
import { ProjectAutocompleteType } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.component';
import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from 'src/app/pages/validators';
import { SetUpOrgRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Project } from 'src/app/proto/generated/zitadel/project_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-app-create',
  templateUrl: './app-create.component.html',
  styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent {
  public projectId: string = '';
  public ProjectAutocompleteType: any = ProjectAutocompleteType;

  constructor(private router: Router, breadcrumbService: BreadcrumbService) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }

  public goToAppCreatePage(): void {
    this.router.navigate(['/projects', this.projectId, 'apps', 'create']);
  }

  public close(): void {
    window.history.back();
  }

  public selectProject(project: Project.AsObject): void {
    if (project.id) {
      this.projectId = project.id;
    }
  }
}
