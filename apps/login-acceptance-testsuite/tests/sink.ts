import axios from "axios";

export async function getOtpFromSink(key: string): Promise<any> {
  try {
    const response = await axios.post(
      process.env.SINK_NOTIFICATION_URL!,
      {
        recipient: key,
      },
      {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${process.env.ZITADEL_SERVICE_USER_TOKEN}`,
        },
      },
    );

    if (response.status >= 400) {
      const error = `HTTP Error: ${response.status} - ${response.statusText}`;
      console.error(error);
      throw new Error(error);
    }
    return response.data.args.oTP;
  } catch (error) {
    console.error("Error making request:", error);
    throw error;
  }
}

export async function getCodeFromSink(key: string): Promise<any> {
  try {
    const response = await axios.post(
      process.env.SINK_NOTIFICATION_URL!,
      {
        recipient: key,
      },
      {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${process.env.ZITADEL_SERVICE_USER_TOKEN}`,
        },
      },
    );

    if (response.status >= 400) {
      const error = `HTTP Error: ${response.status} - ${response.statusText}`;
      console.error(error);
      throw new Error(error);
    }
    return response.data.args.code;
  } catch (error) {
    console.error("Error making request:", error);
    throw error;
  }
}
