import axios from "axios";
import decodeJwtToken from "./JwtService";

const BASE_URL_COURT = process.env.REACT_APP_API_BASE_URL_COURT;

export async function getHearingsByJMBG() {
  const token = localStorage.getItem("token");
  const jmbg = decodeJwtToken(token!).sub;

  try {
    const response = await axios.get(`${BASE_URL_COURT}/hearings/${jmbg}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to retrieve hearings');
  }
};

export async function getSuspensionByJMBG() {
  const token = localStorage.getItem("token");
  const jmbg = decodeJwtToken(token!).sub;

  try {
    const response = await axios.get(`${BASE_URL_COURT}/suspensions/${jmbg}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to retrieve suspensions');
  }
};

export async function getWarrantsByJMBG() {
  const token = localStorage.getItem("token");
  const jmbg = decodeJwtToken(token!).sub;

  try {
    const response = await axios.get(`${BASE_URL_COURT}/warrants/${jmbg}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to retrieve warrants');
  }
};
