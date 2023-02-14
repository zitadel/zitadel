import { Location } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { AddMachineUserRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { AccessTokenType } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-user-create-machine',
  templateUrl: './user-create-machine.component.html',
  styleUrls: ['./user-create-machine.component.scss'],
})
export class UserCreateMachineComponent implements OnDestroy {
  public user: AddMachineUserRequest.AsObject = new AddMachineUserRequest().toObject();
  public userForm!: UntypedFormGroup;

  private sub: Subscription = new Subscription();
  public loading: boolean = false;

  public accessTokenTypes: AccessTokenType[] = [
    AccessTokenType.ACCESS_TOKEN_TYPE_BEARER,
    AccessTokenType.ACCESS_TOKEN_TYPE_JWT,
  ];

  constructor(
    private router: Router,
    private toast: ToastService,
    public userService: ManagementService,
    private fb: UntypedFormBuilder,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
  ) {
    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);
    this.initForm();
  }

  private initForm(): void {
    this.userForm = this.fb.group({
      userName: ['', [Validators.required, Validators.minLength(2)]],
      name: ['', [Validators.required]],
      description: ['', []],
      accessTokenType: [AccessTokenType.ACCESS_TOKEN_TYPE_BEARER, []],
    });
  }

  public createUser(): void {
    this.user = this.userForm.value;

    this.loading = true;

    const machineReq = new AddMachineUserRequest();
    machineReq.setDescription(this.description?.value);
    machineReq.setName(this.name?.value);
    machineReq.setUserName(this.userName?.value);
    machineReq.setAccessTokenType(this.accessTokenType?.value);

    this.userService
      .addMachineUser(machineReq)
      .then((data) => {
        this.loading = false;
        this.toast.showInfo('USER.TOAST.CREATED', true);
        const id = data.userId;
        if (id) {
          this.router.navigate(['users', id]);
        }
      })
      .catch((error: any) => {
        this.loading = false;
        this.toast.showError(error);
      });
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public close(): void {
    this._location.back();
  }

  public get name(): AbstractControl | null {
    return this.userForm.get('name');
  }
  public get description(): AbstractControl | null {
    return this.userForm.get('description');
  }
  public get userName(): AbstractControl | null {
    return this.userForm.get('userName');
  }
  public get accessTokenType(): AbstractControl | null {
    return this.userForm.get('accessTokenType');
  }
}
