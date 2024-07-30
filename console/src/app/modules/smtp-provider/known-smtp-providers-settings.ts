export interface AmazonRegionsEndpoints {
  'US East (Ohio)': string;
  'US East (N. Virginia)': string;
  'US West (N. California)': string;
  'US West (Oregon)': string;
  'Asia Pacific (Mumbai)': string;
  'Asia Pacific (Osaka)': string;
  'Asia Pacific (Seoul)': string;
  'Asia Pacific (Singapore)': string;
  'Asia Pacific (Sydney)': string;
  'Asia Pacific (Tokyo)': string;
  'Canada (Central)': string;
  'Europe (Frankfurt)': string;
  'Europe (London)': string;
  'Europe (Paris)': string;
  'Europe (Stockholm)': string;
  'South America (São Paulo)': string;
}

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
};

export interface ProviderDefaultSettings {
  name: string;
  regions?: AmazonRegionsEndpoints;
  multiHostsLabel?: string;
  requiredTls: boolean;
  host?: string;
  unencryptedPort?: number;
  encryptedPort?: number;
  user: {
    value: string;
    placeholder: string;
  };
  password: {
    value: string;
    placeholder: string;
  };
  senderEmailPlaceholder?: string;
  image?: string;
  routerLinkElement: string;
}

export const AmazonSESDefaultSettings: ProviderDefaultSettings = {
  name: 'amazon SES',
  regions: amazonEndpoints,
  multiHostsLabel: 'Choose your region',
  requiredTls: true,
  encryptedPort: 587,
  user: { value: '', placeholder: 'your Amazon SES credentials for this region' },
  password: { value: '', placeholder: 'your Amazon SES credentials for this region' },
  image: './assets/images/smtp/aws-ses.svg',
  routerLinkElement: 'aws-ses',
};

export const GoogleDefaultSettings: ProviderDefaultSettings = {
  name: 'google',
  requiredTls: true,
  host: 'smtp.gmail.com',
  unencryptedPort: 587,
  encryptedPort: 587,
  user: { value: '', placeholder: 'your complete Google Workspace email address' },
  password: { value: '', placeholder: 'your complete Google Workspace password' },
  image: './assets/images/smtp/google.png',
  routerLinkElement: 'google',
};

export const MailgunDefaultSettings: ProviderDefaultSettings = {
  name: 'mailgun',
  requiredTls: false,
  host: 'smtp.mailgun.org',
  unencryptedPort: 587,
  encryptedPort: 465,
  user: { value: '', placeholder: 'postmaster@YOURDOMAIN' },
  password: { value: '', placeholder: 'Your mailgun smtp password' },
  image: './assets/images/smtp/mailgun.svg',
  routerLinkElement: 'mailgun',
};

export const MailjetDefaultSettings: ProviderDefaultSettings = {
  name: 'mailjet',
  requiredTls: false,
  host: 'in-v3.mailjet.com',
  unencryptedPort: 587,
  encryptedPort: 465,
  user: { value: '', placeholder: 'Your Mailjet API key' },
  password: { value: '', placeholder: 'Your Mailjet Secret key' },
  image: './assets/images/smtp/mailjet.svg',
  senderEmailPlaceholder: 'An authorized domain or email address',
  routerLinkElement: 'mailjet',
};

export const PostmarkDefaultSettings: ProviderDefaultSettings = {
  name: 'postmark',
  requiredTls: false,
  host: 'smtp.postmarkapp.com',
  unencryptedPort: 587,
  encryptedPort: 587,
  user: { value: '', placeholder: 'Your Server API token' },
  password: { value: '', placeholder: 'Your Server API token' },
  image: './assets/images/smtp/postmark.png',
  senderEmailPlaceholder: 'An authorized domain or email address',
  routerLinkElement: 'postmark',
};

export const SendgridDefaultSettings: ProviderDefaultSettings = {
  name: 'sendgrid',
  requiredTls: false,
  host: 'smtp.sendgrid.net',
  unencryptedPort: 587,
  encryptedPort: 465,
  user: { value: 'apikey', placeholder: '' },
  password: { value: '', placeholder: ' Your SendGrid API Key' },
  image: './assets/images/smtp/sendgrid.png',
  routerLinkElement: 'sendgrid',
};

export const MailchimpDefaultSettings: ProviderDefaultSettings = {
  name: 'mailchimp',
  requiredTls: false,
  host: 'smtp.mandrillapp.com',
  unencryptedPort: 587,
  encryptedPort: 465,
  user: { value: '', placeholder: 'Your Mailchimp primary contact email' },
  password: { value: '', placeholder: 'Your Mailchimp Transactional API key' },
  image: './assets/images/smtp/mailchimp.svg',
  senderEmailPlaceholder: 'An authorized domain or email address',
  routerLinkElement: 'mailchimp',
};

export const BrevoDefaultSettings: ProviderDefaultSettings = {
  name: 'brevo',
  requiredTls: false,
  host: 'smtp-relay.sendinblue.com',
  unencryptedPort: 587,
  encryptedPort: 465,
  user: { value: '', placeholder: 'Your SMTP login email address' },
  password: { value: '', placeholder: 'Your SMTP key' },
  image: './assets/images/smtp/brevo.svg',
  routerLinkElement: 'brevo',
};

export const OutlookDefaultSettings: ProviderDefaultSettings = {
  name: 'outlook.com',
  requiredTls: true,
  host: 'smtp-mail.outlook.com',
  unencryptedPort: 587,
  encryptedPort: 587,
  user: { value: '', placeholder: 'Your outlook.com email address' },
  password: { value: '', placeholder: 'Your outlook.com password' },
  image: './assets/images/smtp/outlook.svg',
  senderEmailPlaceholder: 'Your outlook.com email address',
  routerLinkElement: 'outlook',
};

export const GenericDefaultSettings: ProviderDefaultSettings = {
  name: 'generic',
  requiredTls: false,
  user: { value: '', placeholder: 'your SMTP user' },
  password: { value: '', placeholder: 'your SMTP password' },
  routerLinkElement: 'generic',
};

export const SMTPKnownProviders = [
  AmazonSESDefaultSettings,
  BrevoDefaultSettings,
  // GoogleDefaultSettings,
  MailgunDefaultSettings,
  MailchimpDefaultSettings,
  MailjetDefaultSettings,
  PostmarkDefaultSettings,
  SendgridDefaultSettings,
  OutlookDefaultSettings,
  GenericDefaultSettings,
];
