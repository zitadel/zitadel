"use client";

import { TextInput } from "#/ui/Input";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { ClientError } from "nice-grpc";

type Props = {
  userId?: string;
  isMe?: boolean;
  userState?: any; // UserState;
};

export default function Page() {
  const [passwordLoading, setPasswordLoading] = useState<boolean>(false);
  const [policyValid, setPolicyValid] = useState<boolean>(false);

  type Inputs = {
    password?: string;
    newPassword: string;
    confirmPassword: string;
  };

  const { register, handleSubmit, watch, reset, formState } = useForm<Inputs>({
    mode: "onChange",
    reValidateMode: "onChange",
    shouldUseNativeValidation: true,
  });

  const { errors, isValid } = formState;

  const watchNewPassword = watch("newPassword", "");
  const watchConfirmPassword = watch("confirmPassword", "");

  async function updatePassword(value: Inputs) {
    setPasswordLoading(true);

    // const authData: UpdateMyPasswordRequest = {
    //   oldPassword: value.password ?? '',
    //   newPassword: value.newPassword,
    // };

    const response = await fetch(
      `/api/user/password/me` +
        `?${new URLSearchParams({
          resend: `false`,
        })}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({}),
      }
    );

    if (response.ok) {
      setPasswordLoading(false);
      //   toast('password.set');
      // TODO: success info
      reset();
    } else {
      const error = (await response.json()) as ClientError;
      //   toast.error(error.details);
      // TODO: show error
    }
    setPasswordLoading(false);
  }

  async function sendHumanResetPasswordNotification(userId: string) {
    // const mgmtData: SendHumanResetPasswordNotificationRequest = {
    //   type: SendHumanResetPasswordNotificationRequest_Type.TYPE_EMAIL,
    //   userId: userId,
    // };

    const response = await fetch(`/api/user/password/resetlink/${userId}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({}),
    });

    if (response.ok) {
      // TODO: success info
      //   toast(t('sendPasswordResetLinkSent'));
    } else {
      const error = await response.json();
      // TODO: show error
      //   toast.error((error as ClientError).details);
    }
  }

  return (
    <>
      <h1 className="text-center">Set Password</h1>
      <p className="text-center my-4 mb-6 text-14px text-input-light-label dark:text-input-dark-label">
        Enter your new Password according to the requirements listed.
      </p>
      <form>
        <div>
          <TextInput
            type="password"
            required
            {...register("password", { required: true })}
            label="Password"
            error={errors.password?.message}
          />
        </div>
        <div className="mt-3">
          <TextInput
            type="password"
            required
            {...register("newPassword", { required: true })}
            label="New Password"
            error={errors.newPassword?.message}
          />
        </div>
        <div className="mt-3 mb-4">
          <TextInput
            type="password"
            required
            {...register("confirmPassword", {
              required: true,
            })}
            label="Confirm Password"
            error={errors.confirmPassword?.message}
          />
        </div>
      </form>
    </>
  );
}
