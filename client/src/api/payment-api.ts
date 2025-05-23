import { ResponseModel } from "../types";
import { BASE_URL } from "../utils/AppConst";
import { axiosAuth } from "./common";

export const CollectPaymentApi = async (
  token: string
): Promise<ResponseModel> => {
  try {
    const auth = axiosAuth();
    const response = await auth.get(`${BASE_URL}/transactions/payment`);
    return response.data;
  } catch (error) {
    console.log(error);
    return {
      message: "error occured",
    };
  }
};

export const VerifyPayment = async (token: string): Promise<ResponseModel> => {
  try {
    const auth = axiosAuth();
    const response = await auth.get(`${BASE_URL}/transactions/payment/verify`);
    return response.data;
  } catch (error) {
    console.log(error);
    return {
      message: "error occured",
    };
  }
};
