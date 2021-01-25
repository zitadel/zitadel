import { Component, OnInit } from '@angular/core';
import { AuthenticationService } from 'src/app/services/authentication.service';

@Component({
    selector: 'cnsl-onboarding',
    templateUrl: './onboarding.component.html',
    styleUrls: ['./onboarding.component.scss']
})
export class OnboardingComponent implements OnInit {
    public steps = [
        { titleI18nKey: 'ONBOARDING.STEPS.1.TITLE', descI18nKey: 'ONBOARDING.STEPS.1.DESC', docs: "https://docs.zitadel.ch/use" },
        { titleI18nKey: 'ONBOARDING.STEPS.2.TITLE', descI18nKey: 'ONBOARDING.STEPS.2.DESC', docs: "https://docs.zitadel.ch/use" },
        { titleI18nKey: 'ONBOARDING.STEPS.3.TITLE', descI18nKey: 'ONBOARDING.STEPS.3.DESC' },
        { titleI18nKey: 'ONBOARDING.STEPS.4.TITLE', descI18nKey: 'ONBOARDING.STEPS.4.DESC' },
        { titleI18nKey: 'ONBOARDING.STEPS.5.TITLE', descI18nKey: 'ONBOARDING.STEPS.5.DESC' },
    ];
    constructor(public authenticationService: AuthenticationService) { }

    ngOnInit(): void {
    }

}
