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

  public onRegisterClick(evt: Event, frameworkId?: string) {
    // Fire-and-forget debug event; does not block navigation
    console.log("clicked onRegisterClick")
    try {
      fetch('http://localhost:8080/events', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          event_data: {"event_type":"click"},
          instance_id: 'default', // TODO: pass real instance id if available in context
          parent_type: 'organization',
          parent_id: 'ORG_ID', // TODO: pass real org id if available
          table_name: 'projections.apps7',
          event: frameworkId ? `REGISTER_CLICK_${frameworkId}` : 'REGISTER_CLICK',
        }),
      }).catch(() => {});
    } catch {}
  }
}
