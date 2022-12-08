import { Endpoint } from "../../components/Endpoint";

export const endpoint = {
  render: Endpoint,
  children: ["paragraph"],
  attributes: {
    method: {
      type: String,
      default: "GET",
    },
    link: {
      type: String,
    },
  },
};
