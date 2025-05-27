import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { MatButtonModule } from '@angular/material/button';
import type { FrameworkName } from '@netlify/framework-info/lib/generated/frameworkNames';
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';

export type FrameworkDefinition = {
  id?: FrameworkName | string;
  title: string;
  description?: string;
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
