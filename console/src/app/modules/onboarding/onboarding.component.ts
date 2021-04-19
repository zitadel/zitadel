import { Component } from '@angular/core';

@Component({
    selector: 'cnsl-onboarding',
    templateUrl: './onboarding.component.html',
    styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent {
    public steps: Array<{
        titleI18nKey: string;
        descI18nKey: string,
        docs?: string;
        link?: string[];
    }> = [
            { titleI18nKey: 'ONBOARDING.STEPS.1.TITLE', descI18nKey: 'ONBOARDING.STEPS.1.DESC', docs: 'https://docs.zitadel.ch/use', link: ['/projects', 'create'] },
            { titleI18nKey: 'ONBOARDING.STEPS.2.TITLE', descI18nKey: 'ONBOARDING.STEPS.2.DESC', docs: 'https://docs.zitadel.ch/use', link: ['/projects'] },
            { titleI18nKey: 'ONBOARDING.STEPS.3.TITLE', descI18nKey: 'ONBOARDING.STEPS.3.DESC', link: ['/iam', 'policies'] },
        ];
}
