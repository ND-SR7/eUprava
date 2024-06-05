import axios from "axios";
import Credentials from "../models/User/Credentials";
import UserToken from "../models/User/UserToken";
import NewPerson from "../models/User/NewPerson";
import NewLegalEntity from "../models/User/NewLegalEntity";

const BASE_URL = process.env.REACT_APP_API_BASE_URL_SSO;

export async function login(data: Credentials) {
    try {
      const response = await axios.post(`${BASE_URL}/login`, data);
      return response.data as UserToken;
    } catch (error: any) {
      throw new Error(error.response.data.message || 'Failed to login user');
    }
};

export async function logout() {
  try {
    const response = await axios.get(`${BASE_URL}/logout`);
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to logout user');
  }
};

export async function registerPerson(data: NewPerson) {
  try {
    const response = await axios.post(`${BASE_URL}/register-person`, data);
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to register person');
  }
};

export async function registerLegalEntity(data: NewLegalEntity) {
  try {
    const response = await axios.post(`${BASE_URL}/register-entity`, data);
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to register legal entity');
  }
};

export async function sendRecoveryEmail(email: string) {
  try {
    const response = await axios.post(`${BASE_URL}/recover-password`, {email});
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to send recovery email');
  }
};
