import { Component, Inject } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatSelectChange } from '@angular/material/select';
import { SubscriptionService } from 'src/app/services/subscription.service';

import { COUNTRIES, Country } from '../country';

function compare(a: Country, b: Country) {
    if (a.isoCode < b.isoCode) {
        return -1;
    }
    if (a.isoCode > b.isoCode) {
        return 1;
    }
    return 0;
}

@Component({
    selector: 'app-payment-info-dialog',
    templateUrl: './payment-info-dialog.component.html',
    styleUrls: ['./payment-info-dialog.component.scss'],
})
export class PaymentInfoDialogComponent {
    public stripeLoading: boolean = false;
    public COUNTRIES: Country[] = COUNTRIES.sort(compare);
    public form!: FormGroup;

    private orgId: string = '';

    constructor(
        private subService: SubscriptionService,
        private fb: FormBuilder,
        public dialogRef: MatDialogRef<PaymentInfoDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
        this.orgId = data.orgId;

        this.form = this.fb.group({
            contact: ['', [Validators.required]],
            company: ['', []],
            address: ['', [Validators.required]],
            city: ['', [Validators.required]],
            postal_code: ['', [Validators.required]],
            country: ['', [Validators.required]],
        });

        if (data.customer) {
            this.form.patchValue(data.customer);
        }

        if (!data.customer?.country) {
            this.form.get('country')?.setValue('CH');
        }

        this.getLink();
    }

    public getLink(): void {
        if (this.orgId) {
            this.stripeLoading = true;
            this.subService.getLink(this.orgId, window.location.href)
                .then(payload => {
                    this.stripeLoading = false;
                    console.log(payload);
                    if (payload.redirect_url) {
                        window.open(payload.redirect_url, '_blank');
                    }
                })
                .catch(error => {
                    this.stripeLoading = false;
                    console.error(error);
                });
        }
    }

    public changeCountry(selection: MatSelectChange): void {
        const country = COUNTRIES.find(c => c.isoCode === selection.value);
        if (country && country.phoneCode !== undefined && this.phone && this.phone.value !== `+${country.phoneCode}`) {
            this.phone.setValue(`+${country.phoneCode}`);
        }
    }

    submitAndCloseDialog(): void {
        this.dialogRef.close(this.form.value);
    }

    get phone(): AbstractControl | null {
        return this.form.get('phone');
    }
}
