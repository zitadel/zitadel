import { OnboardingActions } from '../services/admin.service';
import { COLORS } from './color';

const reddark: string = COLORS[0][700];
const redlight = COLORS[0][200];

const yellowdark: string = COLORS[3][700];
const yellowlight = COLORS[3][200];

const greendark: string = COLORS[6][700];
const greenlight = COLORS[6][200];

const bluedark: string = COLORS[9][700];
const bluelight = COLORS[9][200];

const purpledark: string = COLORS[12][700];
const purplelight = COLORS[12][200];

const pinkdark: string = COLORS[15][700];
const pinklight = COLORS[15][200];

export const ONBOARDING_EVENTS: OnboardingActions[] = [
  {
    order: 0,
    eventType: 'project.added',
    oneof: ['project.added'],
    link: ['/projects/create'],
    iconClasses: 'las la-database',
    darkcolor: greendark,
    lightcolor: greenlight,
  },
  {
    order: 1,
    eventType: 'project.application.added',
    oneof: ['project.application.added'],
    link: ['/projects/app-create'],
    iconClasses: 'lab la-openid',
    darkcolor: purpledark,
    lightcolor: purplelight,
  },
  {
    order: 2,
    eventType: 'user.human.added',
    oneof: ['user.human.added'],
    link: ['/users/create'],
    iconClasses: 'las la-user',
    darkcolor: bluedark,
    lightcolor: bluelight,
  },
  {
    order: 3,
    eventType: 'user.grant.added',
    oneof: ['user.grant.added'],
    link: ['/grant-create'],
    iconClasses: 'las la-shield-alt',
    darkcolor: reddark,
    lightcolor: redlight,
  },
  {
    order: 4,
    eventType: 'instance.policy.label.added',
    oneof: ['instance.policy.label.added', 'instance.policy.label.changed'],
    link: ['/settings'],
    fragment: 'branding',
    iconClasses: 'las la-swatchbook',
    darkcolor: pinkdark,
    lightcolor: pinklight,
  },
  {
    order: 5,
    eventType: 'instance.smtp.config.added',
    oneof: ['instance.smtp.config.added', 'instance.smtp.config.changed'],
    link: ['/settings'],
    fragment: 'notifications',
    iconClasses: 'las la-envelope',
    darkcolor: yellowdark,
    lightcolor: yellowlight,
  },
];
