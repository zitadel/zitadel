import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { GrantedProject, ProjectGrantState, Role } from 'src/app/proto/generated/zitadel/project_pb';

@Component({
  selector: 'cnsl-project-grant-illustration',
  templateUrl: './project-grant-illustration.component.html',
  styleUrls: ['./project-grant-illustration.component.scss'],
})
export class ProjectGrantIllustrationComponent implements OnInit {
  @Input() public grantedProject!: GrantedProject.AsObject;
  @Input() public projectRoleOptions: Role.AsObject[] = [];
  @Output() public roleRemoved: EventEmitter<string> = new EventEmitter();
  @Output() public editRoleClicked: EventEmitter<void> = new EventEmitter();

  ProjectGrantState: any = ProjectGrantState;
  constructor() {}

  ngOnInit(): void {}

  public removeRole(roleKey: string): void {
    this.roleRemoved.emit(roleKey);
  }
}
