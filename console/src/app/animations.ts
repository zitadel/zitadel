import {
  animate,
  animateChild,
  AnimationTriggerMetadata,
  group,
  query,
  stagger,
  style,
  transition,
  trigger,
} from '@angular/animations';


export const toolbarAnimation: AnimationTriggerMetadata =
  trigger('toolbar', [
    transition(':enter', [
      style({
        transform: 'translateY(-100%)',
        opacity: 0,
      }),
      animate(
        '.2s ease-out',
        style({
          transform: 'translateY(0%)',
          opacity: 1,
        }),
      ),
    ]),
  ]);

export const adminLineAnimation: AnimationTriggerMetadata =
  trigger('adminline', [
    transition(':enter', [
      style({
        transform: 'translateY(100%)',
        opacity: 0.5,
      }),
      animate(
        '.2s ease-out',
        style({
          transform: 'translateY(0%)',
          opacity: 1,
        }),
      ),
    ]),
  ]);

export const accountCard: AnimationTriggerMetadata = trigger('accounts', [
  transition(':enter', [
    style({
      transform: 'scale(.9) translateY(-10%)',
      height: '200px',
      opacity: 0,
    }),
    animate(
      '.1s ease-out',
      style({
        transform: 'scale(1) translateY(0%)',
        height: '*',
        opacity: 1,
      }),
    ),
  ]),
]);

export const navAnimations: Array<AnimationTriggerMetadata> = [
  trigger('navAnimation', [
    transition('* => *', [
      query('@navitem', stagger('50ms', animateChild()), { optional: true }),
    ]),
  ]),
  trigger('navitem', [
    transition(':enter', [
      style({
        opacity: 0,
      }),
      animate(
        '.0s',
        style({
          opacity: 1,
        }),
      ),
    ]),
    transition(':leave', [
      style({
        opacity: 1,
      }),
      animate(
        '.0s',
        style({
          opacity: 0,
        }),
      ),
    ]),
  ]),
];


export const enterAnimations: Array<AnimationTriggerMetadata> = [
  trigger('appearfade', [
    transition(':enter', [
      style({
        transform: 'scale(.9) translateY(-10%)',
        opacity: 0,
      }),
      animate(
        '100ms ease-in-out',
        style({
          transform: 'scale(1) translateY(0%)',
          opacity: 1,
        }),
      ),
    ]),
    transition(':leave', [
      style({
        transform: 'scale(1) translateY(0%)',
        opacity: 1,
      }),
      animate(
        '100ms ease-in-out',
        style({
          transform: 'scale(.9) translateY(-10%)',
          opacity: 0,
        }),
      ),
    ]),
  ]),
];

export const routeAnimations: AnimationTriggerMetadata = trigger('routeAnimations', [
  transition('HomePage => AddPage', [
    style({ transform: 'translateX(100%)', opacity: 0.5 }),
    animate('250ms ease-out', style({ transform: 'translateX(0%)', opacity: 1 })),
  ]),
  transition('AddPage => HomePage',
    [animate('250ms', style({ transform: 'translateX(100%)', opacity: 0.5 }))],
  ),
  transition('HomePage => DetailPage', [
    query(':enter, :leave', style({ position: 'absolute', left: 0, right: 0 }), {
      optional: true,
    }),
    group([
      query(
        ':enter',
        [
          style({
            transform: 'translateX(20%)',
            opacity: 0.5,
          }),
          animate(
            '.35s ease-in',
            style({
              transform: 'translateX(0%)',
              opacity: 1,
            }),
          ),
        ],
        {
          optional: true,
        },
      ),
      query(
        ':leave',
        [style({ opacity: 1, width: '100%' }), animate('.35s ease-out', style({ opacity: 0 }))],
        {
          optional: true,
        },
      ),
    ]),
  ]),
  transition('DetailPage => HomePage', [
    query(':enter, :leave', style({ position: 'absolute', left: 0, right: 0 }), {
      optional: true,
    }),
    group([
      query(
        ':enter',
        [
          style({
            opacity: 0,
          }),
          animate(
            '.35s ease-out',
            style({
              opacity: 1,
            }),
          ),
        ],
        {
          optional: true,
        },
      ),
      query(
        ':leave',
        [
          style({ width: '100%', transform: 'translateX(0%)' }),
          animate('.35s ease-in', style({ transform: 'translateX(30%)', opacity: 0 })),
        ],
        {
          optional: true,
        },
      ),
    ]),
  ]),
]);
