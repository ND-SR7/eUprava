import axios from "axios";

const BASE_URL_MUP = process.env.REACT_APP_API_BASE_URL_MUP;
const BASE_URL_POLICE = process.env.REACT_APP_API_BASE_URL_POLICE;
const BASE_URL_COURT = process.env.REACT_APP_API_BASE_URL_COURT;
const BASE_URL_STATISTICS = process.env.REACT_APP_API_BASE_URL_STATISTICS;

const token = localStorage.getItem("token");

export async function pingMUP() {
  try {
    const response = await axios.get(`${BASE_URL_MUP}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to ping MUP');
  }
};

export async function pingPolice() {
  try {
    const response = await axios.get(`${BASE_URL_POLICE}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to ping traffic police');
  }
};

export async function pingCourt() {
  try {
    const response = await axios.get(`${BASE_URL_COURT}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to ping court');
  }
};

export async function pingStatistics() {
  try {
    const response = await axios.get(`${BASE_URL_STATISTICS}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to ping Institute for Statistics');
  }
};
