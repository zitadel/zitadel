import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import type { getFramework } from '@netlify/build-info';
import { FRAMEWORK_DEFINITION } from 'src/app/utils/framework';

type FrameworkName = Parameters<typeof getFramework>[0];

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
  selector: 'cnsl-quickstart',
  templateUrl: './quickstart.component.html',
  styleUrls: ['./quickstart.component.scss'],
  imports: [TranslateModule, RouterModule, MatButtonModule],
})
export class QuickstartComponent {
  public frameworks: FrameworkDefinition[] = FRAMEWORK_DEFINITION;
}
