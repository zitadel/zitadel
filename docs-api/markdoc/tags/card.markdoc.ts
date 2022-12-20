import { Card } from "../../components/Card";

export const card = {
  render: Card,
  children: [""],
  attributes: {
    title: {
      type: String,
    },
    hasLanguageToggle: {
      type: Boolean,
    },
    hasProtocolToggle: {
      type: Boolean,
    },
    language: {
      type: String,
    },
    protocol: {
      type: String,
    },
  },
};
