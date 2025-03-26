import { CommonModule } from '@angular/common';
import { Component, EventEmitter, OnDestroy, OnInit, Output } from '@angular/core';
import { ActivatedRoute, Params, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { MatButtonModule } from '@angular/material/button';
import { Framework } from '../quickstart/quickstart.component';
import { BehaviorSubject, Subject, takeUntil } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { FrameworkChangeDialogComponent } from './framework-change-dialog.component';

@Component({
  standalone: true,
  selector: 'cnsl-framework-change',
  templateUrl: './framework-change.component.html',
  styleUrls: ['./framework-change.component.scss'],
  imports: [TranslateModule, RouterModule, CommonModule, MatButtonModule],
})
export class FrameworkChangeComponent implements OnInit, OnDestroy {
  private destroy$: Subject<void> = new Subject();
  public framework: BehaviorSubject<Framework | undefined> = new BehaviorSubject<Framework | undefined>(undefined);
  @Output() public frameworkChanged: EventEmitter<Framework> = new EventEmitter();
  public frameworks: Framework[] = frameworkDefinition.map((f) => {
    return {
      ...f,
      fragment: '',
      imgSrcDark: `assets${f.imgSrcDark}`,
      imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
    };
  });

  constructor(
    private activatedRoute: ActivatedRoute,
    private dialog: MatDialog,
  ) {
    this.framework.pipe(takeUntil(this.destroy$)).subscribe((value) => {
      this.frameworkChanged.emit(value);
    });
  }

  public ngOnInit() {
    this.activatedRoute.queryParams.pipe(takeUntil(this.destroy$)).subscribe((params: Params) => {
      const { framework } = params;
      if (framework) {
        this.findFramework(framework);
      }
    });
  }

  public ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public findFramework(id: string) {
    const temp = this.frameworks.find((f) => f.id === id);
    this.framework.next(temp);
    this.frameworkChanged.emit(temp);
  }

  public openDialog(): void {
    const ref = this.dialog.open(FrameworkChangeDialogComponent, {
      width: '400px',
      data: {
        framework: this.framework.value,
        frameworks: this.frameworks,
      },
    });

    ref.afterClosed().subscribe((resp) => {
      if (resp) {
        this.framework.next(resp);
      }
    });
  }
}
