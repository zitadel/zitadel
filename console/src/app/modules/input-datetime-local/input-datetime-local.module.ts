import {forwardRef, NgModule} from "@angular/core";
import {DateTimeLocalInputComponent} from "./input-datetime-local.component";
import {NG_VALUE_ACCESSOR} from "@angular/forms";

@NgModule({
  declarations: [DateTimeLocalInputComponent],
  imports: [],
  exports: [DateTimeLocalInputComponent],
})
export class DateTimeLocalInputModule {}
