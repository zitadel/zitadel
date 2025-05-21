import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import type { FrameworkName } from '@netlify/framework-info/lib/generated/frameworkNames';
import { AddOIDCAppRequest } from 'src/app/proto/generated/zitadel/management_pb';

export type FrameworkDefinition = {
  id?: FrameworkName | string;
  title: string;
  imgSrcDark: string;
  imgSrcLight?: string;
  docsLink: string;
  external?: boolean;
};

export type Framework = FrameworkDefinition & {
  fragment: string;
};

@Component({
  standalone: true,
  selector: 'cnsl-oidc-app-configuration',
  templateUrl: './oidc-configuration.component.html',
  styleUrls: ['./oidc-configuration.component.scss'],
  imports: [TranslateModule, RouterModule, CommonModule, MatButtonModule],
})
export class OIDCConfigurationComponent {
  @Input() public name?: string;
  @Input() public configuration: AddOIDCAppRequest.AsObject = new AddOIDCAppRequest().toObject();
  @Output() public changeName: EventEmitter<string> = new EventEmitter();
}
