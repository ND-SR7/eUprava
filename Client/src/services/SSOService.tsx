import axios from "axios";
import Credentials from "../models/User/Credentials";
import UserToken from "../models/User/UserToken";

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

// TODO: Complete this function
export async function register(data: any) {
  try {
    const response = await axios.post(`${BASE_URL}/register`, data);
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to register user');
  }
};
