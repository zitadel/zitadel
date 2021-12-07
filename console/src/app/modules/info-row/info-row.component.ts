import { HttpClient } from '@angular/common/http';
import { Component, Input, OnInit } from '@angular/core';
import { App, AppState } from 'src/app/proto/generated/zitadel/app_pb';
import { IDP, IDPState } from 'src/app/proto/generated/zitadel/idp_pb';
import { Org, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { GrantedProject, Project, ProjectGrantState, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { User, UserState } from 'src/app/proto/generated/zitadel/user_pb';

@Component({
  selector: 'cnsl-info-row',
  templateUrl: './info-row.component.html',
  styleUrls: ['./info-row.component.scss'],
})
export class InfoRowComponent implements OnInit {
  @Input() public user!: User.AsObject;
  @Input() public org!: Org.AsObject;
  @Input() public app!: App.AsObject;
  @Input() public idp!: IDP.AsObject;
  @Input() public project!: Project.AsObject;
  @Input() public grantedProject!: GrantedProject.AsObject;

  public UserState: any = UserState;
  public OrgState: any = OrgState;
  public AppState: any = AppState;
  public IDPState: any = IDPState;
  public ProjectState: any = ProjectState;
  public ProjectGrantState: any = ProjectGrantState;

  public copied: string = '';

  public environmentMap: { [key: string]: string } = {};

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    if (this.app) {
      this.http
        .get('./assets/environment.json')
        .toPromise()
        .then((env: any) => {
          this.environmentMap = {
            issuer: env.issuer,
            adminServiceUrl: env.adminServiceUrl,
            mgmtServiceUrl: env.mgmtServiceUrl,
            authServiceUrl: env.adminServiceUrl,
          };
        });
    }
  }
}
