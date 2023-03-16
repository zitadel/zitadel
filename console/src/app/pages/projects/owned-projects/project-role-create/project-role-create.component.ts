import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { UntypedFormArray, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { BulkAddProjectRolesRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-project-role-create',
  templateUrl: './project-role-create.component.html',
  styleUrls: ['./project-role-create.component.scss'],
  animations: [
    trigger('list', [transition(':enter', [query('@animate', stagger(80, animateChild()))])]),
    trigger('animate', [
      transition(':enter', [
        style({ opacity: 0, transform: 'translateY(-100%)' }),
        animate('100ms', style({ opacity: 1, transform: 'translateY(0)' })),
      ]),
      transition(':leave', [
        style({ opacity: 1, transform: 'translateY(0)' }),
        animate('100ms', style({ opacity: 0, transform: 'translateY(100%)' })),
      ]),
    ]),
  ],
})
export class ProjectRoleCreateComponent implements OnInit, OnDestroy {
  private subscription: Subscription = new Subscription();
  public projectId: string = '';

  public formArray!: UntypedFormArray;
  public formGroup: UntypedFormGroup = this.fb.group({
    key: new UntypedFormControl('', [requiredValidator]),
    displayName: new UntypedFormControl(''),
    group: new UntypedFormControl(''),
  });

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private toast: ToastService,
    private mgmtService: ManagementService,
    private fb: UntypedFormBuilder,
    private _location: Location,
  ) {
    this.formArray = new UntypedFormArray([this.formGroup]);
  }

  public addEntry(): void {
    const newGroup = this.fb.group({
      key: new UntypedFormControl('', [requiredValidator]),
      displayName: new UntypedFormControl(''),
      group: new UntypedFormControl(''),
    });

    this.formArray.push(newGroup);
  }

  public removeEntry(index: number): void {
    this.formArray.removeAt(index);
  }

  public ngOnInit(): void {
    this.subscription = this.route.params.subscribe((params) => this.getData(params));
  }

  public ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  private getData({ projectid }: Params): void {
    this.projectId = projectid;
  }

  public addRole(): void {
    const rolesToAdd: BulkAddProjectRolesRequest.Role[] = this.formArray.value.map((element: any) => {
      const role = new BulkAddProjectRolesRequest.Role();
      role.setKey(element.key);
      role.setDisplayName(element.displayName);
      role.setGroup(element.group);
      return role;
    });

    this.mgmtService
      .bulkAddProjectRoles(this.projectId, rolesToAdd)
      .then(() => {
        this.toast.showInfo('PROJECT.TOAST.ROLESCREATED', true);
        this.router.navigate(['projects', this.projectId], { queryParams: { id: 'roles' } });
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public close(): void {
    this._location.back();
  }
}
