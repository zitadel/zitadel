import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import { UserService } from 'src/app/services/user.service';
import { FormBuilder, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'cnsl-project-create',
  templateUrl: './project-create.component.html',
  styleUrls: ['./project-create.component.scss'],
  standalone: false,
})
export class ProjectCreateComponent {
  protected readonly form: ReturnType<typeof this.buildForm>;

  constructor(
    private readonly router: Router,
    private readonly toast: ToastService,
    private readonly newMgmtService: NewMgmtService,
    private readonly fb: FormBuilder,
    protected readonly location: Location,
    protected readonly userService: UserService,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);

    this.form = this.buildForm();
  }

  private buildForm() {
    return this.fb.group({
      name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
      selfAccount: new FormControl(true, { nonNullable: true, validators: [Validators.required] }),
    });
  }

  protected saveProject(): void {
    const { name, selfAccount } = this.form.getRawValue();

    this.newMgmtService
      .addProject({
        name,
        admins: selfAccount ? [{ userId: this.userService.userId() }] : undefined,
      })
      .then((resp) => {
        this.toast.showInfo('PROJECT.TOAST.CREATED', true);
        return this.router.navigate(['projects', resp.id], { queryParams: { new: true } });
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }
}
