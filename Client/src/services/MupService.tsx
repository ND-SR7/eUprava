import axios from "axios";
import toast from "react-hot-toast";

const BASE_URL_MUP = process.env.REACT_APP_API_BASE_URL_MUP;

export async function getDrivingBans(){
    const token = localStorage.getItem("token");

    try {
        const response = await axios.get(`${BASE_URL_MUP}/driving-bans`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error: any) {
        throw new Error(error.response.data.message || 'Failed to retrieve driving bans');
    }
};

export async function getVehicles(){
    const token = localStorage.getItem("token");

    try {
        const response = await axios.get(`${BASE_URL_MUP}/persons-vehicles`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error: any) {
        throw new Error(error.response.data.message || 'Failed to retrieve persons vehicles');
    }
};

export async function getDrivingPermit() {
    const token = localStorage.getItem("token");

    try {
        const response = await axios.get(`${BASE_URL_MUP}/persons-driving-permit`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error: any) {
        throw new Error(error.response.data.message || 'Failed to retrieve driving permit');
    }
};

export async function getDrivingPermitRequests() {
    const token = localStorage.getItem("token");

    try {
        const response = await axios.get(`${BASE_URL_MUP}/pending-traffic-permit-requests`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error: any) {
        throw new Error(error.response.data.message || 'Failed to retrieve traffic permit requests');
    }
};

export async function getRegistrationRequests() {
    const token = localStorage.getItem("token");

    try {
        const response = await axios.get(`${BASE_URL_MUP}/pending-registration-requests`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error: any) {
        throw new Error(error.response.data.message || 'Failed to retrieve traffic permit requests');
    }
};

export const registerVehicle = async (vehicleID: string) => {
    const token = localStorage.getItem("token");

    const body = {
        registrationNumber: "registrationNumber",
        issuedDate: new Date().toISOString(),
        expirationDate: new Date().toISOString(),
        vehicleID: vehicleID,
        approved: true,
        owner: "owner",
        plates: "plates"
    };

    try {
        const response = await axios.post(`${BASE_URL_MUP}/registration-request`, body, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error: any) {
        throw new Error(error.response.data.message || 'Failed to register vehicle');
    }
};

export const approveRegistrationRequest = async (
    registrationNumber: string, 
    vehicleID: string, 
    owner: string,
    closeModal: () => void,
    refreshRequests: () => void
) => {
    const token = localStorage.getItem("token");

    const body = {
        registrationNumber: registrationNumber,
        issuedDate: new Date().toISOString(),
        expirationDate: new Date().toISOString(),
        vehicleID: vehicleID,
        approved: true,
        owner: owner,
        plates: ""
    };

    try {
        const response = await axios.post(`${BASE_URL_MUP}/approve-registration-request`, body, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        toast.success("Successfully approved registration request");
        closeModal();
        refreshRequests();
        return response.data;
    } catch (error: any) {
        toast.error("Failed to approve registration request");
        console.error(error);
        throw new Error(error.response.data.message || 'Failed to approve registration request');
    }
};

export const declineRegistrationRequest = async (
    registrationNumber: string,
    closeModal: () => void,
    refreshRequests: () => void
) => {
    const token = localStorage.getItem("token");

    try {
        const response = await axios.delete(`${BASE_URL_MUP}/delete-pending-registration-request/${registrationNumber}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        toast.success("Successfully declined registration request");
        closeModal();
        refreshRequests();
        return response.data;
    } catch (error: any) {
        toast.error("Failed to decline registration request");
        console.error(error);
        throw new Error(error.response.data.message || 'Failed to decline registration request');
    }
};

export const requestTrafficPermit = async () => {
    const token = localStorage.getItem("token");

    const body = {
        id: "",
        number: "",
        issuedDate: "2024-06-07T00:00:00Z",
        expirationDate: "2024-06-07T00:00:00Z",
        approved: false,
        person: ""
    };

    try {
        const response = await axios.post(`${BASE_URL_MUP}/traffic-permit-request`, body, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (error: any) {
        throw new Error(error.response.data.message || 'Failed to request traffic permit');
    }
};

export const approveDrivingPermitRequest = async (
    id: string, 
    number: string, 
    issuedDate: string, 
    expirationDate: string, 
    person: string, 
    closeModal: () => void, 
    refreshRequests: () => void
) => {
    const token = localStorage.getItem("token");

    const body = {
        id: id,
        number: number,
        issuedDate: issuedDate,
        expirationDate: expirationDate,
        approved: true,
        person: person
    };

    try {
        const response = await axios.post(`${BASE_URL_MUP}/approve-traffic-permit-request`, body, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        toast.success("Successfully approved driving permit request");
        closeModal();
        refreshRequests();
        return response.data;
    } catch (error: any) {
        toast.error("Failed to approve driving permit request");
        console.error(error);
        throw new Error(error.response.data.message || 'Failed to approve driving permit request');
    }
};

export const declineDrivingPermitRequest = async (
    id: string, 
    closeModal: () => void, 
    refreshRequests: () => void
) => {
    const token = localStorage.getItem("token");

    try {
        const response = await axios.delete(`${BASE_URL_MUP}/delete-pending-traffic-permit-request/${id}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        toast.success("Successfully declined driving permit request");
        closeModal();
        refreshRequests();
        return response.data;
    } catch (error: any) {
        toast.error("Failed to decline driving permit request");
        console.error(error);
        throw new Error(error.response.data.message || 'Failed to decline driving permit request');
    }
};

export const handleRegister = async (vehicleID: string) => {
    try {
        await registerVehicle(vehicleID);
        toast.success("Registration request sent successfully");
    } catch (error) {
        toast.error("Failed to send registration request");
    }
};

export const handleRequestPermit = async () => {
    try {
        await requestTrafficPermit();
        toast.success("Driving permit requested successfully");
    } catch (error) {
        toast.error("Failed to request driving permit");
    }
};