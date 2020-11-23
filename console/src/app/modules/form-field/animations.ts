import { animate, AnimationTriggerMetadata, state, style, transition, trigger } from '@angular/animations';

/**
 * Animations used by the CnslFormFieldComponent.
 */
export const cnslFormFieldAnimations: {
    readonly transitionMessages: AnimationTriggerMetadata;
} = {
    /** Animation that transitions the form field's error and hint messages. */
    transitionMessages: trigger('transitionMessages', [
        // TODO(mmalerba): Use angular animations for label animation as well.
        state('enter', style({ opacity: 1, transform: 'translateY(0%)' })),
        transition('void => enter', [
            style({ opacity: 0, transform: 'translateY(-100%)' }),
            animate('3000ms cubic-bezier(0.55, 0, 0.55, 0.2)'),
        ]),
    ]),
};
