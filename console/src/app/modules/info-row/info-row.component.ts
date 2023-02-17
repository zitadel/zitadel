import { Component, Input } from '@angular/core';
import { App, AppState } from 'src/app/proto/generated/zitadel/app_pb';
import { IDP, IDPState } from 'src/app/proto/generated/zitadel/idp_pb';
import { InstanceDetail, State } from 'src/app/proto/generated/zitadel/instance_pb';
import { Org, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { LoginPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { GrantedProject, Project, ProjectGrantState, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { User, UserState } from 'src/app/proto/generated/zitadel/user_pb';

@Component({
  selector: 'cnsl-info-row',
  templateUrl: './info-row.component.html',
  styleUrls: ['./info-row.component.scss'],
})
export class InfoRowComponent {
  @Input() public user!: User.AsObject;
  @Input() public org!: Org.AsObject;
  @Input() public instance!: InstanceDetail.AsObject;
  @Input() public app!: App.AsObject;
  @Input() public idp!: IDP.AsObject;
  @Input() public project!: Project.AsObject;
  @Input() public grantedProject!: GrantedProject.AsObject;
  @Input() public loginPolicy?: LoginPolicy.AsObject;

  public UserState: any = UserState;
  public State: any = State;
  public OrgState: any = OrgState;
  public AppState: any = AppState;
  public IDPState: any = IDPState;
  public ProjectState: any = ProjectState;
  public ProjectGrantState: any = ProjectGrantState;

  public copied: string = '';

  constructor() {}

  public get loginMethods(): Set<string> {
    const methods = this.user?.loginNamesList;
    let email: string = '';
    let phone: string = '';
    if (this.loginPolicy) {
      if (
        !this.loginPolicy?.disableLoginWithEmail &&
        this.user.human?.email?.email &&
        this.user.human.email.isEmailVerified
      ) {
        email = this.user.human?.email?.email;
      }
      if (
        !this.loginPolicy?.disableLoginWithPhone &&
        this.user.human?.phone?.phone &&
        this.user.human.phone.isPhoneVerified
      ) {
        phone = this.user.human?.phone?.phone;
      }
    }
    return new Set([email, phone, ...methods].filter((method) => !!method));
  }
}
