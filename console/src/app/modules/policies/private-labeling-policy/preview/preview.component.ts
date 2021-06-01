import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';

import { Preview, Theme } from '../private-labeling-policy.component';

@Component({
  selector: 'cnsl-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit {
  @Input() preview: Preview = Preview.PREVIEW;
  @Input() policy!: LabelPolicy.AsObject;
  @Input() label: string = 'PREVIEW';
  @Input() images: { [imagekey: string]: any; } = {};
  @Input() theme: Theme = Theme.DARK;
  @Input() refresh: Observable<void> = of();
  private destroyed$: Subject<void> = new Subject();
  public Theme: any = Theme;
  public Preview: any = Preview;
  constructor(private chd: ChangeDetectorRef) { }

  public ngOnInit(): void {
    this.refresh.pipe(takeUntil(this.destroyed$)).subscribe(() => {
      this.chd.detectChanges();
    });
  }

  public ngOnDestroy(): void {
    this.destroyed$.next();
    this.destroyed$.complete();
  }
}
