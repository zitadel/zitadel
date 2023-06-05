import type { Meta, StoryObj } from '@storybook/react';

import VerifyEmailForm from './VerifyEmailForm';

const meta: Meta<typeof VerifyEmailForm> = {
  title: 'VerifyEmailForm',
  component: VerifyEmailForm,
};

export default meta;

type Story = StoryObj<typeof VerifyEmailForm>;

export const Default: Story = {
  args: {
    code: 'xyz',
    submit: true,
    userId: "123",
  }
};
