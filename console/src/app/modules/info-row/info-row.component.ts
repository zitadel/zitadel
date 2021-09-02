import { Component, Input, OnInit } from '@angular/core';
import { User, UserState } from 'src/app/proto/generated/zitadel/user_pb';

@Component({
  selector: 'cnsl-info-row',
  templateUrl: './info-row.component.html',
  styleUrls: ['./info-row.component.scss']
})
export class InfoRowComponent implements OnInit {
  @Input() public user!: User.AsObject;
  public UserState: any = UserState;
  public copied: string = '';

  constructor() { }

  ngOnInit(): void {
  }

}
