import { ChangeDetectionStrategy, Component, Input, Pipe, PipeTransform } from '@angular/core';
import { AsyncPipe, NgForOf, NgIf } from '@angular/common';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { LabelModule } from 'src/app/modules/label/label.module';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatOptionModule } from '@angular/material/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TranslateModule } from '@ngx-translate/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { fromEvent, map, mergeWith, Observable } from 'rxjs';
import { startWith } from 'rxjs/operators';

@Pipe({ standalone: true, name: 'filter' })
class Filter implements PipeTransform {
  transform(items: string[] | undefined = [], input: HTMLInputElement): Observable<string[]> {
    const focus$ = fromEvent(input, 'focus').pipe(map(() => ''));

    return fromEvent(input, 'input').pipe(
      startWith(undefined),
      map(() => input.value.toLowerCase()),
      mergeWith(focus$),
      map((input) => items.filter((item) => item.toLowerCase().includes(input))),
    );
  }
}

@Component({
  selector: 'cnsl-actions-two-add-action-autocomplete-input',
  templateUrl: './actions-two-add-action-autocomplete-input.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [
    AsyncPipe,
    Filter,
    FormFieldModule,
    InputModule,
    LabelModule,
    MatAutocompleteModule,
    MatOptionModule,
    MatProgressSpinnerModule,
    TranslateModule,
    NgIf,
    NgForOf,
    ReactiveFormsModule,
  ],
})
export class ActionsTwoAddActionAutocompleteInputComponent {
  @Input({ required: true })
  public label!: string;

  @Input({ required: true })
  public items: string[] | undefined;

  @Input({ required: true })
  public control!: FormControl<string>;
}
