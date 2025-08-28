import { Component, Input } from '@angular/core';
import { App, AppState } from 'src/app/proto/generated/zitadel/app_pb';
import { IDP, IDPState } from 'src/app/proto/generated/zitadel/idp_pb';
import { InstanceDetail, State } from 'src/app/proto/generated/zitadel/instance_pb';
import { Org, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { LoginPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { GrantedProject, Project, ProjectGrantState, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { User as UserV1 } from '@zitadel/proto/zitadel/user_pb';
import { User as UserV2 } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { LoginPolicy as LoginPolicyV2 } from '@zitadel/proto/zitadel/policy_pb';

@Component({
  selector: 'cnsl-info-row',
  templateUrl: './info-row.component.html',
  styleUrls: ['./info-row.component.scss'],
})
export class InfoRowComponent {
  @Input() public user?: User.AsObject | UserV2 | UserV1;
  @Input() public org!: Org.AsObject;
  @Input() public instance!: InstanceDetail.AsObject;
  @Input() public app!: App.AsObject;
  @Input() public idp!: IDP.AsObject;
  @Input() public project!: Project.AsObject;
  @Input() public grantedProject!: GrantedProject.AsObject;
  @Input() public loginPolicy?: LoginPolicy.AsObject | LoginPolicyV2;

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
    if (!this.user) {
      return new Set();
    }

    const methods = '$typeName' in this.user ? this.user.loginNames : this.user.loginNamesList;

    const loginPolicy = this.loginPolicy;
    if (!loginPolicy) {
      return new Set([...methods]);
    }

    let email = !loginPolicy.disableLoginWithEmail ? this.getEmail(this.user) : '';
    let phone = !loginPolicy.disableLoginWithPhone ? this.getPhone(this.user) : '';

    return new Set([email, phone, ...methods].filter(Boolean));
  }

  public get userId() {
    if (!this.user) {
      return undefined;
    }
    if ('$typeName' in this.user && this.user.$typeName === 'zitadel.user.v2.User') {
      return this.user.userId;
    }
    return this.user.id;
  }

  public get changeDate() {
    return this.user?.details?.changeDate;
  }

  public get creationDate() {
    return this.user?.details?.creationDate;
  }

  private getEmail(user: User.AsObject | UserV2 | UserV1) {
    const human = this.human(user);
    if (!human) {
      return '';
    }
    if ('$typeName' in human && human.$typeName === 'zitadel.user.v2.HumanUser') {
      return human.email?.isVerified ? human.email.email : '';
    }
    return human.email?.isEmailVerified ? human.email.email : '';
  }

  private getPhone(user: User.AsObject | UserV2 | UserV1) {
    const human = this.human(user);
    if (!human) {
      return '';
    }
    if ('$typeName' in human && human.$typeName === 'zitadel.user.v2.HumanUser') {
      return human.phone?.isVerified ? human.phone.phone : '';
    }
    return human.phone?.isPhoneVerified ? human.phone.phone : '';
  }

  public human(user: User.AsObject | UserV2 | UserV1) {
    if (!('$typeName' in user)) {
      return user.human;
    }
    return user.type.case === 'human' ? user.type.value : undefined;
  }

  public isV2(user: User.AsObject | UserV2 | UserV1) {
    if ('$typeName' in user) {
      return user.$typeName === 'zitadel.user.v2.User';
    }
    return false;
  }
}
