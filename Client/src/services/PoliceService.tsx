import axios from "axios";

const BASE_URL_POLICE = process.env.REACT_APP_API_BASE_URL_POLICE;

export async function getAllTrafficViolations(){
    const token = localStorage.getItem("token")

    try {
        const response = await axios.get(`${BASE_URL_POLICE}/traffic-violation`, {
            headers: {
                Authorization:`Bearer ${token}`
            }
        });
        return response.data;
    } catch(error: any){
        throw new Error(error.response.data.message || 'Failed to retrieve traffic violations.');
    }
}

export async function getAllTrafficViolationsForUser() {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.get(`${BASE_URL_POLICE}/traffic-violation/jmbg`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to retrieve traffic violations.');
  }
}

export async function checkAlcoholLevel (data: any) {
  const token = localStorage.getItem("token");
  const response = await axios.post(`${BASE_URL_POLICE}/traffic-violation/check-alcohol-level`, data, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return response;
};

export async function checkVehicleTire (data: any) {
  const token = localStorage.getItem("token");
  const response = await axios.post(`${BASE_URL_POLICE}/traffic-violation/check-vehicle-tire`, data, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return response;
};

export async function checkDriverBan(data: any){
  const token = localStorage.getItem("token");
  const response = await axios.post(`${BASE_URL_POLICE}/traffic-violation/check-driver-ban`, data, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return response;
}

export async function CheckDriverPermitValidation(data: any){
  const token = localStorage.getItem("token");
  const response = await axios.post(`${BASE_URL_POLICE}/traffic-violation/check-driver-permit`, data, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return response.data;
}

export async function CheckVehicleRegistration(data: any){
  const token = localStorage.getItem("token");
  const response = await axios.post(`${BASE_URL_POLICE}/traffic-violation/check-vehicle-registration`, data, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return response.data;
}

export async function checkAll(data: any){
  const token = localStorage.getItem("token");
  const response = await axios.post(`${BASE_URL_POLICE}/traffic-violation/check-all`, data, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return response.data;
}