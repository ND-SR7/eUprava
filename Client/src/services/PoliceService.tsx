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