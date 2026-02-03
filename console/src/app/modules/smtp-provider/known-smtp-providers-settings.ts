type ProviderDefaultSettings = {
  name: string;
  host?: string;
  ports?: {
    unencryptedPort?: number;
    encryptedPort: number;
  };
  auth?:
    | {
        case: 'plain';
        user: {
          value: string;
          placeholder: string;
        };
        password: {
          value: string;
          placeholder: string;
        };
      }
    | {
        case: 'xoauth2';
        scopes: string;
      };
  senderEmailPlaceholder?: string;
  image?: string;
};

const amazonEndpoints = {
  'US East (Ohio)': 'email-smtp.us-east-2.amazonaws.com',
  'US East (N. Virginia)': 'email-smtp.us-east-1.amazonaws.com',
  'US West (N. California)': 'email-smtp.us-west-1.amazonaws.com',
  'US West (Oregon)': 'email-smtp.us-west-2.amazonaws.com',
  'Asia Pacific (Mumbai)': 'email-smtp.ap-south-1.amazonaws.com',
  'Asia Pacific (Osaka)': 'email-smtp.ap-northeast-3.amazonaws.com',
  'Asia Pacific (Seoul)': 'email-smtp.ap-northeast-2.amazonaws.com',
  'Asia Pacific (Singapore)': 'email-smtp.ap-southeast-1.amazonaws.com',
  'Asia Pacific (Sydney)': 'email-smtp.ap-southeast-2.amazonaws.com',
  'Asia Pacific (Tokyo)': 'email-smtp.ap-northeast-1.amazonaws.com',
  'Canada (Central)': 'email-smtp.ca-central-1.amazonaws.com',
  'Europe (Frankfurt)': 'email-smtp.eu-central-1.amazonaws.com',
  'Europe (Ireland)': 'email-smtp.eu-west-1.amazonaws.com',
  'Europe (London)': 'email-smtp.eu-west-2.amazonaws.com',
  'Europe (Paris)': 'email-smtp.eu-west-3.amazonaws.com',
  'Europe (Stockholm)': 'email-smtp.eu-north-1.amazonaws.com',
  'South America (São Paulo)': 'email-smtp.sa-east-1.amazonaws.com',
} as const;

export const amazon = {
  name: 'amazon SES',
  regions: amazonEndpoints,
  ports: {
    encryptedPort: 587,
  },
  user: { value: '', placeholder: 'your Amazon SES credentials for this region' },
  password: { value: '', placeholder: 'your Amazon SES credentials for this region' },
  image: './assets/images/smtp/aws-ses.svg',
} as const;

amazon satisfies ProviderDefaultSettings;

export const google = {
  name: 'google',
  requiredTls: true,
  host: 'smtp.gmail.com',
  ports: {
    encryptedPort: 587,
  },
  auth: {
    case: 'plain',
    user: { value: '', placeholder: 'your complete Google Workspace email address' },
    password: { value: '', placeholder: 'your complete Google Workspace password' },
  },
  image: './assets/images/smtp/google.png',
} as const;

google satisfies ProviderDefaultSettings;

export const mailgun = {
  name: 'mailgun',
  requiredTls: false,
  host: 'smtp.mailgun.org',
  ports: {
    unencryptedPort: 587,
    encryptedPort: 465,
  },
  auth: {
    case: 'plain',
    user: { value: '', placeholder: 'postmaster@YOURDOMAIN' },
    password: { value: '', placeholder: 'Your mailgun smtp password' },
  },
  image: './assets/images/smtp/mailgun.svg',
} as const;

mailgun satisfies ProviderDefaultSettings;

export const mailjet = {
  name: 'mailjet',
  requiredTls: false,
  host: 'in-v3.mailjet.com',
  ports: {
    unencryptedPort: 587,
    encryptedPort: 465,
  },
  auth: {
    case: 'plain',
    user: { value: '', placeholder: 'Your Mailjet API key' },
    password: { value: '', placeholder: 'Your Mailjet Secret key' },
  },
  image: './assets/images/smtp/mailjet.svg',
  senderEmailPlaceholder: 'An authorized domain or email address',
} as const;

mailgun satisfies ProviderDefaultSettings;

export const postmark = {
  name: 'postmark',
  requiredTls: false,
  host: 'smtp.postmarkapp.com',
  ports: {
    unencryptedPort: 587,
    encryptedPort: 587,
  },
  auth: {
    case: 'plain',
    user: { value: '', placeholder: 'Your Server API token' },
    password: { value: '', placeholder: 'Your Server API token' },
  },
  image: './assets/images/smtp/postmark.png',
  senderEmailPlaceholder: 'An authorized domain or email address',
} as const;

postmark satisfies ProviderDefaultSettings;

export const sendgrid = {
  name: 'sendgrid',
  requiredTls: false,
  host: 'smtp.sendgrid.net',
  ports: {
    unencryptedPort: 587,
    encryptedPort: 465,
  },
  auth: {
    case: 'plain',
    user: { value: 'apikey', placeholder: '' },
    password: { value: '', placeholder: ' Your SendGrid API Key' },
  },
  image: './assets/images/smtp/sendgrid.png',
} as const;

sendgrid satisfies ProviderDefaultSettings;

export const mailchimp = {
  name: 'mailchimp',
  requiredTls: false,
  host: 'smtp.mandrillapp.com',
  ports: {
    unencryptedPort: 587,
    encryptedPort: 465,
  },
  auth: {
    case: 'plain',
    user: { value: '', placeholder: 'Your Mailchimp primary contact email' },
    password: { value: '', placeholder: 'Your Mailchimp Transactional API key' },
  },
  image: './assets/images/smtp/mailchimp.svg',
  senderEmailPlaceholder: 'An authorized domain or email address',
} as const;

mailchimp satisfies ProviderDefaultSettings;

export const brevo = {
  name: 'brevo',
  requiredTls: false,
  host: 'smtp-relay.sendinblue.com',
  ports: {
    unencryptedPort: 587,
    encryptedPort: 465,
  },
  auth: {
    case: 'plain',
    user: { value: '', placeholder: 'Your SMTP login email address' },
    password: { value: '', placeholder: 'Your SMTP key' },
  },
  image: './assets/images/smtp/brevo.svg',
} as const;

brevo satisfies ProviderDefaultSettings;

export const outlook = {
  name: 'outlook.com',
  requiredTls: true,
  host: 'smtp-mail.outlook.com',
  ports: {
    encryptedPort: 587,
  },
  auth: {
    case: 'xoauth2',
    scopes: 'https://outlook.office.com/SMTP.Send',
  },
  image: './assets/images/smtp/outlook.svg',
  senderEmailPlaceholder: 'Your outlook.com email address',
} as const;

outlook satisfies ProviderDefaultSettings;

export const generic = {
  name: 'generic',
  requiredTls: false,
} as const;

generic satisfies ProviderDefaultSettings;
