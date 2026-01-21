import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, effect, OnInit } from '@angular/core';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { map, ReplaySubject } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import {
  FrameworkChangeDialogComponent,
  FrameworkChangeDialogData,
  FrameworkChangeDialogResult,
} from './framework-change-dialog.component';
import { outputFromObservable, takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop';
import { frameworks } from 'src/app/utils/framework';
import { filter } from 'rxjs/operators';

type Framework = (typeof frameworks)[number];

@Component({
  selector: 'cnsl-framework-change',
  templateUrl: './framework-change.component.html',
  styleUrls: ['./framework-change.component.scss'],
  imports: [TranslateModule, RouterModule, CommonModule, MatButtonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FrameworkChangeComponent {
  protected readonly frameworks = frameworks;
  protected readonly framework$ = new ReplaySubject<Framework>(1);
  public readonly frameworkChanged = outputFromObservable(this.framework$);

  constructor(
    private activatedRoute: ActivatedRoute,
    private dialog: MatDialog,
  ) {
    const frameworkSignal = toSignal(this.activatedRoute.queryParamMap.pipe(map((params) => params.get('framework'))), {
      initialValue: null,
    });
    effect(() => {
      const framework = this.frameworks.find((f) => f.id === frameworkSignal());
      if (framework) {
        this.framework$.next(framework);
      }
    });
  }

  public openDialog(framework: Framework | null) {
    this.dialog
      .open<FrameworkChangeDialogComponent, FrameworkChangeDialogData, FrameworkChangeDialogResult>(
        FrameworkChangeDialogComponent,
        {
          width: '400px',
          data: framework,
        },
      )
      .afterClosed()
      .pipe(filter(Boolean), takeUntilDestroyed())
      .subscribe((framework) => {
        this.framework$.next(framework);
      });
  }
}
