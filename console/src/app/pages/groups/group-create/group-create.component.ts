import { Location } from '@angular/common';
import { ChangeDetectorRef, Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, ValidatorFn, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subject } from 'rxjs';
import { AddGroupRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Domain } from 'src/app/proto/generated/zitadel/org_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { minLengthValidator, requiredValidator } from '../../../modules/form-field/validators/validators';
import { LanguagesService } from '../../../services/languages.service';

@Component({
  selector: 'cnsl-group-create',
  templateUrl: './group-create.component.html',
  styleUrls: ['./group-create.component.scss'],
})
export class GroupCreateComponent implements OnInit, OnDestroy {
  public group: AddGroupRequest.AsObject = new AddGroupRequest().toObject();
  public groupForm!: UntypedFormGroup;
  private destroyed$: Subject<void> = new Subject();

  public loading: boolean = false;

  @ViewChild('suffix') public suffix!: any;
  private primaryDomain!: Domain.AsObject;
  public policy!: PasswordComplexityPolicy.AsObject;

  constructor(
    private router: Router,
    private toast: ToastService,
    private fb: UntypedFormBuilder,
    private mgmtService: ManagementService,
    private changeDetRef: ChangeDetectorRef,
    private _location: Location,
    public langSvc: LanguagesService,
    breadcrumbService: BreadcrumbService,
  ) {
    this.initForm();
    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);

  }

  public close(): void {
    this._location.back();
  }

  private initForm(): void {
    this.groupForm = this.fb.group({
      name: ['', [requiredValidator, minLengthValidator(1)]],
    });

    const validators: Validators[] = [requiredValidator];
  }

  public createGroup(): void {
    this.group = this.groupForm.value;

    this.loading = true;

    const humanReq = new AddGroupRequest();
    humanReq.setName(this.name?.value);
    humanReq.setDescription(this.description?.value);

    this.mgmtService
      .addGroup(humanReq)
      .then((data) => {
        this.loading = false;
        this.toast.showInfo('GROUP.TOAST.CREATED', true);
        this.router.navigate(['groups'], { queryParams: { new: true } });
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
  }

  ngOnInit(): void {}

  ngOnDestroy(): void {
    this.destroyed$.next();
    this.destroyed$.complete();
  }

  public get name(): AbstractControl | null {
    return this.groupForm.get('name');
  }
  public get description(): AbstractControl | null {
    return this.groupForm.get('description');
  }
}
