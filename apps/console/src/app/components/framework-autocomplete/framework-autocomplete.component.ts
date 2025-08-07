import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Observable, map, of, startWith, switchMap, tap } from 'rxjs';
import { Framework } from '../quickstart/quickstart.component';

@Component({
  standalone: true,
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
export class FrameworkAutocompleteComponent implements OnInit {
  public isLoading = signal(false);
  @Input() public frameworkId?: string;
  @Input() public frameworks: Framework[] = [];
  @Input() public withCustom: boolean = false;
  public myControl: FormControl = new FormControl();
  @Output() public selectionChanged: EventEmitter<string> = new EventEmitter();
  public filteredOptions: Observable<Framework[]> = of([]);

  constructor() {}

  public ngOnInit() {
    this.filteredOptions = this.myControl.valueChanges.pipe(
      startWith(''),
      map((value) => {
        return this._filter(value || '');
      }),
    );
  }

  private _filter(value: string): Framework[] {
    const filterValue = value.toLowerCase();
    return this.frameworks
      .filter((option) => option.id)
      .filter((option) => option.title.toLowerCase().includes(filterValue));
  }

  public selected(event: MatAutocompleteSelectedEvent): void {
    this.selectionChanged.emit(event.option.value);
  }
}
