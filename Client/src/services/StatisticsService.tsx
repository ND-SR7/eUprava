import axios from "axios";
import BrandCount from "../models/Statistics/BrandCount";

const BASE_URL_STATISTICS = process.env.REACT_APP_API_BASE_URL_STATISTICS;

export const getVehicleStatisticsByYear = async (): Promise<{ [year: number]: number }> => {
  const token = localStorage.getItem("token");
  const response = await axios.get(`${BASE_URL_STATISTICS}/vehicle-statistics-by-year`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.data;
};

export const getTrafficStatistics = async () => {
  const token = localStorage.getItem("token");
  const response = await axios.get(`${BASE_URL_STATISTICS}/traffic-statistic`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.data;
};

export const getRegisteredVehiclesByYear = async (year: string) => {
  const token = localStorage.getItem("token");
  try {
    const response = await axios.get(`${BASE_URL_STATISTICS}/registered-vehicles/${year}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.data;
  } catch (error) {
    console.error('Error fetching registered vehicles:', error);
    throw error;
  }
};

export const getTrafficViolationsByYear = async (year: string) => {
  const token = localStorage.getItem("token");
  try {
    const response = await axios.get(`${BASE_URL_STATISTICS}/traffic-violations-report/${year}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.data;
  } catch (error) {
    console.error('Error fetching traffic violations:', error);
    throw error;
  }
};

export const getMostPopularBrands = async (year: string): Promise<BrandCount[]> => {
  const token = localStorage.getItem("token");
  try {
    const response = await axios.get(`${BASE_URL_STATISTICS}/most-popular-brands/${year}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.data || [];
  } catch (error) {
    console.error('Error fetching most popular brands:', error);
    throw error;
  }
};
