import { OnboardingActions } from '../services/admin.service';
import { COLORS } from './color';
import { MilestoneType } from '../proto/generated/zitadel/milestone/v1/milestone_pb';

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

const sthdark: string = COLORS[18][700];
const sthlight = COLORS[18][200];

export const ONBOARDING_MILESTONES: OnboardingActions[] = [
  {
    order: 0,
    milestoneType: MilestoneType.MILESTONE_TYPE_PROJECT_CREATED,
    link: '/projects/create',
    iconClasses: 'las la-database',
    darkcolor: greendark,
    lightcolor: greenlight,
  },
  {
    order: 1,
    milestoneType: MilestoneType.MILESTONE_TYPE_APPLICATION_CREATED,
    link: '/projects/app-create',
    iconClasses: 'lab la-openid',
    darkcolor: purpledark,
    lightcolor: purplelight,
  },
  {
    order: 3,
    milestoneType: MilestoneType.MILESTONE_TYPE_AUTHENTICATION_SUCCEEDED_ON_APPLICATION,
    link: 'https://zitadel.com/docs/guides/integrate/login-users',
    externalLink: true,
    iconClasses: 'las la-sign-in-alt',
    darkcolor: sthdark,
    lightcolor: sthlight,
  } /*
  {
    order: 4,
    milestoneType: 'user.human.added',
    link: '/users/create',
    iconClasses: 'las la-user',
    darkcolor: bluedark,
    lightcolor: bluelight,
  },
  {
    order: 5,
    milestoneType: 'user.grant.added',
    link: '/grant-create',
    iconClasses: 'las la-shield-alt',
    darkcolor: reddark,
    lightcolor: redlight,
  },
  {
    order: 6,
    milestoneType: 'instance.policy.label.added',
    link: '/settings',
    fragment: 'branding',
    iconClasses: 'las la-swatchbook',
    darkcolor: pinkdark,
    lightcolor: pinklight,
  },
  {
    order: 7,
    milestoneType: 'instance.smtp.config.added',
    link: '/settings',
    fragment: 'smtpprovider',
    iconClasses: 'las la-envelope',
    darkcolor: yellowdark,
    lightcolor: yellowlight,
  },*/,
];
