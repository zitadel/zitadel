import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, input, output, Signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { frameworks } from '../../utils/framework';
import { toSignal } from '@angular/core/rxjs-interop';

@Component({
  changeDetection: ChangeDetectionStrategy.OnPush,
  selector: 'cnsl-framework-autocomplete',
  templateUrl: './framework-autocomplete.component.html',
  styleUrls: ['./framework-autocomplete.component.scss'],
  imports: [
    TranslateModule,
    RouterModule,
    MatSelectModule,
    MatAutocompleteModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    FormsModule,
    CommonModule,
    MatButtonModule,
    InputModule,
  ],
})
export class FrameworkAutocompleteComponent {
  public readonly withCustom = input<boolean>(false);
  public selectionChanged = output<string>();

  protected readonly control: FormControl = new FormControl<string>('', { nonNullable: true });
  protected filteredOptions: Signal<typeof frameworks>;

  constructor() {
    const controlValue = toSignal(this.control.valueChanges, { initialValue: this.control.value });
    this.filteredOptions = computed(() => {
      const value = controlValue().toLowerCase();
      return frameworks.filter((option) => option.title.toLowerCase().includes(value));
    });
  }

  public selected(event: MatAutocompleteSelectedEvent): void {
    this.selectionChanged.emit(event.option.value);
  }
}
