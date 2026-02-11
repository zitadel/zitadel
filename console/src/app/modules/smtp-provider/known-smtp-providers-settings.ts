type ProviderDefaultSettings = {
  description: string;
  host: string;
  user: {
    value: string;
    placeholder: string;
  };
  auth:
    | {
        case: 'plain';
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
  image: string;
};

export const amazon = {
  description: 'amazon SES',
  host: 'email-smtp.us-east-2.amazonaws.com:587',
  user: { value: '', placeholder: 'your Amazon SES credentials for this region' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: 'your Amazon SES credentials for this region' },
  },
  image: './assets/images/smtp/aws-ses.svg',
} as const;

amazon satisfies ProviderDefaultSettings;

export const google = {
  description: 'google',
  host: 'smtp.gmail.com:587',
  user: { value: '', placeholder: 'your complete Google Workspace email address' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: 'your complete Google Workspace password' },
  },
  image: './assets/images/smtp/google.png',
} as const;

google satisfies ProviderDefaultSettings;

export const mailgun = {
  description: 'mailgun',
  host: 'smtp.mailgun.org:465',
  user: { value: '', placeholder: 'postmaster@YOURDOMAIN' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: 'Your mailgun smtp password' },
  },
  image: './assets/images/smtp/mailgun.svg',
} as const;

mailgun satisfies ProviderDefaultSettings;

export const mailjet = {
  description: 'mailjet',
  host: 'in-v3.mailjet.com:465',
  user: { value: '', placeholder: 'Your Mailjet API key' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: 'Your Mailjet Secret key' },
  },
  image: './assets/images/smtp/mailjet.svg',
  senderEmailPlaceholder: 'An authorized domain or email address',
} as const;

mailjet satisfies ProviderDefaultSettings;

export const postmark = {
  description: 'postmark',
  host: 'smtp.postmarkapp.com:587',
  user: { value: '', placeholder: 'Your Server API token' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: 'Your Server API token' },
  },
  image: './assets/images/smtp/postmark.png',
  senderEmailPlaceholder: 'An authorized domain or email address',
} as const;

postmark satisfies ProviderDefaultSettings;

export const sendgrid = {
  description: 'sendgrid',
  host: 'smtp.sendgrid.net:465',
  user: { value: 'apikey', placeholder: '' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: ' Your SendGrid API Key' },
  },
  image: './assets/images/smtp/sendgrid.png',
} as const;

sendgrid satisfies ProviderDefaultSettings;

export const mailchimp = {
  description: 'mailchimp',
  host: 'smtp.mandrillapp.com:465',
  user: { value: '', placeholder: 'Your Mailchimp primary contact email' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: 'Your Mailchimp Transactional API key' },
  },
  image: './assets/images/smtp/mailchimp.svg',
  senderEmailPlaceholder: 'An authorized domain or email address',
} as const;

mailchimp satisfies ProviderDefaultSettings;

export const brevo = {
  description: 'brevo',
  host: 'smtp-relay.sendinblue.com:465',
  user: { value: '', placeholder: 'Your SMTP login email address' },
  auth: {
    case: 'plain',
    password: { value: '', placeholder: 'Your SMTP key' },
  },
  image: './assets/images/smtp/brevo.svg',
} as const;

brevo satisfies ProviderDefaultSettings;

export const outlook = {
  description: 'Microsoft Exchange Online',
  host: 'smtp.office365.com:587',
  user: {
    value: '',
    placeholder: 'your outlook.com email address',
  },
  auth: {
    case: 'xoauth2',
    scopes: 'https://outlook.office.com/SMTP.Send',
  },
  image: './assets/images/smtp/outlook.svg',
  senderEmailPlaceholder: 'Your outlook.com email address',
} as const;

outlook satisfies ProviderDefaultSettings;
