import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { MatButtonModule } from '@angular/material/button';
import type { getFramework } from '@netlify/build-info';
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';

type FrameworkName = Parameters<typeof getFramework>[0];

export type FrameworkDefinition = {
  id?: FrameworkName | string;
  title: string;
  description?: string;
  ecosystem?: string;
  imgSrcDark: string;
  imgSrcLight?: string;
  docsLink: string;
  external?: boolean;
  client?: boolean;
  sdk?: boolean;
  sdkLink?: string;
  sdkName?: string;
  sdkPackage?: string;
  sdkCommand?: string;
  localhost?: string;
  buildCommand?: string;
  startCommand?: string;
  example?: string;
  exampleLink?: string;
  envSetup?: {
    type: string;
    filename?: string;
    description: string;
    variables: Array<{
      name: string;
      description: string;
      placeholder: string;
      required: boolean;
    }>;
  };
};

export type Framework = FrameworkDefinition & {
  fragment: string;
};

@Component({
  selector: 'cnsl-quickstart',
  templateUrl: './quickstart.component.html',
  styleUrls: ['./quickstart.component.scss'],
  imports: [TranslateModule, RouterModule, CommonModule, MatButtonModule],
})
export class QuickstartComponent {
  public frameworks: FrameworkDefinition[] = frameworkDefinition
    .filter((f) => f.id && OIDC_CONFIGURATIONS[f.id])
    .map((f) => {
      return {
        ...f,
        imgSrcDark: `assets${f.imgSrcDark}`,
        imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
      };
    });
}
