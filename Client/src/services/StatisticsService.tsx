import axios from "axios";

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
