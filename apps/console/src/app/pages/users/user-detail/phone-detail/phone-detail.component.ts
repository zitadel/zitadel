import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { formatPhone } from 'src/app/utils/formatPhone';

@Component({
  selector: 'cnsl-phone-detail',
  templateUrl: './phone-detail.component.html',
  styleUrls: ['./phone-detail.component.scss'],
})
export class PhoneDetailComponent implements OnChanges {
  @Input() phone: string | undefined;
  public country: string | undefined;

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['phone'].currentValue) {
      const phoneNumber = formatPhone(changes['phone'].currentValue);
      if (phoneNumber && this.phone !== phoneNumber.phone) {
        this.phone = phoneNumber.phone;
        this.country = phoneNumber.country;
      }
    } else {
      this.country = undefined;
    }
  }
}
