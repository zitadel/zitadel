import { CommonModule } from '@angular/common';
import { Component, EventEmitter, OnDestroy, OnInit, Output, effect, signal } from '@angular/core';
import { ActivatedRoute, Params, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { MatButtonModule } from '@angular/material/button';
import { listFrameworks, hasFramework, getFramework } from '@netlify/framework-info';
import { FrameworkName } from '@netlify/framework-info/lib/generated/frameworkNames';
import { FrameworkAutocompleteComponent } from '../framework-autocomplete/framework-autocomplete.component';
import { Framework } from '../quickstart/quickstart.component';
import { BehaviorSubject, Subject, takeUntil } from 'rxjs';

@Component({
  standalone: true,
  selector: 'cnsl-framework-change',
  templateUrl: './framework-change.component.html',
  styleUrls: ['./framework-change.component.scss'],
  imports: [TranslateModule, RouterModule, CommonModule, MatButtonModule, FrameworkAutocompleteComponent],
})
export class FrameworkChangeComponent implements OnInit, OnDestroy {
  private destroy$: Subject<void> = new Subject();
  public framework: BehaviorSubject<Framework | undefined> = new BehaviorSubject<Framework | undefined>(undefined);
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

  constructor(private activatedRoute: ActivatedRoute) {
    this.framework.pipe(takeUntil(this.destroy$)).subscribe((value) => {
      this.frameworkChanged.emit(value);
    });
  }

  public ngOnInit() {
    this.activatedRoute.queryParams.pipe(takeUntil(this.destroy$)).subscribe((params: Params) => {
      const { framework } = params;
      console.log(framework);
      if (framework) {
        console.log(this.frameworks);
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
}
