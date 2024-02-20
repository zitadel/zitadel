import { CommonModule } from '@angular/common';
import { Component, EventEmitter, OnDestroy, Output, effect, signal } from '@angular/core';
import { ActivatedRoute, Params, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { MatButtonModule } from '@angular/material/button';
import { listFrameworks, hasFramework, getFramework } from '@netlify/framework-info';
import { FrameworkName } from '@netlify/framework-info/lib/generated/frameworkNames';
import { FrameworkAutocompleteComponent } from '../framework-autocomplete/framework-autocomplete.component';
import { Framework } from '../quickstart/quickstart.component';
import { Subject, takeUntil } from 'rxjs';

@Component({
  standalone: true,
  selector: 'cnsl-framework-change',
  templateUrl: './framework-change.component.html',
  styleUrls: ['./framework-change.component.scss'],
  imports: [TranslateModule, RouterModule, CommonModule, MatButtonModule, FrameworkAutocompleteComponent],
})
export class FrameworkChangeComponent implements OnDestroy {
  private destroy$: Subject<void> = new Subject();
  public framework = signal<Framework | undefined>(undefined);
  public showFrameworkAutocomplete = signal<boolean>(false);
  @Output() public frameworkChanged: EventEmitter<Framework> = new EventEmitter();
  public frameworks: Framework[] = frameworkDefinition.map((f) => {
    return {
      ...f,
      fragment: '',
      imgSrcDark: `assets${f.imgSrcDark}`,
      imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
    };
  });

  constructor(activatedRoute: ActivatedRoute) {
    effect(() => {
      this.frameworkChanged.emit(this.framework());
    });

    activatedRoute.queryParams.pipe(takeUntil(this.destroy$)).subscribe((params: Params) => {
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
    this.framework.set(temp);
  }
}
