import axios from "axios";
import decodeJwtToken from "./JwtService";
import CourtHearing, { RescheduleCourtHearing } from "../models/Court/CourtHearing";
import Suspension from "../models/Court/Suspension";
import Warrant from "../models/Court/Warrant";

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

export async function createHearingPerson(hearing: CourtHearing) {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.post(`${BASE_URL_COURT}/create-hearing-person`, hearing,
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to create hearing for person');
  }
};

export async function createHearingLegalEntity(hearing: CourtHearing) {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.post(`${BASE_URL_COURT}/create-hearing-entity`, hearing,
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to create hearing for entity');
  }
};

export async function createSuspension(suspension: Suspension) {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.post(`${BASE_URL_COURT}/suspensions`, suspension,
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to create suspension');
  }
};

export async function createWarrant(warrant: Warrant) {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.post(`${BASE_URL_COURT}/warrants`, warrant,
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to create warrant');
  }
};

export async function updateHearingPerson(hearing: RescheduleCourtHearing) {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.put(`${BASE_URL_COURT}/update-hearing-person`, hearing,
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to update person hearing');
  }
};

export async function updateHearingLegalEntity(hearing: RescheduleCourtHearing) {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.put(`${BASE_URL_COURT}/update-hearing-entity`, hearing,
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );
    return response.data;
  } catch (error: any) {
    throw new Error(error.response.data.message || 'Failed to update entity hearing');
  }
};
